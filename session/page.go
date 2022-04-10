package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	apiErrors "github.com/heat1q/boardsite/api/errors"

	"github.com/heat1q/boardsite/redis"

	"github.com/heat1q/boardsite/api/log"
	"github.com/heat1q/boardsite/api/types"
)

const (
	QueryKeyUpdate = "update"
)

const (
	updateOperationMeta   = "meta"
	updateOperationClear  = "clear"
	updateOperationDelete = "delete"
)

// PageStyle declares the style of the page background.
type PageBackground struct {
	// page background
	Paper    int `json:"paper,omitempty"`
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

type Page struct {
	PageId  string     `json:"pageId"`
	Meta    *PageMeta  `json:"meta"`
	Strokes *[]*Stroke `json:"strokes,omitempty"` //nullable
}

// PageRequest declares the message content for page requests.
type PageRequest struct {
	PageID  []string                       `json:"pageId"`
	Index   []int                          `json:"index,omitempty"`
	Meta    map[string]*PageMeta           `json:"meta"`
	Strokes *map[string]map[string]*Stroke `json:"strokes,omitempty"`
}

type PageSync struct {
	PageRank []string         `json:"pageRank"`
	Pages    map[string]*Page `json:"pages"`
}

func (scb *controlBlock) GetPageRank(ctx context.Context) ([]string, error) {
	return scb.cache.GetPageRank(ctx, scb.cfg.ID)
}

func (scb *controlBlock) GetPage(ctx context.Context, pageId string, withStrokes bool) (*Page, error) {
	page := Page{
		PageId: pageId,
		Meta:   &PageMeta{},
	}

	err := scb.cache.GetPageMeta(ctx, scb.cfg.ID, pageId, page.Meta)
	if err != nil {
		return nil, err
	}

	if withStrokes {
		strokes, err := scb.getStrokes(ctx, pageId)
		if err != nil {
			return nil, err
		}
		page.Strokes = &strokes
	}

	return &page, nil
}

// AddPages adds a page with pageID to the session and broadcasts
// the change to all connected clients.
func (scb *controlBlock) AddPages(ctx context.Context, pageRequest PageRequest) error {
	if len(pageRequest.PageID) != len(pageRequest.Index) {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("cannot find page index"))
	}
	if scb.IsValidPage(ctx, pageRequest.PageID...) {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("some pages already exist"))
	}

	defer scb.broadcastPageSync(ctx, pageRequest.PageID, pageRequest.Strokes != nil)

	for i, pid := range pageRequest.PageID {
		pMeta, ok := pageRequest.Meta[pid]
		if !ok {
			return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("no meta given for page %s", pid))
		}
		if err := scb.cache.AddPage(ctx, scb.cfg.ID, pid, pageRequest.Index[i], pMeta); err != nil {
			return errors.New("cannot add page")
		}
		if pageRequest.Strokes != nil {
			strokeMap, ok := (*pageRequest.Strokes)[pid]
			if !ok {
				continue
			}
			strokes := make([]redis.Stroke, 0, len(strokeMap))
			for _, s := range strokeMap {
				strokes = append(strokes, s)
			}
			if err := scb.cache.UpdateStrokes(ctx, scb.cfg.ID, strokes...); err != nil {
				return fmt.Errorf("update page strokes: %w", err)
			}
		}
	}

	return nil
}

// UpdatePages modifies the page meta data and/or clears the content.
func (scb *controlBlock) UpdatePages(ctx context.Context, pageRequest PageRequest, operation string) error {
	switch operation {
	case updateOperationMeta:
		return scb.updatePagesMeta(ctx, pageRequest.Meta)

	case updateOperationDelete:
		return scb.deletePages(ctx, pageRequest.PageID...)

	case updateOperationClear:
		return scb.clearPages(ctx, pageRequest.PageID...)

	default:
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("unknown operation: %s", operation))
	}
}

func (scb *controlBlock) GetPageSync(ctx context.Context, pageIds []string, withStrokes bool) (*PageSync, error) {
	var (
		sync PageSync
		err  error
	)

	sync.PageRank, err = scb.cache.GetPageRank(ctx, scb.cfg.ID)
	if err != nil {
		return nil, err
	}

	sync.Pages, err = scb.getPages(ctx, pageIds, withStrokes)
	if err != nil {
		return nil, err
	}

	return &sync, nil
}

func (scb *controlBlock) SyncSession(ctx context.Context, sync PageSync) error {
	if err := scb.cache.ClearSession(ctx, scb.cfg.ID); err != nil {
		return err
	}

	defer scb.broadcastPageSync(ctx, sync.PageRank, true)

	for _, pid := range sync.PageRank {
		page, ok := sync.Pages[pid]
		if !ok {
			return fmt.Errorf("page %s not found", pid)
		}

		if err := scb.cache.AddPage(ctx, scb.cfg.ID, pid, -1, page.Meta); err != nil {
			return err
		}

		var strokes []redis.Stroke
		if page.Strokes != nil {
			strokes = make([]redis.Stroke, 0, len(*page.Strokes))
			for _, s := range *page.Strokes {
				strokes = append(strokes, s)
			}
		}

		if err := scb.cache.UpdateStrokes(ctx, scb.cfg.ID, strokes...); err != nil {
			return err
		}
	}

	return nil
}

