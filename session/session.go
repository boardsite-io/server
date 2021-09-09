package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/heat1q/boardsite/api/types"
)

// GetStrokes fetches all stroke data for specified page.
func (scb *ControlBlock) GetStrokes(ctx context.Context, pageID string) ([]types.Stroke, error) {
	strokesRaw, err := scb.cache.FetchStrokesRaw(ctx, scb.ID, pageID)
	if err != nil {
		return nil, errors.New("unable to fetch strokes")
	}

	strokes := make([]types.Stroke, len(strokesRaw))
	for i, s := range strokesRaw {
		if err := json.Unmarshal(s, &strokes[i]); err != nil {
			return nil, err
		}
	}
	return strokes, nil
}

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

// NewUser generate a new user struct based on
// the alias and color attribute
//
// Does some sanitize checks.
func (scb *ControlBlock) NewUser(alias, color string) (*types.User, error) {
	if len(alias) > 24 {
		alias = alias[:24]
	}
	//TODO check if html color ?
	if len(color) != 7 {
		return nil, fmt.Errorf("incorrect html color")
	}

	id, err := gonanoid.New(16)
	if err != nil {
		return nil, err
	}
	user := &types.User{
		ID:    id,
		Alias: alias,
		Color: color,
	}
	// set user waiting
	scb.UserReady(user)
	return user, err
}

// Receive is the entry point when a message is received in
// the session via the websocket.
func (scb *ControlBlock) Receive(ctx context.Context, msg *types.Message) error {
	if !scb.IsUserConnected(msg.Sender) {
		return errors.New("invalid sender userId")
	}

	var err error
	switch msg.Type {
	case types.MessageTypeStroke:
		err = scb.sanitizeStrokes(ctx, msg)

	case types.MessageTypeMouseMove:
		err = scb.mouseMove(msg)

	default:
		err = fmt.Errorf("message type not recognized: %s", msg.Type)
	}
	return err
}

// sanitizeStrokes parses the stroke content of the message.
//
// It further checks if the strokes have a valid pageId and userId.
func (scb *ControlBlock) sanitizeStrokes(ctx context.Context, msg *types.Message) error {
	var strokes []*types.Stroke
	if err := msg.UnmarshalContent(&strokes); err != nil {
		return err
	}

	validStrokes := make([]*types.Stroke, 0, len(strokes))
	pageIDs := scb.GetPagesSet(ctx)

	for _, stroke := range strokes {
		if _, ok := pageIDs[stroke.GetPageID()]; ok { // valid pageID
			if stroke.GetUserID() == msg.Sender { // valid userID
				validStrokes = append(validStrokes, stroke)
			}
		}
	}
	if len(validStrokes) > 0 {
		scb.updateStrokes(msg.Sender, validStrokes)
		return nil
	}
	return errors.New("strokes not validated")
}

// updateStrokes updates the strokes in the session with sessionID.
//
// userID indicates the initiator of the message, which is
// to be excluded in the broadcast. The strokes are scheduled for an
// update to Redis.
func (scb *ControlBlock) updateStrokes(userID string, strokes []*types.Stroke) {
	// broadcast changes
	scb.broadcast <- &types.Message{
		Type:    types.MessageTypeStroke,
		Sender:  userID,
		Content: strokes,
	}

	// save to database
	scb.dbUpdate <- strokes
}

// mouseMove broadcast mouse move events.
func (scb *ControlBlock) mouseMove(msg *types.Message) error {
	var mouseUpdate types.ContentMouseMove
	if err := msg.UnmarshalContent(&mouseUpdate); err != nil {
		return err
	}
	scb.broadcast <- &types.Message{
		Type:    types.MessageTypeMouseMove,
		Sender:  msg.Sender,
		Content: mouseUpdate,
	}
	return nil
}
