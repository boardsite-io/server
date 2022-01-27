// Code generated by counterfeiter. DO NOT EDIT.
package sessionfakes

import (
	"context"
	"sync"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/attachment"
	"github.com/heat1q/boardsite/session"
)

type FakeController struct {
	AddPagesStub        func(context.Context, session.PageRequest) error
	addPagesMutex       sync.RWMutex
	addPagesArgsForCall []struct {
		arg1 context.Context
		arg2 session.PageRequest
	}
	addPagesReturns struct {
		result1 error
	}
	addPagesReturnsOnCall map[int]struct {
		result1 error
	}
	AttachmentsStub        func() attachment.Handler
	attachmentsMutex       sync.RWMutex
	attachmentsArgsForCall []struct {
	}
	attachmentsReturns struct {
		result1 attachment.Handler
	}
	attachmentsReturnsOnCall map[int]struct {
		result1 attachment.Handler
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	GetPageStub        func(context.Context, string, bool) (*session.Page, error)
	getPageMutex       sync.RWMutex
	getPageArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 bool
	}
	getPageReturns struct {
		result1 *session.Page
		result2 error
	}
	getPageReturnsOnCall map[int]struct {
		result1 *session.Page
		result2 error
	}
	GetStrokesStub        func(context.Context, string) ([]*session.Stroke, error)
	getStrokesMutex       sync.RWMutex
	getStrokesArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	getStrokesReturns struct {
		result1 []*session.Stroke
		result2 error
	}
	getStrokesReturnsOnCall map[int]struct {
		result1 []*session.Stroke
		result2 error
	}
	GetUserReadyStub        func(string) (*session.User, error)
	getUserReadyMutex       sync.RWMutex
	getUserReadyArgsForCall []struct {
		arg1 string
	}
	getUserReadyReturns struct {
		result1 *session.User
		result2 error
	}
	getUserReadyReturnsOnCall map[int]struct {
		result1 *session.User
		result2 error
	}
	GetUsersStub        func() map[string]*session.User
	getUsersMutex       sync.RWMutex
	getUsersArgsForCall []struct {
	}
	getUsersReturns struct {
		result1 map[string]*session.User
	}
	getUsersReturnsOnCall map[int]struct {
		result1 map[string]*session.User
	}
	IDStub        func() string
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
	}
	iDReturns struct {
		result1 string
	}
	iDReturnsOnCall map[int]struct {
		result1 string
	}
	IsUserConnectedStub        func(string) bool
	isUserConnectedMutex       sync.RWMutex
	isUserConnectedArgsForCall []struct {
		arg1 string
	}
	isUserConnectedReturns struct {
		result1 bool
	}
	isUserConnectedReturnsOnCall map[int]struct {
		result1 bool
	}
	IsUserReadyStub        func(string) bool
	isUserReadyMutex       sync.RWMutex
	isUserReadyArgsForCall []struct {
		arg1 string
	}
	isUserReadyReturns struct {
		result1 bool
	}
	isUserReadyReturnsOnCall map[int]struct {
		result1 bool
	}
	IsValidPageStub        func(context.Context, ...string) bool
	isValidPageMutex       sync.RWMutex
	isValidPageArgsForCall []struct {
		arg1 context.Context
		arg2 []string
	}
	isValidPageReturns struct {
		result1 bool
	}
	isValidPageReturnsOnCall map[int]struct {
		result1 bool
	}
	NewUserStub        func(string, string) (*session.User, error)
	newUserMutex       sync.RWMutex
	newUserArgsForCall []struct {
		arg1 string
		arg2 string
	}
	newUserReturns struct {
		result1 *session.User
		result2 error
	}
	newUserReturnsOnCall map[int]struct {
		result1 *session.User
		result2 error
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
	ReceiveStub        func(context.Context, *types.Message) error
	receiveMutex       sync.RWMutex
	receiveArgsForCall []struct {
		arg1 context.Context
		arg2 *types.Message
	}
	receiveReturns struct {
		result1 error
	}
	receiveReturnsOnCall map[int]struct {
		result1 error
	}
	UpdatePagesStub        func(context.Context, session.PageRequest, string) error
	updatePagesMutex       sync.RWMutex
	updatePagesArgsForCall []struct {
		arg1 context.Context
		arg2 session.PageRequest
		arg3 string
	}
	updatePagesReturns struct {
		result1 error
	}
	updatePagesReturnsOnCall map[int]struct {
		result1 error
	}
	UserConnectStub        func(*session.User)
	userConnectMutex       sync.RWMutex
	userConnectArgsForCall []struct {
		arg1 *session.User
	}
	UserDisconnectStub        func(context.Context, string)
	userDisconnectMutex       sync.RWMutex
	userDisconnectArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	UserReadyStub        func(*session.User) error
	userReadyMutex       sync.RWMutex
	userReadyArgsForCall []struct {
		arg1 *session.User
	}
	userReadyReturns struct {
		result1 error
	}
	userReadyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeController) AddPages(arg1 context.Context, arg2 session.PageRequest) error {
	fake.addPagesMutex.Lock()
	ret, specificReturn := fake.addPagesReturnsOnCall[len(fake.addPagesArgsForCall)]
	fake.addPagesArgsForCall = append(fake.addPagesArgsForCall, struct {
		arg1 context.Context
		arg2 session.PageRequest
	}{arg1, arg2})
	stub := fake.AddPagesStub
	fakeReturns := fake.addPagesReturns
	fake.recordInvocation("AddPages", []interface{}{arg1, arg2})
	fake.addPagesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) AddPagesCallCount() int {
	fake.addPagesMutex.RLock()
	defer fake.addPagesMutex.RUnlock()
	return len(fake.addPagesArgsForCall)
}

