// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package transport

import (
	"context"
	"github.com/ivch/dynasty/server/handlers/users"
	"sync"
)

var (
	lockUsersServiceMockAddFamilyMember    sync.RWMutex
	lockUsersServiceMockDeleteFamilyMember sync.RWMutex
	lockUsersServiceMockListFamilyMembers  sync.RWMutex
	lockUsersServiceMockRegister           sync.RWMutex
	lockUsersServiceMockUserByID           sync.RWMutex
)

// Ensure, that UsersServiceMock does implement UsersService.
// If this is not the case, regenerate this file with moq.
var _ UsersService = &UsersServiceMock{}

// UsersServiceMock is a mock implementation of UsersService.
//
//     func TestSomethingThatUsesUsersService(t *testing.T) {
//
//         // make and configure a mocked UsersService
//         mockedUsersService := &UsersServiceMock{
//             AddFamilyMemberFunc: func(ctx context.Context, r *users.User) (*users.User, error) {
// 	               panic("mock out the AddFamilyMember method")
//             },
//             DeleteFamilyMemberFunc: func(ctx context.Context, ownerID uint, memberID uint) error {
// 	               panic("mock out the DeleteFamilyMember method")
//             },
//             ListFamilyMembersFunc: func(ctx context.Context, id uint) ([]*users.User, error) {
// 	               panic("mock out the ListFamilyMembers method")
//             },
//             RegisterFunc: func(ctx context.Context, req *users.User) (*users.User, error) {
// 	               panic("mock out the Register method")
//             },
//             UserByIDFunc: func(ctx context.Context, id uint) (*users.User, error) {
// 	               panic("mock out the UserByID method")
//             },
//         }
//
//         // use mockedUsersService in code that requires UsersService
//         // and then make assertions.
//
//     }
type UsersServiceMock struct {
	// AddFamilyMemberFunc mocks the AddFamilyMember method.
	AddFamilyMemberFunc func(ctx context.Context, r *users.User) (*users.User, error)

	// DeleteFamilyMemberFunc mocks the DeleteFamilyMember method.
	DeleteFamilyMemberFunc func(ctx context.Context, ownerID uint, memberID uint) error

	// ListFamilyMembersFunc mocks the ListFamilyMembers method.
	ListFamilyMembersFunc func(ctx context.Context, id uint) ([]*users.User, error)

	// RegisterFunc mocks the Register method.
	RegisterFunc func(ctx context.Context, req *users.User) (*users.User, error)

	// UserByIDFunc mocks the UserByID method.
	UserByIDFunc func(ctx context.Context, id uint) (*users.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddFamilyMember holds details about calls to the AddFamilyMember method.
		AddFamilyMember []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *users.User
		}
		// DeleteFamilyMember holds details about calls to the DeleteFamilyMember method.
		DeleteFamilyMember []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// OwnerID is the ownerID argument value.
			OwnerID uint
			// MemberID is the memberID argument value.
			MemberID uint
		}
		// ListFamilyMembers holds details about calls to the ListFamilyMembers method.
		ListFamilyMembers []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID uint
		}
		// Register holds details about calls to the Register method.
		Register []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Req is the req argument value.
			Req *users.User
		}
		// UserByID holds details about calls to the UserByID method.
		UserByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID uint
		}
	}
}

