// Code generated by counterfeiter. DO NOT EDIT.
package sessionfakes

import (
	"sync"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
	"github.com/heat1q/boardsite/session"
)

type FakeBroadcaster struct {
	BindStub        func(session.Controller) session.Broadcaster
	bindMutex       sync.RWMutex
	bindArgsForCall []struct {
		arg1 session.Controller
	}
	bindReturns struct {
		result1 session.Broadcaster
	}
	bindReturnsOnCall map[int]struct {
		result1 session.Broadcaster
	}
	BroadcastStub        func() chan<- types.Message
	broadcastMutex       sync.RWMutex
	broadcastArgsForCall []struct {
	}
	broadcastReturns struct {
		result1 chan<- types.Message
	}
	broadcastReturnsOnCall map[int]struct {
		result1 chan<- types.Message
	}
	CacheStub        func() chan<- []redis.Stroke
	cacheMutex       sync.RWMutex
	cacheArgsForCall []struct {
	}
	cacheReturns struct {
		result1 chan<- []redis.Stroke
	}
	cacheReturnsOnCall map[int]struct {
		result1 chan<- []redis.Stroke
	}
	CloseStub        func() chan<- struct{}
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 chan<- struct{}
	}
	closeReturnsOnCall map[int]struct {
		result1 chan<- struct{}
	}
	SendStub        func() chan<- types.Message
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
	}
	sendReturns struct {
		result1 chan<- types.Message
	}
	sendReturnsOnCall map[int]struct {
		result1 chan<- types.Message
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBroadcaster) Bind(arg1 session.Controller) session.Broadcaster {
	fake.bindMutex.Lock()
	ret, specificReturn := fake.bindReturnsOnCall[len(fake.bindArgsForCall)]
	fake.bindArgsForCall = append(fake.bindArgsForCall, struct {
		arg1 session.Controller
	}{arg1})
	stub := fake.BindStub
	fakeReturns := fake.bindReturns
	fake.recordInvocation("Bind", []interface{}{arg1})
	fake.bindMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBroadcaster) BindCallCount() int {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	return len(fake.bindArgsForCall)
}

func (fake *FakeBroadcaster) BindCalls(stub func(session.Controller) session.Broadcaster) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = stub
}

func (fake *FakeBroadcaster) BindArgsForCall(i int) session.Controller {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	argsForCall := fake.bindArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBroadcaster) BindReturns(result1 session.Broadcaster) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = nil
	fake.bindReturns = struct {
		result1 session.Broadcaster
	}{result1}
}

func (fake *FakeBroadcaster) BindReturnsOnCall(i int, result1 session.Broadcaster) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = nil
	if fake.bindReturnsOnCall == nil {
		fake.bindReturnsOnCall = make(map[int]struct {
			result1 session.Broadcaster
		})
	}
	fake.bindReturnsOnCall[i] = struct {
		result1 session.Broadcaster
	}{result1}
}

func (fake *FakeBroadcaster) Broadcast() chan<- types.Message {
	fake.broadcastMutex.Lock()
	ret, specificReturn := fake.broadcastReturnsOnCall[len(fake.broadcastArgsForCall)]
	fake.broadcastArgsForCall = append(fake.broadcastArgsForCall, struct {
	}{})
	stub := fake.BroadcastStub
	fakeReturns := fake.broadcastReturns
	fake.recordInvocation("Broadcast", []interface{}{})
	fake.broadcastMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBroadcaster) BroadcastCallCount() int {
	fake.broadcastMutex.RLock()
	defer fake.broadcastMutex.RUnlock()
	return len(fake.broadcastArgsForCall)
}

func (fake *FakeBroadcaster) BroadcastCalls(stub func() chan<- types.Message) {
	fake.broadcastMutex.Lock()
	defer fake.broadcastMutex.Unlock()
	fake.BroadcastStub = stub
}

func (fake *FakeBroadcaster) BroadcastReturns(result1 chan<- types.Message) {
	fake.broadcastMutex.Lock()
	defer fake.broadcastMutex.Unlock()
	fake.BroadcastStub = nil
	fake.broadcastReturns = struct {
		result1 chan<- types.Message
	}{result1}
}