func (fake *FakeController) AddPagesCalls(stub func(context.Context, session.PageRequest) error) {
	fake.addPagesMutex.Lock()
	defer fake.addPagesMutex.Unlock()
	fake.AddPagesStub = stub
}

func (fake *FakeController) AddPagesArgsForCall(i int) (context.Context, session.PageRequest) {
	fake.addPagesMutex.RLock()
	defer fake.addPagesMutex.RUnlock()
	argsForCall := fake.addPagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) AddPagesReturns(result1 error) {
	fake.addPagesMutex.Lock()
	defer fake.addPagesMutex.Unlock()
	fake.AddPagesStub = nil
	fake.addPagesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddPagesReturnsOnCall(i int, result1 error) {
	fake.addPagesMutex.Lock()
	defer fake.addPagesMutex.Unlock()
	fake.AddPagesStub = nil
	if fake.addPagesReturnsOnCall == nil {
		fake.addPagesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addPagesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) Attachments() attachment.Handler {
	fake.attachmentsMutex.Lock()
	ret, specificReturn := fake.attachmentsReturnsOnCall[len(fake.attachmentsArgsForCall)]
	fake.attachmentsArgsForCall = append(fake.attachmentsArgsForCall, struct {
	}{})
	stub := fake.AttachmentsStub
	fakeReturns := fake.attachmentsReturns
	fake.recordInvocation("Attachments", []interface{}{})
	fake.attachmentsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) AttachmentsCallCount() int {
	fake.attachmentsMutex.RLock()
	defer fake.attachmentsMutex.RUnlock()
	return len(fake.attachmentsArgsForCall)
}

func (fake *FakeController) AttachmentsCalls(stub func() attachment.Handler) {
	fake.attachmentsMutex.Lock()
	defer fake.attachmentsMutex.Unlock()
	fake.AttachmentsStub = stub
}

func (fake *FakeController) AttachmentsReturns(result1 attachment.Handler) {
	fake.attachmentsMutex.Lock()
	defer fake.attachmentsMutex.Unlock()
	fake.AttachmentsStub = nil
	fake.attachmentsReturns = struct {
		result1 attachment.Handler
	}{result1}
}

func (fake *FakeController) AttachmentsReturnsOnCall(i int, result1 attachment.Handler) {
	fake.attachmentsMutex.Lock()
	defer fake.attachmentsMutex.Unlock()
	fake.AttachmentsStub = nil
	if fake.attachmentsReturnsOnCall == nil {
		fake.attachmentsReturnsOnCall = make(map[int]struct {
			result1 attachment.Handler
		})
	}
	fake.attachmentsReturnsOnCall[i] = struct {
		result1 attachment.Handler
	}{result1}
}

func (fake *FakeController) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		fake.CloseStub()
	}
}

