package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/heat1q/boardsite/api/types"
)

// PageStyle declares the style of the page background.
type PageBackground struct {
	// page background
	Style    string `json:"style,omitempty"`
	PageNum  int    `json:"documentPageNum"`
	AttachId string `json:"attachId"`
}

type PageSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// PageMeta declares some page meta data.
type PageMeta struct {
	PageSize   PageSize       `json:"size"`
	Background PageBackground `json:"background"`
}

// ContentPageRequest declares the message content for page requests.
type ContentPageRequest struct {
	PageID []string             `json:"pageId"`
	Index  []int                `json:"index,omitempty"`
	Clear  bool                 `json:"clear,omitempty"`
	Meta   map[string]*PageMeta `json:"meta"`
}

// ContentPageSync message content for page sync.
type ContentPageSync struct {
	PageRank []string             `json:"pageRank"`
	Meta     map[string]*PageMeta `json:"meta"`
}

//type Page struct {
//	PageId  string          `json:"pageId"`
//	Strokes []*Stroke       `json:"strokes"`
//	Meta    *PageMeta `json:"meta"`
//}

// GetPages returns all pageIDs in order.
func (scb *controlBlock) GetPages(ctx context.Context) ([]string, map[string]*PageMeta, error) {
	pageRank, err := scb.cache.GetPageRank(ctx, scb.id)
	if err != nil {
		return nil, nil, fmt.Errorf("cache: get page rank: %w", err)
	}
	meta := make(map[string]*PageMeta, len(pageRank))
	for _, pid := range pageRank {
		var m PageMeta
		if err := scb.cache.GetPageMeta(ctx, scb.id, pid, &m); err != nil {
			return nil, nil, fmt.Errorf("cache: get page meta: %w", err)
		}
		meta[pid] = &m
	}
	return pageRank, meta, nil
}

// GetPagesSet returns all pageIDs in a map for fast verification.
func (scb *controlBlock) GetPagesSet(ctx context.Context) map[string]struct{} {
	pageIDs, _ := scb.cache.GetPageRank(ctx, scb.id)
	pageIDSet := make(map[string]struct{})

	for _, pid := range pageIDs {
		pageIDSet[pid] = struct{}{}
	}

	return pageIDSet
}

// IsValidPage checks if a pageID is valid, i.e. the page exists.
func (scb *controlBlock) IsValidPage(ctx context.Context, pageID ...string) bool {
	pages := scb.GetPagesSet(ctx)
	for _, pid := range pageID {
		if _, ok := pages[pid]; !ok {
			return false
		}
	}
	return true
}

// AddPages adds a page with pageID to the session and broadcasts
// the change to all connected clients.
func (scb *controlBlock) AddPages(ctx context.Context, pageIDs []string, index []int, meta map[string]*PageMeta) error {
	if len(pageIDs) != len(index) {
		return errors.New("cannot find page index")
	}
	if scb.IsValidPage(ctx, pageIDs...) {
		return errors.New("some pages already exist")
	}

	defer scb.SyncPages(ctx)

	for i, pid := range pageIDs {
		pMeta, ok := meta[pid]
		if !ok {
			return fmt.Errorf("no meta given for page %s", pid)
		}
		if err := scb.cache.AddPage(ctx, scb.id, pid, index[i], pMeta); err != nil {
			return errors.New("cannot add page")
		}
	}

	return nil
}

// DeletePages delete pages with pageID and broadcasts
// the change to all connected clients.
func (scb *controlBlock) DeletePages(ctx context.Context, pageID ...string) error {
	defer scb.SyncPages(ctx)

	var sb strings.Builder
	// go through all pages even if some fail
	for _, pid := range pageID {
		if !scb.IsValidPage(ctx, pid) {
			sb.WriteString(fmt.Sprintf(": page %s does not exist", pid))
			continue
		}
		if err := scb.cache.DeletePage(ctx, scb.id, pid); err != nil {
			sb.WriteString(fmt.Sprintf(": cannot delete page %s", pageID))
		}
	}

	if sb.Len() > 0 {
		return fmt.Errorf("error deleting pages%s", sb.String())
	}

	return nil
}

// UpdatePages modifies the page meta data and/or clears the content.
func (scb *controlBlock) UpdatePages(ctx context.Context, pageIDs []string, meta map[string]*PageMeta, clear bool) error {
	var pageIDsUpdate []string

	defer func() {
		if len(pageIDsUpdate) == 0 {
			return
		}
		scb.broadcast <- &types.Message{
			Type:   types.MessageTypePageUpdate,
			Sender: "", // send to all clients
			Content: ContentPageRequest{
				PageID: pageIDsUpdate,
				Clear:  clear,
				Meta:   meta,
			},
		}
	}()

	for _, pid := range pageIDs {
		if !scb.IsValidPage(ctx, pid) {
			return fmt.Errorf("page %s does not exits", pid)
		}
		if clear {
			if err := scb.cache.ClearPage(ctx, scb.id, pid); err != nil {
				return fmt.Errorf("cannot clear page %s", pid)
			}
		} else {
			pMeta, ok := meta[pid]
			if !ok {
				return fmt.Errorf("no meta given for page %s", pid)
			}
			// update db
			var newMeta PageMeta
			if err := scb.cache.GetPageMeta(ctx, scb.id, pid, &newMeta); err != nil {
				return err
			}
			tmp, err := json.Marshal(pMeta)
			if err != nil {
				return err
			}
			// overwrite non-zero values
			if err := json.Unmarshal(tmp, &newMeta); err != nil {
				return err
			}
			if err := scb.cache.SetPageMeta(ctx, scb.id, pid, newMeta); err != nil {
				return err
			}
		}
		pageIDsUpdate = append(pageIDsUpdate, pid)
	}

	return nil
}

// SyncPages broadcasts the current PageRank to all connected
// clients indicating an update in the pages (or ordering).
func (scb *controlBlock) SyncPages(ctx context.Context) error {
	pageRank, meta, err := scb.GetPages(ctx)
	if err != nil {
		return err
	}

	scb.broadcast <- &types.Message{
		Type:   types.MessageTypePageSync,
		Sender: "", // send to all clients
		Content: &ContentPageSync{
			PageRank: pageRank,
			Meta:     meta,
		},
	}
	return nil
}
