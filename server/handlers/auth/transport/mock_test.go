// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package transport

import (
	"context"
	"github.com/ivch/dynasty/server/handlers/auth"
	"sync"
)

// Ensure, that AuthServiceMock does implement AuthService.
// If this is not the case, regenerate this file with moq.
var _ AuthService = &AuthServiceMock{}

// AuthServiceMock is a mock implementation of AuthService.
//
//	func TestSomethingThatUsesAuthService(t *testing.T) {
//
//		// make and configure a mocked AuthService
//		mockedAuthService := &AuthServiceMock{
//			GwfaFunc: func(token string) (uint, error) {
//				panic("mock out the Gwfa method")
//			},
//			LoginFunc: func(ctx context.Context, phone string, password string) (*auth.Tokens, error) {
//				panic("mock out the Login method")
//			},
//			LogoutFunc: func(ctx context.Context, id uint) error {
//				panic("mock out the Logout method")
//			},
//			RefreshFunc: func(ctx context.Context, token string) (*auth.Tokens, error) {
//				panic("mock out the Refresh method")
//			},
//		}
//
//		// use mockedAuthService in code that requires AuthService
//		// and then make assertions.
//
//	}
type AuthServiceMock struct {
	// GwfaFunc mocks the Gwfa method.
	GwfaFunc func(token string) (uint, error)

	// LoginFunc mocks the Login method.
	LoginFunc func(ctx context.Context, phone string, password string) (*auth.Tokens, error)

	// LogoutFunc mocks the Logout method.
	LogoutFunc func(ctx context.Context, id uint) error

	// RefreshFunc mocks the Refresh method.
	RefreshFunc func(ctx context.Context, token string) (*auth.Tokens, error)

	// calls tracks calls to the methods.
	calls struct {
		// Gwfa holds details about calls to the Gwfa method.
		Gwfa []struct {
			// Token is the token argument value.
			Token string
		}
		// Login holds details about calls to the Login method.
		Login []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Phone is the phone argument value.
			Phone string
			// Password is the password argument value.
			Password string
		}
		// Logout holds details about calls to the Logout method.
		Logout []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID uint
		}
		// Refresh holds details about calls to the Refresh method.
		Refresh []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Token is the token argument value.
			Token string
		}
	}
	lockGwfa    sync.RWMutex
	lockLogin   sync.RWMutex
	lockLogout  sync.RWMutex
	lockRefresh sync.RWMutex
}

// Gwfa calls GwfaFunc.
func (mock *AuthServiceMock) Gwfa(token string) (uint, error) {
	if mock.GwfaFunc == nil {
		panic("AuthServiceMock.GwfaFunc: method is nil but AuthService.Gwfa was just called")
	}
	callInfo := struct {
		Token string
	}{
		Token: token,
	}
	mock.lockGwfa.Lock()
	mock.calls.Gwfa = append(mock.calls.Gwfa, callInfo)
	mock.lockGwfa.Unlock()
	return mock.GwfaFunc(token)
}

// GwfaCalls gets all the calls that were made to Gwfa.
// Check the length with:
//
//	len(mockedAuthService.GwfaCalls())
func (mock *AuthServiceMock) GwfaCalls() []struct {
	Token string
} {
	var calls []struct {
		Token string
	}
	mock.lockGwfa.RLock()
	calls = mock.calls.Gwfa
	mock.lockGwfa.RUnlock()
	return calls
}

// Login calls LoginFunc.
func (mock *AuthServiceMock) Login(ctx context.Context, phone string, password string) (*auth.Tokens, error) {
	if mock.LoginFunc == nil {
		panic("AuthServiceMock.LoginFunc: method is nil but AuthService.Login was just called")
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
	mock.lockLogin.Lock()
	mock.calls.Login = append(mock.calls.Login, callInfo)
	mock.lockLogin.Unlock()
	return mock.LoginFunc(ctx, phone, password)
}

// LoginCalls gets all the calls that were made to Login.
// Check the length with:
//
//	len(mockedAuthService.LoginCalls())
func (mock *AuthServiceMock) LoginCalls() []struct {
	Ctx      context.Context
	Phone    string
	Password string
} {
	var calls []struct {
		Ctx      context.Context
		Phone    string
		Password string
	}
	mock.lockLogin.RLock()
	calls = mock.calls.Login
	mock.lockLogin.RUnlock()
	return calls
}

// Logout calls LogoutFunc.
func (mock *AuthServiceMock) Logout(ctx context.Context, id uint) error {
	if mock.LogoutFunc == nil {
		panic("AuthServiceMock.LogoutFunc: method is nil but AuthService.Logout was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  uint
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockLogout.Lock()
	mock.calls.Logout = append(mock.calls.Logout, callInfo)
	mock.lockLogout.Unlock()
	return mock.LogoutFunc(ctx, id)
}

// LogoutCalls gets all the calls that were made to Logout.
// Check the length with:
//
//	len(mockedAuthService.LogoutCalls())
func (mock *AuthServiceMock) LogoutCalls() []struct {
	Ctx context.Context
	ID  uint
} {
	var calls []struct {
		Ctx context.Context
		ID  uint
	}
	mock.lockLogout.RLock()
	calls = mock.calls.Logout
	mock.lockLogout.RUnlock()
	return calls
}

// Refresh calls RefreshFunc.
func (mock *AuthServiceMock) Refresh(ctx context.Context, token string) (*auth.Tokens, error) {
	if mock.RefreshFunc == nil {
		panic("AuthServiceMock.RefreshFunc: method is nil but AuthService.Refresh was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Token string
	}{
		Ctx:   ctx,
		Token: token,
	}
	mock.lockRefresh.Lock()
	mock.calls.Refresh = append(mock.calls.Refresh, callInfo)
	mock.lockRefresh.Unlock()
	return mock.RefreshFunc(ctx, token)
}

// RefreshCalls gets all the calls that were made to Refresh.
// Check the length with:
//
//	len(mockedAuthService.RefreshCalls())
func (mock *AuthServiceMock) RefreshCalls() []struct {
	Ctx   context.Context
	Token string
} {
	var calls []struct {
		Ctx   context.Context
		Token string
	}
	mock.lockRefresh.RLock()
	calls = mock.calls.Refresh
	mock.lockRefresh.RUnlock()
	return calls
}
