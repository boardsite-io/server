package session_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/heat1q/boardsite/api/types"

	"github.com/heat1q/boardsite/attachment/attachmentfakes"

	"github.com/heat1q/boardsite/redis/redisfakes"

	"github.com/heat1q/boardsite/session"
	"github.com/heat1q/boardsite/session/sessionfakes"
	"github.com/stretchr/testify/assert"
)

func Test_controlBlock_AddPages(t *testing.T) {
	ctx := context.Background()
	sessionId := "sid1"
	fakeDispatcher := &sessionfakes.FakeDispatcher{}
	fakeBroadcaster := &sessionfakes.FakeBroadcaster{}
	fakeCache := &redisfakes.FakeHandler{}
	fakeAttachments := &attachmentfakes.FakeHandler{}

	scb, err := session.NewControlBlock(sessionId, session.WithCache(fakeCache), session.WithAttachments(fakeAttachments),
		session.WithDispatcher(fakeDispatcher), session.WithBroadcaster(fakeBroadcaster))
	assert.NoError(t, err)

	broadcast := make(chan types.Message, 1)
	defer close(broadcast)
	fakeBroadcaster.BroadcastReturns(broadcast)

	t.Run("successful", func(t *testing.T) {
		meta := &session.PageMeta{PageSize: session.PageSize{768, 1024}, Background: session.PageBackground{Style: "ruled"}}
		pageRequest := session.PageRequest{
			PageID: []string{"pid1"},
			Index:  []int{-1},
			Meta:   map[string]*session.PageMeta{"pid1": meta},
		}
		fakeCache.AddPageCalls(func(_ context.Context, sid string, pid string, index int, meta interface{}) error {
			assert.Equal(t, sessionId, sid)
			assert.Equal(t, "pid1", pid)
			assert.Equal(t, -1, index)
			assert.Equal(t, pageRequest.Meta["pid1"], meta)
			return nil
		})

		err := scb.AddPages(ctx, pageRequest)

		assert.NoError(t, err)
	})
}

func Test_controlBlock_GetPageSync(t *testing.T) {
	ctx := context.Background()
	sessionId := "sid1"
	fakeDispatcher := &sessionfakes.FakeDispatcher{}
	fakeBroadcaster := &sessionfakes.FakeBroadcaster{}
	fakeCache := &redisfakes.FakeHandler{}
	fakeAttachments := &attachmentfakes.FakeHandler{}

	scb, err := session.NewControlBlock(sessionId, session.WithCache(fakeCache), session.WithAttachments(fakeAttachments),
		session.WithDispatcher(fakeDispatcher), session.WithBroadcaster(fakeBroadcaster))
	assert.NoError(t, err)

	t.Run("without strokes", func(t *testing.T) {
		want := session.PageSync{
			PageRank: []string{"pid1", "pid2"},
			Pages: map[string]*session.Page{
				"pid1": {
					PageId:  "pid1",
					Meta:    &session.PageMeta{PageSize: session.PageSize{768, 1024}, Background: session.PageBackground{Style: "ruled"}},
					Strokes: &[]*session.Stroke{{ID: "stroke1"}},
				},
				"pid2": {
					PageId:  "pid2",
					Meta:    &session.PageMeta{PageSize: session.PageSize{768, 1024}, Background: session.PageBackground{Style: "ruled"}},
					Strokes: &[]*session.Stroke{{ID: "stroke2"}},
				},
			},
		}
		fakeCache.GetPageRankReturns(want.PageRank, nil)
		calls := 0
		fakeCache.GetPageMetaCalls(func(_ context.Context, sid string, pid string, i interface{}) error {
			defer func() { calls++ }()
			meta := i.(*session.PageMeta)
			if calls == 0 {
				meta.PageSize = session.PageSize{768, 1024}
				meta.Background.Style = "ruled"
				return nil
			}
			if calls == 1 {
				meta.PageSize = session.PageSize{768, 1024}
				meta.Background.Style = "ruled"
				return nil
			}
			return assert.AnError
		})
		strokepid1, _ := json.Marshal((*want.Pages["pid1"].Strokes)[0])
		fakeCache.GetPageStrokesReturnsOnCall(0, [][]byte{strokepid1}, nil)
		strokepid2, _ := json.Marshal((*want.Pages["pid2"].Strokes)[0])
		fakeCache.GetPageStrokesReturnsOnCall(1, [][]byte{strokepid2}, nil)

		got, err := scb.GetPageSync(ctx, []string{"pid1", "pid2"}, true)

		assert.NoError(t, err)
		assert.Equal(t, &want, got)
	})

}

