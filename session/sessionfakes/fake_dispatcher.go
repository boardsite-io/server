// Code generated by counterfeiter. DO NOT EDIT.
package sessionfakes

import (
	"context"
	"sync"

	"github.com/heat1q/boardsite/session"
)

type FakeDispatcher struct {
	CloseStub        func(context.Context, string) error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	CreateStub        func(context.Context, session.CreateConfig) (string, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 context.Context
		arg2 session.CreateConfig
	}
	createReturns struct {
		result1 string
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetSCBStub        func(string) (session.Controller, error)
	getSCBMutex       sync.RWMutex
	getSCBArgsForCall []struct {
		arg1 string
	}
	getSCBReturns struct {
		result1 session.Controller
		result2 error
	}
	getSCBReturnsOnCall map[int]struct {
		result1 session.Controller
		result2 error
	}
	IsValidStub        func(string) bool
	isValidMutex       sync.RWMutex
	isValidArgsForCall []struct {
		arg1 string
	}
	isValidReturns struct {
		result1 bool
	}
	isValidReturnsOnCall map[int]struct {
		result1 bool
	}
	NumSessionsStub        func() int
	numSessionsMutex       sync.RWMutex
	numSessionsArgsForCall []struct {
	}
	numSessionsReturns struct {
		result1 int
	}
	numSessionsReturnsOnCall map[int]struct {
		result1 int
	}
	NumUsersStub        func() int
	numUsersMutex       sync.RWMutex
	numUsersArgsForCall []struct {
	}
	numUsersReturns struct {
		result1 int
	}
	numUsersReturnsOnCall map[int]struct {
		result1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDispatcher) Close(arg1 context.Context, arg2 string) error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{arg1, arg2})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDispatcher) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeDispatcher) CloseCalls(stub func(context.Context, string) error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeDispatcher) CloseArgsForCall(i int) (context.Context, string) {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	argsForCall := fake.closeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDispatcher) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDispatcher) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeDispatcher) Create(arg1 context.Context, arg2 session.CreateConfig) (string, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 context.Context
		arg2 session.CreateConfig
	}{arg1, arg2})
	stub := fake.CreateStub
	fakeReturns := fake.createReturns
	fake.recordInvocation("Create", []interface{}{arg1, arg2})
	fake.createMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDispatcher) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeDispatcher) CreateCalls(stub func(context.Context, session.CreateConfig) (string, error)) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeDispatcher) CreateArgsForCall(i int) (context.Context, session.CreateConfig) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDispatcher) CreateReturns(result1 string, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeDispatcher) CreateReturnsOnCall(i int, result1 string, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeDispatcher) GetSCB(arg1 string) (session.Controller, error) {
	fake.getSCBMutex.Lock()
	ret, specificReturn := fake.getSCBReturnsOnCall[len(fake.getSCBArgsForCall)]
	fake.getSCBArgsForCall = append(fake.getSCBArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetSCBStub
	fakeReturns := fake.getSCBReturns
	fake.recordInvocation("GetSCB", []interface{}{arg1})
	fake.getSCBMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDispatcher) GetSCBCallCount() int {
	fake.getSCBMutex.RLock()
	defer fake.getSCBMutex.RUnlock()
	return len(fake.getSCBArgsForCall)
}

func (fake *FakeDispatcher) GetSCBCalls(stub func(string) (session.Controller, error)) {
	fake.getSCBMutex.Lock()
	defer fake.getSCBMutex.Unlock()
	fake.GetSCBStub = stub
}

func (fake *FakeDispatcher) GetSCBArgsForCall(i int) string {
	fake.getSCBMutex.RLock()
	defer fake.getSCBMutex.RUnlock()
	argsForCall := fake.getSCBArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDispatcher) GetSCBReturns(result1 session.Controller, result2 error) {
	fake.getSCBMutex.Lock()
	defer fake.getSCBMutex.Unlock()
	fake.GetSCBStub = nil
	fake.getSCBReturns = struct {
		result1 session.Controller
		result2 error
	}{result1, result2}
}

func (fake *FakeDispatcher) GetSCBReturnsOnCall(i int, result1 session.Controller, result2 error) {
	fake.getSCBMutex.Lock()
	defer fake.getSCBMutex.Unlock()
	fake.GetSCBStub = nil
	if fake.getSCBReturnsOnCall == nil {
		fake.getSCBReturnsOnCall = make(map[int]struct {
			result1 session.Controller
			result2 error
		})
	}
	fake.getSCBReturnsOnCall[i] = struct {
		result1 session.Controller
		result2 error
	}{result1, result2}
}

func (fake *FakeDispatcher) IsValid(arg1 string) bool {
	fake.isValidMutex.Lock()
	ret, specificReturn := fake.isValidReturnsOnCall[len(fake.isValidArgsForCall)]
	fake.isValidArgsForCall = append(fake.isValidArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.IsValidStub
	fakeReturns := fake.isValidReturns
	fake.recordInvocation("IsValid", []interface{}{arg1})
	fake.isValidMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDispatcher) IsValidCallCount() int {
	fake.isValidMutex.RLock()
	defer fake.isValidMutex.RUnlock()
	return len(fake.isValidArgsForCall)
}

func (fake *FakeDispatcher) IsValidCalls(stub func(string) bool) {
	fake.isValidMutex.Lock()
	defer fake.isValidMutex.Unlock()
	fake.IsValidStub = stub
}

func (fake *FakeDispatcher) IsValidArgsForCall(i int) string {
	fake.isValidMutex.RLock()
	defer fake.isValidMutex.RUnlock()
	argsForCall := fake.isValidArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDispatcher) IsValidReturns(result1 bool) {
	fake.isValidMutex.Lock()
	defer fake.isValidMutex.Unlock()
	fake.IsValidStub = nil
	fake.isValidReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeDispatcher) IsValidReturnsOnCall(i int, result1 bool) {
	fake.isValidMutex.Lock()
	defer fake.isValidMutex.Unlock()
	fake.IsValidStub = nil
	if fake.isValidReturnsOnCall == nil {
		fake.isValidReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isValidReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeDispatcher) NumSessions() int {
	fake.numSessionsMutex.Lock()
	ret, specificReturn := fake.numSessionsReturnsOnCall[len(fake.numSessionsArgsForCall)]
	fake.numSessionsArgsForCall = append(fake.numSessionsArgsForCall, struct {
	}{})
	stub := fake.NumSessionsStub
	fakeReturns := fake.numSessionsReturns
	fake.recordInvocation("NumSessions", []interface{}{})
	fake.numSessionsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDispatcher) NumSessionsCallCount() int {
	fake.numSessionsMutex.RLock()
	defer fake.numSessionsMutex.RUnlock()
	return len(fake.numSessionsArgsForCall)
}

func (fake *FakeDispatcher) NumSessionsCalls(stub func() int) {
	fake.numSessionsMutex.Lock()
	defer fake.numSessionsMutex.Unlock()
	fake.NumSessionsStub = stub
}

func (fake *FakeDispatcher) NumSessionsReturns(result1 int) {
	fake.numSessionsMutex.Lock()
	defer fake.numSessionsMutex.Unlock()
	fake.NumSessionsStub = nil
	fake.numSessionsReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeDispatcher) NumSessionsReturnsOnCall(i int, result1 int) {
	fake.numSessionsMutex.Lock()
	defer fake.numSessionsMutex.Unlock()
	fake.NumSessionsStub = nil
	if fake.numSessionsReturnsOnCall == nil {
		fake.numSessionsReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.numSessionsReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeDispatcher) NumUsers() int {
	fake.numUsersMutex.Lock()
	ret, specificReturn := fake.numUsersReturnsOnCall[len(fake.numUsersArgsForCall)]
	fake.numUsersArgsForCall = append(fake.numUsersArgsForCall, struct {
	}{})
	stub := fake.NumUsersStub
	fakeReturns := fake.numUsersReturns
	fake.recordInvocation("NumUsers", []interface{}{})
	fake.numUsersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDispatcher) NumUsersCallCount() int {
	fake.numUsersMutex.RLock()
	defer fake.numUsersMutex.RUnlock()
	return len(fake.numUsersArgsForCall)
}

func (fake *FakeDispatcher) NumUsersCalls(stub func() int) {
	fake.numUsersMutex.Lock()
	defer fake.numUsersMutex.Unlock()
	fake.NumUsersStub = stub
}

func (fake *FakeDispatcher) NumUsersReturns(result1 int) {
	fake.numUsersMutex.Lock()
	defer fake.numUsersMutex.Unlock()
	fake.NumUsersStub = nil
	fake.numUsersReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeDispatcher) NumUsersReturnsOnCall(i int, result1 int) {
	fake.numUsersMutex.Lock()
	defer fake.numUsersMutex.Unlock()
	fake.NumUsersStub = nil
	if fake.numUsersReturnsOnCall == nil {
		fake.numUsersReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.numUsersReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeDispatcher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.getSCBMutex.RLock()
	defer fake.getSCBMutex.RUnlock()
	fake.isValidMutex.RLock()
	defer fake.isValidMutex.RUnlock()
	fake.numSessionsMutex.RLock()
	defer fake.numSessionsMutex.RUnlock()
	fake.numUsersMutex.RLock()
	defer fake.numUsersMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDispatcher) recordInvocation(key string, args []interface{}) {
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

var _ session.Dispatcher = new(FakeDispatcher)