func (fake *FakeController) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeController) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeController) GetPage(arg1 context.Context, arg2 string, arg3 bool) (*session.Page, error) {
	fake.getPageMutex.Lock()
	ret, specificReturn := fake.getPageReturnsOnCall[len(fake.getPageArgsForCall)]
	fake.getPageArgsForCall = append(fake.getPageArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 bool
	}{arg1, arg2, arg3})
	stub := fake.GetPageStub
	fakeReturns := fake.getPageReturns
	fake.recordInvocation("GetPage", []interface{}{arg1, arg2, arg3})
	fake.getPageMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeController) GetPageCallCount() int {
	fake.getPageMutex.RLock()
	defer fake.getPageMutex.RUnlock()
	return len(fake.getPageArgsForCall)
}

func (fake *FakeController) GetPageCalls(stub func(context.Context, string, bool) (*session.Page, error)) {
	fake.getPageMutex.Lock()
	defer fake.getPageMutex.Unlock()
	fake.GetPageStub = stub
}

func (fake *FakeController) GetPageArgsForCall(i int) (context.Context, string, bool) {
	fake.getPageMutex.RLock()
	defer fake.getPageMutex.RUnlock()
	argsForCall := fake.getPageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeController) GetPageReturns(result1 *session.Page, result2 error) {
	fake.getPageMutex.Lock()
	defer fake.getPageMutex.Unlock()
	fake.GetPageStub = nil
	fake.getPageReturns = struct {
		result1 *session.Page
		result2 error
	}{result1, result2}
}

func (fake *FakeController) GetPageReturnsOnCall(i int, result1 *session.Page, result2 error) {
	fake.getPageMutex.Lock()
	defer fake.getPageMutex.Unlock()
	fake.GetPageStub = nil
	if fake.getPageReturnsOnCall == nil {
		fake.getPageReturnsOnCall = make(map[int]struct {
			result1 *session.Page
			result2 error
		})
	}
	fake.getPageReturnsOnCall[i] = struct {
		result1 *session.Page
		result2 error
	}{result1, result2}
}

