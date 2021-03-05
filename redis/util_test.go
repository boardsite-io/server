package redis

import (
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/heat1q/boardsite/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupConn() error {
	os.Setenv("B_REDIS_HOST", "localhost")
	os.Setenv("B_REDIS_PORT", "6379")

	return InitPool()
}

func genRandStroke(id, pageID string, strokeType int) *types.Stroke {
	pts := make([]float64, 20)
	for i := range pts {
		pts[i] = math.Floor(rand.Float64()*1e3) / 10.0
	}

	return &types.Stroke{
		ID:     id,
		PageID: pageID,
		Type:   strokeType,
		X:      math.Floor(rand.Float64()*1e3) / 10.0,
		Y:      math.Floor(rand.Float64()*1e3) / 10.0,
		Points: pts,
		Style:  types.Style{Color: "#00beef", Width: 3.0},
	}
}

func TestAddPages(t *testing.T) {
	if err := setupConn(); err != nil {
		t.Log("cannot connect to local Redis instance")
		t.SkipNow()
	}
	defer ClosePool()

	sid := "sid2"
	ClearSession(sid)

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
		AddPage(sid, test.pid, test.index, nil)
		pids, err := GetPages(sid)
		assert.NoError(t, err)
		assert.Equal(t, test.want, pids, "pageRank is not correct")
	}
}

func TestMetaPages(t *testing.T) {
	if err := setupConn(); err != nil {
		t.Log("cannot connect to local Redis instance")
		t.SkipNow()
	}
	defer ClosePool()

	sid := "sid2"
	ClearSession(sid)

	tests := []struct {
		meta  types.PageMeta
		index int
		want  types.PageMeta
	}{
		{types.PageMeta{Background: "somebackground1"}, 0, types.PageMeta{Background: "somebackground1"}},
	}

	for _, test := range tests {
		AddPage(sid, "pid1", test.index, &test.meta)
		meta, err := GetPagesMeta(sid, []string{"pid1"})
		assert.NoError(t, err)
		assert.Equal(t, test.want, *meta[0], "pageRank is not correct")
	}
}

func TestUpdateAndFetchStroke(t *testing.T) {
	if err := setupConn(); err != nil {
		t.Log("cannot connect to local Redis instance")
		t.SkipNow()
	}
	defer ClosePool()

	ClearSession("sid1")
	refStroke := genRandStroke("id1", "pid1", 1)

	setData := []*types.Stroke{
		refStroke,
		genRandStroke("id2", "pid1", rand.Intn(10)+1),
		genRandStroke("id3", "pid1", rand.Intn(10)+1),
		genRandStroke("id4", "pid1", rand.Intn(10)+1),
		genRandStroke("id5", "pid1", rand.Intn(10)+1),
	}
	require.NoError(t, Update("sid1", setData))

	clearData := []*types.Stroke{
		// delete strokes
		&types.Stroke{ID: "id2", PageID: "pid1", Type: 0},
		&types.Stroke{ID: "id5", PageID: "pid1", Type: 0},
		&types.Stroke{ID: "id4", PageID: "pid1", Type: 0},
		&types.Stroke{ID: "id3", PageID: "pid1", Type: 0},
	}
	require.NoError(t, Update("sid1", clearData))

	refStrokeStr, errRef := refStroke.JSONStringify()
	require.NoError(t, errRef)

	raw, errFetch := FetchStrokesRaw("sid1", "pid1")
	require.NoError(t, errFetch)
	assert.Equal(
		t,
		refStrokeStr,
		raw[0],
		"incorrect json stringified array of strokes objects",
	)
}

func TestFetchStrokesRaw(t *testing.T) {
	if err := setupConn(); err != nil {
		t.Log("cannot connect to local Redis instance")
		t.SkipNow()
	}
	defer ClosePool()

	setData := []*types.Stroke{
		genRandStroke("id1", "pid1", rand.Intn(10)+1),
		genRandStroke("id2", "pid1", rand.Intn(10)+1),
		genRandStroke("id3", "pid2", rand.Intn(10)+1),
		genRandStroke("id4", "pid1", rand.Intn(10)+1),
		genRandStroke("id5", "pid2", rand.Intn(10)+1),
	}
	require.NoError(t, Update("sid1", setData))

	tests := []struct {
		sid     string
		pid     string
		wantLen int
	}{
		{"", "", 0},
		{"sid0", "pid1", 0},
		{"sid1", "pid1", 3},
		{"sid1", "pid2", 2},
		{"sid1", "pid3", 0},
	}

	for _, test := range tests {
		raw, err := FetchStrokesRaw(test.sid, test.pid)
		assert.NoError(t, err)
		assert.Equal(t, test.wantLen, len(raw), "wrong number of fetched strokes")
	}
}