// AddFamilyMember calls AddFamilyMemberFunc.
func (mock *UsersServiceMock) AddFamilyMember(ctx context.Context, r *users.User) (*users.User, error) {
	if mock.AddFamilyMemberFunc == nil {
		panic("UsersServiceMock.AddFamilyMemberFunc: method is nil but UsersService.AddFamilyMember was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *users.User
	}{
		Ctx: ctx,
		R:   r,
	}
	lockUsersServiceMockAddFamilyMember.Lock()
	mock.calls.AddFamilyMember = append(mock.calls.AddFamilyMember, callInfo)
	lockUsersServiceMockAddFamilyMember.Unlock()
	return mock.AddFamilyMemberFunc(ctx, r)
}

// AddFamilyMemberCalls gets all the calls that were made to AddFamilyMember.
// Check the length with:
//     len(mockedUsersService.AddFamilyMemberCalls())
func (mock *UsersServiceMock) AddFamilyMemberCalls() []struct {
	Ctx context.Context
	R   *users.User
} {
	var calls []struct {
		Ctx context.Context
		R   *users.User
	}
	lockUsersServiceMockAddFamilyMember.RLock()
	calls = mock.calls.AddFamilyMember
	lockUsersServiceMockAddFamilyMember.RUnlock()
	return calls
}

// DeleteFamilyMember calls DeleteFamilyMemberFunc.
func (mock *UsersServiceMock) DeleteFamilyMember(ctx context.Context, ownerID uint, memberID uint) error {
	if mock.DeleteFamilyMemberFunc == nil {
		panic("UsersServiceMock.DeleteFamilyMemberFunc: method is nil but UsersService.DeleteFamilyMember was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		OwnerID  uint
		MemberID uint
	}{
		Ctx:      ctx,
		OwnerID:  ownerID,
		MemberID: memberID,
	}
	lockUsersServiceMockDeleteFamilyMember.Lock()
	mock.calls.DeleteFamilyMember = append(mock.calls.DeleteFamilyMember, callInfo)
	lockUsersServiceMockDeleteFamilyMember.Unlock()
	return mock.DeleteFamilyMemberFunc(ctx, ownerID, memberID)
}

// DeleteFamilyMemberCalls gets all the calls that were made to DeleteFamilyMember.
// Check the length with:
//     len(mockedUsersService.DeleteFamilyMemberCalls())
func (mock *UsersServiceMock) DeleteFamilyMemberCalls() []struct {
	Ctx      context.Context
	OwnerID  uint
	MemberID uint
} {
	var calls []struct {
		Ctx      context.Context
		OwnerID  uint
		MemberID uint
	}
	lockUsersServiceMockDeleteFamilyMember.RLock()
	calls = mock.calls.DeleteFamilyMember
	lockUsersServiceMockDeleteFamilyMember.RUnlock()
	return calls
}

// ListFamilyMembers calls ListFamilyMembersFunc.
func (mock *UsersServiceMock) ListFamilyMembers(ctx context.Context, id uint) ([]*users.User, error) {
	if mock.ListFamilyMembersFunc == nil {
		panic("UsersServiceMock.ListFamilyMembersFunc: method is nil but UsersService.ListFamilyMembers was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  uint
	}{
		Ctx: ctx,
		ID:  id,
	}
	lockUsersServiceMockListFamilyMembers.Lock()
	mock.calls.ListFamilyMembers = append(mock.calls.ListFamilyMembers, callInfo)
	lockUsersServiceMockListFamilyMembers.Unlock()
	return mock.ListFamilyMembersFunc(ctx, id)
}

// ListFamilyMembersCalls gets all the calls that were made to ListFamilyMembers.
// Check the length with:
//     len(mockedUsersService.ListFamilyMembersCalls())
func (mock *UsersServiceMock) ListFamilyMembersCalls() []struct {
	Ctx context.Context
	ID  uint
} {
	var calls []struct {
		Ctx context.Context
		ID  uint
	}
	lockUsersServiceMockListFamilyMembers.RLock()
	calls = mock.calls.ListFamilyMembers
	lockUsersServiceMockListFamilyMembers.RUnlock()
	return calls
}

// Register calls RegisterFunc.
func (mock *UsersServiceMock) Register(ctx context.Context, req *users.User) (*users.User, error) {
	if mock.RegisterFunc == nil {
		panic("UsersServiceMock.RegisterFunc: method is nil but UsersService.Register was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Req *users.User
	}{
		Ctx: ctx,
		Req: req,
	}
	lockUsersServiceMockRegister.Lock()
	mock.calls.Register = append(mock.calls.Register, callInfo)
	lockUsersServiceMockRegister.Unlock()
	return mock.RegisterFunc(ctx, req)
}

// RegisterCalls gets all the calls that were made to Register.
// Check the length with:
//     len(mockedUsersService.RegisterCalls())
func (mock *UsersServiceMock) RegisterCalls() []struct {
	Ctx context.Context
	Req *users.User
} {
	var calls []struct {
		Ctx context.Context
		Req *users.User
	}
	lockUsersServiceMockRegister.RLock()
	calls = mock.calls.Register
	lockUsersServiceMockRegister.RUnlock()
	return calls
}

// UserByID calls UserByIDFunc.
func (mock *UsersServiceMock) UserByID(ctx context.Context, id uint) (*users.User, error) {
	if mock.UserByIDFunc == nil {
		panic("UsersServiceMock.UserByIDFunc: method is nil but UsersService.UserByID was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  uint
	}{
		Ctx: ctx,
		ID:  id,
	}
	lockUsersServiceMockUserByID.Lock()
	mock.calls.UserByID = append(mock.calls.UserByID, callInfo)
	lockUsersServiceMockUserByID.Unlock()
	return mock.UserByIDFunc(ctx, id)
}

// UserByIDCalls gets all the calls that were made to UserByID.
// Check the length with:
//     len(mockedUsersService.UserByIDCalls())
func (mock *UsersServiceMock) UserByIDCalls() []struct {
	Ctx context.Context
	ID  uint
} {
	var calls []struct {
		Ctx context.Context
		ID  uint
	}
	lockUsersServiceMockUserByID.RLock()
	calls = mock.calls.UserByID
	lockUsersServiceMockUserByID.RUnlock()
	return calls
}