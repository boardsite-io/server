package session

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/heat1q/boardsite/api/types"
)

// GetPages returns all pageIDs in order.
func (scb *ControlBlock) GetPages(ctx context.Context) ([]string, map[string]*types.PageMeta, error) {
	pageRank, err := scb.cache.GetPages(ctx, scb.ID)
	if err != nil {
		return nil, nil, errors.New("unable to fetch pages")
	}
	meta, err := scb.cache.GetPagesMeta(ctx, scb.ID, pageRank...)
	if err != nil {
		return nil, nil, errors.New("unable to fetch pages meta data")
	}
	return pageRank, meta, nil
}

// GetPagesSet returns all pageIDs in a map for fast verification.
func (scb *ControlBlock) GetPagesSet(ctx context.Context) map[string]struct{} {
	pageIDs, _ := scb.cache.GetPages(ctx, scb.ID)
	pageIDSet := make(map[string]struct{})

	for _, pid := range pageIDs {
		pageIDSet[pid] = struct{}{}
	}

	return pageIDSet
}

// IsValidPage checks if a pageID is valid, i.e. the page exists.
func (scb *ControlBlock) IsValidPage(ctx context.Context, pageID ...string) bool {
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
func (scb *ControlBlock) AddPages(ctx context.Context, pageIDs []string, index []int, meta map[string]*types.PageMeta) error {
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
		if err := scb.cache.AddPage(ctx, scb.ID, pid, index[i], pMeta); err != nil {
			return errors.New("cannot add page")
		}
	}

	return nil
}

// DeletePages delete pages with pageID and broadcasts
// the change to all connected clients.
func (scb *ControlBlock) DeletePages(ctx context.Context, pageID ...string) error {
	defer scb.SyncPages(ctx)

	var sb strings.Builder
	// go through all pages even if some fail
	for _, pid := range pageID {
		if !scb.IsValidPage(ctx, pid) {
			sb.WriteString(fmt.Sprintf(": page %s does not exist", pid))
			continue
		}
		if err := scb.cache.DeletePage(ctx, scb.ID, pid); err != nil {
			sb.WriteString(fmt.Sprintf(": cannot delete page %s", pageID))
		}
	}

	if sb.Len() > 0 {
		return fmt.Errorf("error deleting pages%s", sb.String())
	}

	return nil
}

// UpdatePages modifies the page meta data and/or clears the content.
func (scb *ControlBlock) UpdatePages(ctx context.Context, pageIDs []string, meta map[string]*types.PageMeta, clear bool) error {
	var pageIDsUpdate []string

	defer func() {
		if len(pageIDsUpdate) == 0 {
			return
		}
		scb.broadcast <- &types.Message{
			Type:   types.MessageTypePageUpdate,
			Sender: "", // send to all clients
			Content: types.ContentPageRequest{
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
			if err := scb.cache.ClearPage(ctx, scb.ID, pid); err != nil {
				return fmt.Errorf("cannot clear page %s", pid)
			}
		} else {
			pMeta, ok := meta[pid]
			if !ok {
				return fmt.Errorf("no meta given for page %s", pid)
			}
			// update db
			if err := scb.cache.UpdatePageMeta(ctx, scb.ID, pid, pMeta); err != nil {
				return err
			}
		}
		pageIDsUpdate = append(pageIDsUpdate, pid)
	}

	return nil
}

// SyncPages broadcasts the current PageRank to all connected
// clients indicating an update in the pages (or ordering).
func (scb *ControlBlock) SyncPages(ctx context.Context) error {
	pageRank, meta, err := scb.GetPages(ctx)
	if err != nil {
		return err
	}

	scb.broadcast <- &types.Message{
		Type:   types.MessageTypePageSync,
		Sender: "", // send to all clients
		Content: &types.ContentPageSync{
			PageRank: pageRank,
			Meta:     meta,
		},
	}
	return nil
}