func (fake *FakeController) GetStrokes(arg1 context.Context, arg2 string) ([]*session.Stroke, error) {
	fake.getStrokesMutex.Lock()
	ret, specificReturn := fake.getStrokesReturnsOnCall[len(fake.getStrokesArgsForCall)]
	fake.getStrokesArgsForCall = append(fake.getStrokesArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.GetStrokesStub
	fakeReturns := fake.getStrokesReturns
	fake.recordInvocation("GetStrokes", []interface{}{arg1, arg2})
	fake.getStrokesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeController) GetStrokesCallCount() int {
	fake.getStrokesMutex.RLock()
	defer fake.getStrokesMutex.RUnlock()
	return len(fake.getStrokesArgsForCall)
}

func (fake *FakeController) GetStrokesCalls(stub func(context.Context, string) ([]*session.Stroke, error)) {
	fake.getStrokesMutex.Lock()
	defer fake.getStrokesMutex.Unlock()
	fake.GetStrokesStub = stub
}

func (fake *FakeController) GetStrokesArgsForCall(i int) (context.Context, string) {
	fake.getStrokesMutex.RLock()
	defer fake.getStrokesMutex.RUnlock()
	argsForCall := fake.getStrokesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) GetStrokesReturns(result1 []*session.Stroke, result2 error) {
	fake.getStrokesMutex.Lock()
	defer fake.getStrokesMutex.Unlock()
	fake.GetStrokesStub = nil
	fake.getStrokesReturns = struct {
		result1 []*session.Stroke
		result2 error
	}{result1, result2}
}

func (fake *FakeController) GetStrokesReturnsOnCall(i int, result1 []*session.Stroke, result2 error) {
	fake.getStrokesMutex.Lock()
	defer fake.getStrokesMutex.Unlock()
	fake.GetStrokesStub = nil
	if fake.getStrokesReturnsOnCall == nil {
		fake.getStrokesReturnsOnCall = make(map[int]struct {
			result1 []*session.Stroke
			result2 error
		})
	}
	fake.getStrokesReturnsOnCall[i] = struct {
		result1 []*session.Stroke
		result2 error
	}{result1, result2}
}

func (fake *FakeController) GetUserReady(arg1 string) (*session.User, error) {
	fake.getUserReadyMutex.Lock()
	ret, specificReturn := fake.getUserReadyReturnsOnCall[len(fake.getUserReadyArgsForCall)]
	fake.getUserReadyArgsForCall = append(fake.getUserReadyArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetUserReadyStub
	fakeReturns := fake.getUserReadyReturns
	fake.recordInvocation("GetUserReady", []interface{}{arg1})
	fake.getUserReadyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeController) GetUserReadyCallCount() int {
	fake.getUserReadyMutex.RLock()
	defer fake.getUserReadyMutex.RUnlock()
	return len(fake.getUserReadyArgsForCall)
}

func (fake *FakeController) GetUserReadyCalls(stub func(string) (*session.User, error)) {
	fake.getUserReadyMutex.Lock()
	defer fake.getUserReadyMutex.Unlock()
	fake.GetUserReadyStub = stub
}

func (fake *FakeController) GetUserReadyArgsForCall(i int) string {
	fake.getUserReadyMutex.RLock()
	defer fake.getUserReadyMutex.RUnlock()
	argsForCall := fake.getUserReadyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) GetUserReadyReturns(result1 *session.User, result2 error) {
	fake.getUserReadyMutex.Lock()
	defer fake.getUserReadyMutex.Unlock()
	fake.GetUserReadyStub = nil
	fake.getUserReadyReturns = struct {
		result1 *session.User
		result2 error
	}{result1, result2}
}

func (fake *FakeController) GetUserReadyReturnsOnCall(i int, result1 *session.User, result2 error) {
	fake.getUserReadyMutex.Lock()
	defer fake.getUserReadyMutex.Unlock()
	fake.GetUserReadyStub = nil
	if fake.getUserReadyReturnsOnCall == nil {
		fake.getUserReadyReturnsOnCall = make(map[int]struct {
			result1 *session.User
			result2 error
		})
	}
	fake.getUserReadyReturnsOnCall[i] = struct {
		result1 *session.User
		result2 error
	}{result1, result2}
}

func (fake *FakeController) GetUsers() map[string]*session.User {
	fake.getUsersMutex.Lock()
	ret, specificReturn := fake.getUsersReturnsOnCall[len(fake.getUsersArgsForCall)]
	fake.getUsersArgsForCall = append(fake.getUsersArgsForCall, struct {
	}{})
	stub := fake.GetUsersStub
	fakeReturns := fake.getUsersReturns
	fake.recordInvocation("GetUsers", []interface{}{})
	fake.getUsersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) GetUsersCallCount() int {
	fake.getUsersMutex.RLock()
	defer fake.getUsersMutex.RUnlock()
	return len(fake.getUsersArgsForCall)
}

func (fake *FakeController) GetUsersCalls(stub func() map[string]*session.User) {
	fake.getUsersMutex.Lock()
	defer fake.getUsersMutex.Unlock()
	fake.GetUsersStub = stub
}

func (fake *FakeController) GetUsersReturns(result1 map[string]*session.User) {
	fake.getUsersMutex.Lock()
	defer fake.getUsersMutex.Unlock()
	fake.GetUsersStub = nil
	fake.getUsersReturns = struct {
		result1 map[string]*session.User
	}{result1}
}

func (fake *FakeController) GetUsersReturnsOnCall(i int, result1 map[string]*session.User) {
	fake.getUsersMutex.Lock()
	defer fake.getUsersMutex.Unlock()
	fake.GetUsersStub = nil
	if fake.getUsersReturnsOnCall == nil {
		fake.getUsersReturnsOnCall = make(map[int]struct {
			result1 map[string]*session.User
		})
	}
	fake.getUsersReturnsOnCall[i] = struct {
		result1 map[string]*session.User
	}{result1}
}

func (fake *FakeController) ID() string {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct {
	}{})
	stub := fake.IDStub
	fakeReturns := fake.iDReturns
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeController) IDCalls(stub func() string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *FakeController) IDReturns(result1 string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeController) IDReturnsOnCall(i int, result1 string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeController) IsUserConnected(arg1 string) bool {
	fake.isUserConnectedMutex.Lock()
	ret, specificReturn := fake.isUserConnectedReturnsOnCall[len(fake.isUserConnectedArgsForCall)]
	fake.isUserConnectedArgsForCall = append(fake.isUserConnectedArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.IsUserConnectedStub
	fakeReturns := fake.isUserConnectedReturns
	fake.recordInvocation("IsUserConnected", []interface{}{arg1})
	fake.isUserConnectedMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) IsUserConnectedCallCount() int {
	fake.isUserConnectedMutex.RLock()
	defer fake.isUserConnectedMutex.RUnlock()
	return len(fake.isUserConnectedArgsForCall)
}

func (fake *FakeController) IsUserConnectedCalls(stub func(string) bool) {
	fake.isUserConnectedMutex.Lock()
	defer fake.isUserConnectedMutex.Unlock()
	fake.IsUserConnectedStub = stub
}

func (fake *FakeController) IsUserConnectedArgsForCall(i int) string {
	fake.isUserConnectedMutex.RLock()
	defer fake.isUserConnectedMutex.RUnlock()
	argsForCall := fake.isUserConnectedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) IsUserConnectedReturns(result1 bool) {
	fake.isUserConnectedMutex.Lock()
	defer fake.isUserConnectedMutex.Unlock()
	fake.IsUserConnectedStub = nil
	fake.isUserConnectedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) IsUserConnectedReturnsOnCall(i int, result1 bool) {
	fake.isUserConnectedMutex.Lock()
	defer fake.isUserConnectedMutex.Unlock()
	fake.IsUserConnectedStub = nil
	if fake.isUserConnectedReturnsOnCall == nil {
		fake.isUserConnectedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isUserConnectedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) IsUserReady(arg1 string) bool {
	fake.isUserReadyMutex.Lock()
	ret, specificReturn := fake.isUserReadyReturnsOnCall[len(fake.isUserReadyArgsForCall)]
	fake.isUserReadyArgsForCall = append(fake.isUserReadyArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.IsUserReadyStub
	fakeReturns := fake.isUserReadyReturns
	fake.recordInvocation("IsUserReady", []interface{}{arg1})
	fake.isUserReadyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) IsUserReadyCallCount() int {
	fake.isUserReadyMutex.RLock()
	defer fake.isUserReadyMutex.RUnlock()
	return len(fake.isUserReadyArgsForCall)
}

func (fake *FakeController) IsUserReadyCalls(stub func(string) bool) {
	fake.isUserReadyMutex.Lock()
	defer fake.isUserReadyMutex.Unlock()
	fake.IsUserReadyStub = stub
}

func (fake *FakeController) IsUserReadyArgsForCall(i int) string {
	fake.isUserReadyMutex.RLock()
	defer fake.isUserReadyMutex.RUnlock()
	argsForCall := fake.isUserReadyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) IsUserReadyReturns(result1 bool) {
	fake.isUserReadyMutex.Lock()
	defer fake.isUserReadyMutex.Unlock()
	fake.IsUserReadyStub = nil
	fake.isUserReadyReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) IsUserReadyReturnsOnCall(i int, result1 bool) {
	fake.isUserReadyMutex.Lock()
	defer fake.isUserReadyMutex.Unlock()
	fake.IsUserReadyStub = nil
	if fake.isUserReadyReturnsOnCall == nil {
		fake.isUserReadyReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isUserReadyReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) IsValidPage(arg1 context.Context, arg2 ...string) bool {
	fake.isValidPageMutex.Lock()
	ret, specificReturn := fake.isValidPageReturnsOnCall[len(fake.isValidPageArgsForCall)]
	fake.isValidPageArgsForCall = append(fake.isValidPageArgsForCall, struct {
		arg1 context.Context
		arg2 []string
	}{arg1, arg2})
	stub := fake.IsValidPageStub
	fakeReturns := fake.isValidPageReturns
	fake.recordInvocation("IsValidPage", []interface{}{arg1, arg2})
	fake.isValidPageMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) IsValidPageCallCount() int {
	fake.isValidPageMutex.RLock()
	defer fake.isValidPageMutex.RUnlock()
	return len(fake.isValidPageArgsForCall)
}

func (fake *FakeController) IsValidPageCalls(stub func(context.Context, ...string) bool) {
	fake.isValidPageMutex.Lock()
	defer fake.isValidPageMutex.Unlock()
	fake.IsValidPageStub = stub
}

func (fake *FakeController) IsValidPageArgsForCall(i int) (context.Context, []string) {
	fake.isValidPageMutex.RLock()
	defer fake.isValidPageMutex.RUnlock()
	argsForCall := fake.isValidPageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) IsValidPageReturns(result1 bool) {
	fake.isValidPageMutex.Lock()
	defer fake.isValidPageMutex.Unlock()
	fake.IsValidPageStub = nil
	fake.isValidPageReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) IsValidPageReturnsOnCall(i int, result1 bool) {
	fake.isValidPageMutex.Lock()
	defer fake.isValidPageMutex.Unlock()
	fake.IsValidPageStub = nil
	if fake.isValidPageReturnsOnCall == nil {
		fake.isValidPageReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isValidPageReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) NewUser(arg1 string, arg2 string) (*session.User, error) {
	fake.newUserMutex.Lock()
	ret, specificReturn := fake.newUserReturnsOnCall[len(fake.newUserArgsForCall)]
	fake.newUserArgsForCall = append(fake.newUserArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.NewUserStub
	fakeReturns := fake.newUserReturns
	fake.recordInvocation("NewUser", []interface{}{arg1, arg2})
	fake.newUserMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeController) NewUserCallCount() int {
	fake.newUserMutex.RLock()
	defer fake.newUserMutex.RUnlock()
	return len(fake.newUserArgsForCall)
}

func (fake *FakeController) NewUserCalls(stub func(string, string) (*session.User, error)) {
	fake.newUserMutex.Lock()
	defer fake.newUserMutex.Unlock()
	fake.NewUserStub = stub
}

func (fake *FakeController) NewUserArgsForCall(i int) (string, string) {
	fake.newUserMutex.RLock()
	defer fake.newUserMutex.RUnlock()
	argsForCall := fake.newUserArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) NewUserReturns(result1 *session.User, result2 error) {
	fake.newUserMutex.Lock()
	defer fake.newUserMutex.Unlock()
	fake.NewUserStub = nil
	fake.newUserReturns = struct {
		result1 *session.User
		result2 error
	}{result1, result2}
}

func (fake *FakeController) NewUserReturnsOnCall(i int, result1 *session.User, result2 error) {
	fake.newUserMutex.Lock()
	defer fake.newUserMutex.Unlock()
	fake.NewUserStub = nil
	if fake.newUserReturnsOnCall == nil {
		fake.newUserReturnsOnCall = make(map[int]struct {
			result1 *session.User
			result2 error
		})
	}
	fake.newUserReturnsOnCall[i] = struct {
		result1 *session.User
		result2 error
	}{result1, result2}
}

func (fake *FakeController) NumUsers() int {
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

func (fake *FakeController) NumUsersCallCount() int {
	fake.numUsersMutex.RLock()
	defer fake.numUsersMutex.RUnlock()
	return len(fake.numUsersArgsForCall)
}

func (fake *FakeController) NumUsersCalls(stub func() int) {
	fake.numUsersMutex.Lock()
	defer fake.numUsersMutex.Unlock()
	fake.NumUsersStub = stub
}

func (fake *FakeController) NumUsersReturns(result1 int) {
	fake.numUsersMutex.Lock()
	defer fake.numUsersMutex.Unlock()
	fake.NumUsersStub = nil
	fake.numUsersReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeController) NumUsersReturnsOnCall(i int, result1 int) {
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

func (fake *FakeController) Receive(arg1 context.Context, arg2 *types.Message) error {
	fake.receiveMutex.Lock()
	ret, specificReturn := fake.receiveReturnsOnCall[len(fake.receiveArgsForCall)]
	fake.receiveArgsForCall = append(fake.receiveArgsForCall, struct {
		arg1 context.Context
		arg2 *types.Message
	}{arg1, arg2})
	stub := fake.ReceiveStub
	fakeReturns := fake.receiveReturns
	fake.recordInvocation("Receive", []interface{}{arg1, arg2})
	fake.receiveMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) ReceiveCallCount() int {
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	return len(fake.receiveArgsForCall)
}

func (fake *FakeController) ReceiveCalls(stub func(context.Context, *types.Message) error) {
	fake.receiveMutex.Lock()
	defer fake.receiveMutex.Unlock()
	fake.ReceiveStub = stub
}

func (fake *FakeController) ReceiveArgsForCall(i int) (context.Context, *types.Message) {
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	argsForCall := fake.receiveArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) ReceiveReturns(result1 error) {
	fake.receiveMutex.Lock()
	defer fake.receiveMutex.Unlock()
	fake.ReceiveStub = nil
	fake.receiveReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) ReceiveReturnsOnCall(i int, result1 error) {
	fake.receiveMutex.Lock()
	defer fake.receiveMutex.Unlock()
	fake.ReceiveStub = nil
	if fake.receiveReturnsOnCall == nil {
		fake.receiveReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.receiveReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) UpdatePages(arg1 context.Context, arg2 session.PageRequest, arg3 string) error {
	fake.updatePagesMutex.Lock()
	ret, specificReturn := fake.updatePagesReturnsOnCall[len(fake.updatePagesArgsForCall)]
	fake.updatePagesArgsForCall = append(fake.updatePagesArgsForCall, struct {
		arg1 context.Context
		arg2 session.PageRequest
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.UpdatePagesStub
	fakeReturns := fake.updatePagesReturns
	fake.recordInvocation("UpdatePages", []interface{}{arg1, arg2, arg3})
	fake.updatePagesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) UpdatePagesCallCount() int {
	fake.updatePagesMutex.RLock()
	defer fake.updatePagesMutex.RUnlock()
	return len(fake.updatePagesArgsForCall)
}

func (fake *FakeController) UpdatePagesCalls(stub func(context.Context, session.PageRequest, string) error) {
	fake.updatePagesMutex.Lock()
	defer fake.updatePagesMutex.Unlock()
	fake.UpdatePagesStub = stub
}

func (fake *FakeController) UpdatePagesArgsForCall(i int) (context.Context, session.PageRequest, string) {
	fake.updatePagesMutex.RLock()
	defer fake.updatePagesMutex.RUnlock()
	argsForCall := fake.updatePagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeController) UpdatePagesReturns(result1 error) {
	fake.updatePagesMutex.Lock()
	defer fake.updatePagesMutex.Unlock()
	fake.UpdatePagesStub = nil
	fake.updatePagesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) UpdatePagesReturnsOnCall(i int, result1 error) {
	fake.updatePagesMutex.Lock()
	defer fake.updatePagesMutex.Unlock()
	fake.UpdatePagesStub = nil
	if fake.updatePagesReturnsOnCall == nil {
		fake.updatePagesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updatePagesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) UserConnect(arg1 *session.User) {
	fake.userConnectMutex.Lock()
	fake.userConnectArgsForCall = append(fake.userConnectArgsForCall, struct {
		arg1 *session.User
	}{arg1})
	stub := fake.UserConnectStub
	fake.recordInvocation("UserConnect", []interface{}{arg1})
	fake.userConnectMutex.Unlock()
	if stub != nil {
		fake.UserConnectStub(arg1)
	}
}

func (fake *FakeController) UserConnectCallCount() int {
	fake.userConnectMutex.RLock()
	defer fake.userConnectMutex.RUnlock()
	return len(fake.userConnectArgsForCall)
}

func (fake *FakeController) UserConnectCalls(stub func(*session.User)) {
	fake.userConnectMutex.Lock()
	defer fake.userConnectMutex.Unlock()
	fake.UserConnectStub = stub
}

func (fake *FakeController) UserConnectArgsForCall(i int) *session.User {
	fake.userConnectMutex.RLock()
	defer fake.userConnectMutex.RUnlock()
	argsForCall := fake.userConnectArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) UserDisconnect(arg1 context.Context, arg2 string) {
	fake.userDisconnectMutex.Lock()
	fake.userDisconnectArgsForCall = append(fake.userDisconnectArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.UserDisconnectStub
	fake.recordInvocation("UserDisconnect", []interface{}{arg1, arg2})
	fake.userDisconnectMutex.Unlock()
	if stub != nil {
		fake.UserDisconnectStub(arg1, arg2)
	}
}

func (fake *FakeController) UserDisconnectCallCount() int {
	fake.userDisconnectMutex.RLock()
	defer fake.userDisconnectMutex.RUnlock()
	return len(fake.userDisconnectArgsForCall)
}

func (fake *FakeController) UserDisconnectCalls(stub func(context.Context, string)) {
	fake.userDisconnectMutex.Lock()
	defer fake.userDisconnectMutex.Unlock()
	fake.UserDisconnectStub = stub
}

func (fake *FakeController) UserDisconnectArgsForCall(i int) (context.Context, string) {
	fake.userDisconnectMutex.RLock()
	defer fake.userDisconnectMutex.RUnlock()
	argsForCall := fake.userDisconnectArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) UserReady(arg1 *session.User) error {
	fake.userReadyMutex.Lock()
	ret, specificReturn := fake.userReadyReturnsOnCall[len(fake.userReadyArgsForCall)]
	fake.userReadyArgsForCall = append(fake.userReadyArgsForCall, struct {
		arg1 *session.User
	}{arg1})
	stub := fake.UserReadyStub
	fakeReturns := fake.userReadyReturns
	fake.recordInvocation("UserReady", []interface{}{arg1})
	fake.userReadyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeController) UserReadyCallCount() int {
	fake.userReadyMutex.RLock()
	defer fake.userReadyMutex.RUnlock()
	return len(fake.userReadyArgsForCall)
}

func (fake *FakeController) UserReadyCalls(stub func(*session.User) error) {
	fake.userReadyMutex.Lock()
	defer fake.userReadyMutex.Unlock()
	fake.UserReadyStub = stub
}

func (fake *FakeController) UserReadyArgsForCall(i int) *session.User {
	fake.userReadyMutex.RLock()
	defer fake.userReadyMutex.RUnlock()
	argsForCall := fake.userReadyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) UserReadyReturns(result1 error) {
	fake.userReadyMutex.Lock()
	defer fake.userReadyMutex.Unlock()
	fake.UserReadyStub = nil
	fake.userReadyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) UserReadyReturnsOnCall(i int, result1 error) {
	fake.userReadyMutex.Lock()
	defer fake.userReadyMutex.Unlock()
	fake.UserReadyStub = nil
	if fake.userReadyReturnsOnCall == nil {
		fake.userReadyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.userReadyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addPagesMutex.RLock()
	defer fake.addPagesMutex.RUnlock()
	fake.attachmentsMutex.RLock()
	defer fake.attachmentsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.getPageMutex.RLock()
	defer fake.getPageMutex.RUnlock()
	fake.getStrokesMutex.RLock()
	defer fake.getStrokesMutex.RUnlock()
	fake.getUserReadyMutex.RLock()
	defer fake.getUserReadyMutex.RUnlock()
	fake.getUsersMutex.RLock()
	defer fake.getUsersMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.isUserConnectedMutex.RLock()
	defer fake.isUserConnectedMutex.RUnlock()
	fake.isUserReadyMutex.RLock()
	defer fake.isUserReadyMutex.RUnlock()
	fake.isValidPageMutex.RLock()
	defer fake.isValidPageMutex.RUnlock()
	fake.newUserMutex.RLock()
	defer fake.newUserMutex.RUnlock()
	fake.numUsersMutex.RLock()
	defer fake.numUsersMutex.RUnlock()
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	fake.updatePagesMutex.RLock()
	defer fake.updatePagesMutex.RUnlock()
	fake.userConnectMutex.RLock()
	defer fake.userConnectMutex.RUnlock()
	fake.userDisconnectMutex.RLock()
	defer fake.userDisconnectMutex.RUnlock()
	fake.userReadyMutex.RLock()
	defer fake.userReadyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeController) recordInvocation(key string, args []interface{}) {
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

var _ session.Controller = new(FakeController)
