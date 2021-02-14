package redis

import (
	"fmt"
	"github.com/heat1q/boardsite/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"math/rand"
	"os"
	"testing"
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
		AddPage(sid, test.pid, test.index)
		assert.Equal(t, test.want, GetPages(sid), "pageRank is not correct")
	}
}

func TestUpdateAndFetchStroke(t *testing.T) {
	if err := setupConn(); err != nil {
		t.Log("cannot connect to local Redis instance")
		t.SkipNow()
	}
	defer ClosePool()

	ClearSession("sid")
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

	refStrokeStr, _ := refStroke.JSONStringify()
	assert.Equal(
		t,
		fmt.Sprintf("[%s]", refStrokeStr),
		FetchStrokes("sid1", "pid1"),
		"incorrect json stringified array of strokes objects",
	)
}