func (fake *FakeBroadcaster) BroadcastReturnsOnCall(i int, result1 chan<- types.Message) {
	fake.broadcastMutex.Lock()
	defer fake.broadcastMutex.Unlock()
	fake.BroadcastStub = nil
	if fake.broadcastReturnsOnCall == nil {
		fake.broadcastReturnsOnCall = make(map[int]struct {
			result1 chan<- types.Message
		})
	}
	fake.broadcastReturnsOnCall[i] = struct {
		result1 chan<- types.Message
	}{result1}
}

func (fake *FakeBroadcaster) Cache() chan<- []redis.Stroke {
	fake.cacheMutex.Lock()
	ret, specificReturn := fake.cacheReturnsOnCall[len(fake.cacheArgsForCall)]
	fake.cacheArgsForCall = append(fake.cacheArgsForCall, struct {
	}{})
	stub := fake.CacheStub
	fakeReturns := fake.cacheReturns
	fake.recordInvocation("Cache", []interface{}{})
	fake.cacheMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBroadcaster) CacheCallCount() int {
	fake.cacheMutex.RLock()
	defer fake.cacheMutex.RUnlock()
	return len(fake.cacheArgsForCall)
}

func (fake *FakeBroadcaster) CacheCalls(stub func() chan<- []redis.Stroke) {
	fake.cacheMutex.Lock()
	defer fake.cacheMutex.Unlock()
	fake.CacheStub = stub
}

func (fake *FakeBroadcaster) CacheReturns(result1 chan<- []redis.Stroke) {
	fake.cacheMutex.Lock()
	defer fake.cacheMutex.Unlock()
	fake.CacheStub = nil
	fake.cacheReturns = struct {
		result1 chan<- []redis.Stroke
	}{result1}
}

func (fake *FakeBroadcaster) CacheReturnsOnCall(i int, result1 chan<- []redis.Stroke) {
	fake.cacheMutex.Lock()
	defer fake.cacheMutex.Unlock()
	fake.CacheStub = nil
	if fake.cacheReturnsOnCall == nil {
		fake.cacheReturnsOnCall = make(map[int]struct {
			result1 chan<- []redis.Stroke
		})
	}
	fake.cacheReturnsOnCall[i] = struct {
		result1 chan<- []redis.Stroke
	}{result1}
}

func (fake *FakeBroadcaster) Close() chan<- struct{} {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBroadcaster) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeBroadcaster) CloseCalls(stub func() chan<- struct{}) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeBroadcaster) CloseReturns(result1 chan<- struct{}) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 chan<- struct{}
	}{result1}
}

func (fake *FakeBroadcaster) CloseReturnsOnCall(i int, result1 chan<- struct{}) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 chan<- struct{}
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 chan<- struct{}
	}{result1}
}

func (fake *FakeBroadcaster) Send() chan<- types.Message {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
	}{})
	stub := fake.SendStub
	fakeReturns := fake.sendReturns
	fake.recordInvocation("Send", []interface{}{})
	fake.sendMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBroadcaster) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *FakeBroadcaster) SendCalls(stub func() chan<- types.Message) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *FakeBroadcaster) SendReturns(result1 chan<- types.Message) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 chan<- types.Message
	}{result1}
}

func (fake *FakeBroadcaster) SendReturnsOnCall(i int, result1 chan<- types.Message) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	if fake.sendReturnsOnCall == nil {
		fake.sendReturnsOnCall = make(map[int]struct {
			result1 chan<- types.Message
		})
	}
	fake.sendReturnsOnCall[i] = struct {
		result1 chan<- types.Message
	}{result1}
}

func (fake *FakeBroadcaster) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	fake.broadcastMutex.RLock()
	defer fake.broadcastMutex.RUnlock()
	fake.cacheMutex.RLock()
	defer fake.cacheMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBroadcaster) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ session.Broadcaster = new(FakeBroadcaster)