func Test_controlBlock_UpdatePages(t *testing.T) {
	ctx := context.Background()
	sessionId := "sid1"
	fakeDispatcher := &sessionfakes.FakeDispatcher{}
	fakeBroadcaster := &sessionfakes.FakeBroadcaster{}
	fakeAttachments := &attachmentfakes.FakeHandler{}

	broadcast := make(chan types.Message, 10)
	defer close(broadcast)
	fakeBroadcaster.BroadcastReturns(broadcast)

	t.Run("update pages meta", func(t *testing.T) {
		fakeCache := &redisfakes.FakeHandler{}
		scb, err := session.NewControlBlock(sessionId, session.WithCache(fakeCache), session.WithAttachments(fakeAttachments),
			session.WithDispatcher(fakeDispatcher), session.WithBroadcaster(fakeBroadcaster))
		assert.NoError(t, err)

		pageRequest := session.PageRequest{
			Meta: map[string]*session.PageMeta{
				"pid1": {PageSize: session.PageSize{1234, 5678}},
			},
		}
		want := &session.PageMeta{PageSize: session.PageSize{1234, 5678}, Background: session.PageBackground{Style: "ruled"}}

		fakeCache.GetPageRankReturns([]string{"pid1"}, nil)
		calls := 0
		fakeCache.GetPageMetaCalls(func(_ context.Context, _ string, _ string, i interface{}) error {
			calls++
			meta := i.(*session.PageMeta)
			meta.PageSize = session.PageSize{768, 1024}
			meta.Background.Style = "ruled"
			return nil
		})

		fakeCache.SetPageMetaCalls(func(_ context.Context, sid string, pid string, meta interface{}) error {
			assert.Equal(t, sessionId, sid)
			assert.Equal(t, "pid1", pid)
			assert.Equal(t, *want, meta)
			return nil
		})

		err = scb.UpdatePages(ctx, pageRequest, "meta")

		assert.NoError(t, err)
		assert.Equal(t, 2, calls)
	})

	t.Run("deletes pages", func(t *testing.T) {
		fakeCache := &redisfakes.FakeHandler{}
		scb, err := session.NewControlBlock(sessionId, session.WithCache(fakeCache), session.WithAttachments(fakeAttachments),
			session.WithDispatcher(fakeDispatcher), session.WithBroadcaster(fakeBroadcaster))
		assert.NoError(t, err)

		pageRequest := session.PageRequest{
			PageID: []string{"pid1", "pid2"},
		}
		fakeCache.GetPageRankReturnsOnCall(0, []string{"pid1", "pid2"}, nil)
		fakeCache.GetPageRankReturnsOnCall(1, []string{"pid1", "pid2"}, nil)
		fakeCache.GetPageRankReturnsOnCall(2, []string{}, nil)

		err = scb.UpdatePages(ctx, pageRequest, "delete")

		assert.NoError(t, err)

		_, sid, pid := fakeCache.DeletePageArgsForCall(0)
		assert.Equal(t, sessionId, sid)
		assert.Equal(t, "pid1", pid)

		_, sid, pid = fakeCache.DeletePageArgsForCall(1)
		assert.Equal(t, sessionId, sid)
		assert.Equal(t, "pid2", pid)

		assert.Equal(t, 0, fakeCache.GetPageMetaCallCount())
	})

	t.Run("clears pages", func(t *testing.T) {
		fakeCache := &redisfakes.FakeHandler{}
		scb, err := session.NewControlBlock(sessionId, session.WithCache(fakeCache), session.WithAttachments(fakeAttachments),
			session.WithDispatcher(fakeDispatcher), session.WithBroadcaster(fakeBroadcaster))
		assert.NoError(t, err)

		pageRequest := session.PageRequest{
			PageID: []string{"pid1", "pid2"},
		}
		fakeCache.GetPageRankReturns([]string{"pid1", "pid2"}, nil)

		err = scb.UpdatePages(ctx, pageRequest, "clear")

		assert.NoError(t, err)

		_, sid, pid := fakeCache.ClearPageArgsForCall(0)
		assert.Equal(t, sessionId, sid)
		assert.Equal(t, "pid1", pid)

		_, sid, pid = fakeCache.ClearPageArgsForCall(1)
		assert.Equal(t, sessionId, sid)
		assert.Equal(t, "pid2", pid)

		assert.Equal(t, 2, fakeCache.GetPageMetaCallCount())
	})

	t.Run("unknown operation", func(t *testing.T) {
		fakeCache := &redisfakes.FakeHandler{}
		scb, err := session.NewControlBlock(sessionId, session.WithCache(fakeCache), session.WithAttachments(fakeAttachments),
			session.WithDispatcher(fakeDispatcher), session.WithBroadcaster(fakeBroadcaster))
		assert.NoError(t, err)

		err = scb.UpdatePages(ctx, session.PageRequest{}, "test")
		assert.Error(t, err)
	})
}
