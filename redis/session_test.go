package redis_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/heat1q/boardsite/session"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/redis"
)

func setupHandler(t *testing.T) (*miniredis.Miniredis, redis.Handler) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err.Error())
	}
	port, err := strconv.ParseInt(mr.Port(), 10, 32)
	if err != nil {
		t.Fatal(err.Error())
	}
	h, err := redis.New(mr.Host(), uint16(port))
	if err != nil {
		t.Fatal(err.Error())
	}
	return mr, h
}

func genStroke(id, pageID string, strokeType int) *session.Stroke {
	return &session.Stroke{
		ID:     id,
		PageID: pageID,
		Type:   strokeType,
		X:      2.32,
		Y:      5.23,
		Points: []float64{324.42, 426.23, 123.34, 316.3, 324.42, 426.23, 123.34, 316.3},
		Style:  session.Style{Color: "#00beef", Width: 3.0},
	}
}

func Test_handler_UpdateStrokes(t *testing.T) {
	ctx := context.Background()
	mr, h := setupHandler(t)
	defer mr.Close()
	defer h.ClosePool()
	sid := "sid"
	pageId := "pageId"

	tests := []struct {
		name    string
		strokes []redis.Stroke
		updates []redis.Stroke
		want    []*session.Stroke
		wantErr bool
	}{
		{
			name: "add a stroke",
			updates: []redis.Stroke{
				genStroke("stroke1", pageId, 1),
			},
			want: []*session.Stroke{
				genStroke("stroke1", pageId, 1),
			},
		},
		{
			name: "delete a stroke",
			strokes: []redis.Stroke{
				genStroke("stroke1", pageId, 1),
			},
			updates: []redis.Stroke{
				genStroke("stroke1", pageId, 0),
			},
			want: []*session.Stroke{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr.FlushAll()
			_ = h.UpdateStrokes(ctx, sid, tt.strokes...)

			err := h.UpdateStrokes(ctx, sid, tt.updates...)

			assert.NoError(t, err)
			strokes, err := h.GetPageStrokes(ctx, sid, pageId)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(strokes))
			for i, s := range strokes {
				var got session.Stroke
				err = json.Unmarshal(s, &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want[i], &got)
			}
		})
	}
}

func Test_handler_GetStrokesRaw(t *testing.T) {
	ctx := context.Background()
	mr, h := setupHandler(t)
	defer mr.Close()
	defer h.ClosePool()
	sid := "sid"
	pageId := "pageId"
	want := genStroke("stroke1", pageId, 1)

	err := h.AddPage(ctx, sid, pageId, -1, nil)
	assert.NoError(t, err)
	err = h.UpdateStrokes(ctx, sid, want)
	assert.NoError(t, err)

	strokes, err := h.GetPageStrokes(ctx, sid, pageId)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(strokes))

	var got session.Stroke
	err = json.Unmarshal(strokes[0], &got)
	assert.NoError(t, err)
	assert.Equal(t, want, &got)
}

func Test_handler_AddPage(t *testing.T) {
	ctx := context.Background()
	mr, h := setupHandler(t)
	defer mr.Close()
	defer h.ClosePool()

	sid := "sid2"
	_ = h.ClearSession(ctx, sid)

	tests := []struct {
		pid   string
		index int
		want  []string
	}{
		{"pid1", 0, []string{"pid1"}},
		{"pid2", 1, []string{"pid1", "pid2"}},
		{"pid3", -1, []string{"pid1", "pid2", "pid3"}},
		{"pid4", 0, []string{"pid4", "pid1", "pid2", "pid3"}},
		{"pid5", 999, []string{"pid4", "pid1", "pid2", "pid3", "pid5"}},
		{"pid6", 2, []string{"pid4", "pid1", "pid6", "pid2", "pid3", "pid5"}},
	}

	for _, test := range tests {
		_ = h.AddPage(ctx, sid, test.pid, test.index, nil)
		pids, err := h.GetPageRank(ctx, sid)
		assert.NoError(t, err)
		assert.Equal(t, test.want, pids, "pageRank is not correct")
	}
}

func Test_handler_Get_SetPageMeta(t *testing.T) {
	ctx := context.Background()
	mr, h := setupHandler(t)
	defer mr.Close()
	defer h.ClosePool()
	sid := "sid"
	pageId := "pageId"
	want := session.PageMeta{
		PageSize: session.PageSize{
			Width:  786,
			Height: 1024,
		},
		Background: session.PageBackground{
			Paper:    "doc",
			PageNum:  0,
			AttachId: "attachId",
		},
	}

	err := h.AddPage(ctx, sid, pageId, -1, nil)
	assert.NoError(t, err)

	err = h.SetPageMeta(ctx, sid, pageId, &want)
	assert.NoError(t, err)

	var got session.PageMeta
	err = h.GetPageMeta(ctx, sid, pageId, &got)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_handler_DeletePage(t *testing.T) {
	ctx := context.Background()
	mr, h := setupHandler(t)
	defer mr.Close()
	defer h.ClosePool()
	sid := "sid"
	pageId := "pageId"

	err := h.AddPage(ctx, sid, pageId, -1, nil)
	assert.NoError(t, err)
	err = h.UpdateStrokes(ctx, sid, genStroke("stroke1", pageId, 1))
	assert.NoError(t, err)

	err = h.DeletePage(ctx, sid, pageId)

	assert.NoError(t, err)
	pageRank, err := h.GetPageRank(ctx, sid)
	assert.NoError(t, err)
	assert.Empty(t, pageRank)
	strokes, err := h.GetPageStrokes(ctx, sid, pageId)
	assert.NoError(t, err)
	assert.Empty(t, strokes)
}

func Test_handler_ClearSession(t *testing.T) {
	ctx := context.Background()
	mr, h := setupHandler(t)
	defer mr.Close()
	defer h.ClosePool()
	sid := "sid"
	pageId := "pageId"

	err := h.AddPage(ctx, sid, pageId, -1, nil)
	assert.NoError(t, err)
	err = h.UpdateStrokes(ctx, sid, genStroke("stroke1", pageId, 1))
	assert.NoError(t, err)

	err = h.ClearSession(ctx, sid)

	assert.NoError(t, err)
	pageRank, err := h.GetPageRank(ctx, sid)
	assert.NoError(t, err)
	assert.Empty(t, pageRank)
	strokes, err := h.GetPageStrokes(ctx, sid, pageId)
	assert.NoError(t, err)
	assert.Empty(t, strokes)
}