func (scb *controlBlock) getPages(ctx context.Context, pageIds []string, withStrokes bool) (map[string]*Page, error) {
	pages := make(map[string]*Page, len(pageIds))
	for _, pid := range pageIds {
		page, err := scb.GetPage(ctx, pid, withStrokes)
		if err != nil {
			return nil, err
		}
		pages[pid] = page
	}

	return pages, nil
}

func (scb *controlBlock) getStrokes(ctx context.Context, pageId string) ([]*Stroke, error) {
	strokeBytes, err := scb.cache.GetPageStrokes(ctx, scb.cfg.ID, pageId)
	if err != nil {
		return nil, err
	}

	strokes := make([]*Stroke, len(strokeBytes))
	for i, s := range strokeBytes {
		var stroke Stroke
		if err := json.Unmarshal(s, &stroke); err != nil {
			return nil, err
		}
		strokes[i] = &stroke
	}

	return strokes, nil
}

// IsValidPage checks if a pageID is valid, i.e. the page exists.
func (scb *controlBlock) IsValidPage(ctx context.Context, pageID ...string) bool {
	pages := scb.getPagesSet(ctx)
	for _, pid := range pageID {
		if _, ok := pages[pid]; !ok {
			return false
		}
	}
	return true
}

// GetPagesSet returns all pageIDs in a map for fast verification.
func (scb *controlBlock) getPagesSet(ctx context.Context) map[string]struct{} {
	pageIDs, _ := scb.cache.GetPageRank(ctx, scb.cfg.ID)
	pageIDSet := make(map[string]struct{})

	for _, pid := range pageIDs {
		pageIDSet[pid] = struct{}{}
	}

	return pageIDSet
}

func (scb *controlBlock) updatePagesMeta(ctx context.Context, meta map[string]*PageMeta) error {
	updates := make([]string, 0, len(meta))
	for pid, m := range meta {
		if !scb.IsValidPage(ctx, pid) {
			continue
		}

		// update db
		var newMeta PageMeta
		if err := scb.cache.GetPageMeta(ctx, scb.cfg.ID, pid, &newMeta); err != nil {
			return err
		}
		tmp, err := json.Marshal(m)
		if err != nil {
			return err
		}
		// overwrite non-zero values
		if err := json.Unmarshal(tmp, &newMeta); err != nil {
			return err
		}
		if err := scb.cache.SetPageMeta(ctx, scb.cfg.ID, pid, newMeta); err != nil {
			return err
		}

		updates = append(updates, pid)
	}

	scb.broadcastPageSync(ctx, updates, false)

	return nil
}

// DeletePages delete pages with pageID and broadcasts
// the change to all connected clients.
func (scb *controlBlock) deletePages(ctx context.Context, pageID ...string) error {
	defer scb.broadcastPageSync(ctx, nil, false)

	var sb strings.Builder
	// go through all pages even if some fail
	for _, pid := range pageID {
		if !scb.IsValidPage(ctx, pid) {
			sb.WriteString(fmt.Sprintf(": page %s does not exist", pid))
			continue
		}
		if err := scb.cache.DeletePage(ctx, scb.cfg.ID, pid); err != nil {
			sb.WriteString(fmt.Sprintf(": cannot delete page %s", pageID))
		}
	}

	if sb.Len() > 0 {
		return fmt.Errorf("error deleting pages%s", sb.String())
	}

	return nil
}

func (scb *controlBlock) clearPages(ctx context.Context, pageIds ...string) error {
	defer scb.broadcastPageSync(ctx, pageIds, true)
	for _, pid := range pageIds {
		if err := scb.cache.ClearPage(ctx, scb.cfg.ID, pid); err != nil {
			return fmt.Errorf("clear page %s: %w", pid, err)
		}
	}
	return nil
}

// SyncPages broadcasts the current PageRank to all connected
// clients indicating an update in the pages (or ordering).
// page with ids specified in the pageIds slice will be broadcasted
func (scb *controlBlock) broadcastPageSync(ctx context.Context, pageIds []string, withStrokes bool) {
	sync, err := scb.GetPageSync(ctx, pageIds, withStrokes)
	if err != nil {
		log.Ctx(ctx).Errorf("failed to broadcast page sync: get sync: %v", err)
		return
	}

	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypePageSync,
		Sender:  "", // send to all clients
		Content: sync,
	}
}
