// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package users

import (
	"context"
	"github.com/ivch/dynasty/models/dto"
	"sync"
)

var (
	lockuserServiceMockRegister               sync.RWMutex
	lockuserServiceMockUserByID               sync.RWMutex
	lockuserServiceMockUserByPhoneAndPassword sync.RWMutex
)

// Ensure, that userServiceMock does implement userService.
// If this is not the case, regenerate this file with moq.
var _ userService = &userServiceMock{}

// userServiceMock is a mock implementation of userService.
//
//     func TestSomethingThatUsesuserService(t *testing.T) {
//
//         // make and configure a mocked userService
//         mockeduserService := &userServiceMock{
//             RegisterFunc: func(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
// 	               panic("mock out the Register method")
//             },
//             UserByIDFunc: func(ctx context.Context, id uint) (*dto.UserByIDResponse, error) {
// 	               panic("mock out the UserByID method")
//             },
//             UserByPhoneAndPasswordFunc: func(ctx context.Context, phone string, password string) (*dto.UserAuthResponse, error) {
// 	               panic("mock out the UserByPhoneAndPassword method")
//             },
//         }
//
//         // use mockeduserService in code that requires userService
//         // and then make assertions.
//
//     }
type userServiceMock struct {
	// RegisterFunc mocks the Register method.
	RegisterFunc func(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)

	// UserByIDFunc mocks the UserByID method.
	UserByIDFunc func(ctx context.Context, id uint) (*dto.UserByIDResponse, error)

	// UserByPhoneAndPasswordFunc mocks the UserByPhoneAndPassword method.
	UserByPhoneAndPasswordFunc func(ctx context.Context, phone string, password string) (*dto.UserAuthResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// Register holds details about calls to the Register method.
		Register []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Req is the req argument value.
			Req *dto.UserRegisterRequest
		}
		// UserByID holds details about calls to the UserByID method.
		UserByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID uint
		}
		// UserByPhoneAndPassword holds details about calls to the UserByPhoneAndPassword method.
		UserByPhoneAndPassword []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Phone is the phone argument value.
			Phone string
			// Password is the password argument value.
			Password string
		}
	}
}

// Register calls RegisterFunc.
func (mock *userServiceMock) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	if mock.RegisterFunc == nil {
		panic("userServiceMock.RegisterFunc: method is nil but userService.Register was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Req *dto.UserRegisterRequest
	}{
		Ctx: ctx,
		Req: req,
	}
	lockuserServiceMockRegister.Lock()
	mock.calls.Register = append(mock.calls.Register, callInfo)
	lockuserServiceMockRegister.Unlock()
	return mock.RegisterFunc(ctx, req)
}

// RegisterCalls gets all the calls that were made to Register.
// Check the length with:
//     len(mockeduserService.RegisterCalls())
func (mock *userServiceMock) RegisterCalls() []struct {
	Ctx context.Context
	Req *dto.UserRegisterRequest
} {
	var calls []struct {
		Ctx context.Context
		Req *dto.UserRegisterRequest
	}
	lockuserServiceMockRegister.RLock()
	calls = mock.calls.Register
	lockuserServiceMockRegister.RUnlock()
	return calls
}

// UserByID calls UserByIDFunc.
func (mock *userServiceMock) UserByID(ctx context.Context, id uint) (*dto.UserByIDResponse, error) {
	if mock.UserByIDFunc == nil {
		panic("userServiceMock.UserByIDFunc: method is nil but userService.UserByID was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  uint
	}{
		Ctx: ctx,
		ID:  id,
	}
	lockuserServiceMockUserByID.Lock()
	mock.calls.UserByID = append(mock.calls.UserByID, callInfo)
	lockuserServiceMockUserByID.Unlock()
	return mock.UserByIDFunc(ctx, id)
}

// UserByIDCalls gets all the calls that were made to UserByID.
// Check the length with:
//     len(mockeduserService.UserByIDCalls())
func (mock *userServiceMock) UserByIDCalls() []struct {
	Ctx context.Context
	ID  uint
} {
	var calls []struct {
		Ctx context.Context
		ID  uint
	}
	lockuserServiceMockUserByID.RLock()
	calls = mock.calls.UserByID
	lockuserServiceMockUserByID.RUnlock()
	return calls
}

// UserByPhoneAndPassword calls UserByPhoneAndPasswordFunc.
func (mock *userServiceMock) UserByPhoneAndPassword(ctx context.Context, phone string, password string) (*dto.UserAuthResponse, error) {
	if mock.UserByPhoneAndPasswordFunc == nil {
		panic("userServiceMock.UserByPhoneAndPasswordFunc: method is nil but userService.UserByPhoneAndPassword was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Phone    string
		Password string
	}{
		Ctx:      ctx,
		Phone:    phone,
		Password: password,
	}
	lockuserServiceMockUserByPhoneAndPassword.Lock()
	mock.calls.UserByPhoneAndPassword = append(mock.calls.UserByPhoneAndPassword, callInfo)
	lockuserServiceMockUserByPhoneAndPassword.Unlock()
	return mock.UserByPhoneAndPasswordFunc(ctx, phone, password)
}

// UserByPhoneAndPasswordCalls gets all the calls that were made to UserByPhoneAndPassword.
// Check the length with:
//     len(mockeduserService.UserByPhoneAndPasswordCalls())
func (mock *userServiceMock) UserByPhoneAndPasswordCalls() []struct {
	Ctx      context.Context
	Phone    string
	Password string
} {
	var calls []struct {
		Ctx      context.Context
		Phone    string
		Password string
	}
	lockuserServiceMockUserByPhoneAndPassword.RLock()
	calls = mock.calls.UserByPhoneAndPassword
	lockuserServiceMockUserByPhoneAndPassword.RUnlock()
	return calls
}