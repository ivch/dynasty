// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package transport

import (
	"context"
	"github.com/ivch/dynasty/server/handlers/requests"
	"sync"
)

// Ensure, that RequestsServiceMock does implement RequestsService.
// If this is not the case, regenerate this file with moq.
var _ RequestsService = &RequestsServiceMock{}

// RequestsServiceMock is a mock implementation of RequestsService.
//
//	func TestSomethingThatUsesRequestsService(t *testing.T) {
//
//		// make and configure a mocked RequestsService
//		mockedRequestsService := &RequestsServiceMock{
//			CreateFunc: func(ctx context.Context, r *requests.Request) (*requests.Request, error) {
//				panic("mock out the Create method")
//			},
//			DeleteFunc: func(ctx context.Context, r *requests.Request) error {
//				panic("mock out the Delete method")
//			},
//			DeleteImageFunc: func(ctx context.Context, r *requests.Image) error {
//				panic("mock out the DeleteImage method")
//			},
//			GetFunc: func(ctx context.Context, r *requests.Request) (*requests.Request, error) {
//				panic("mock out the Get method")
//			},
//			GuardRequestListFunc: func(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, int, error) {
//				panic("mock out the GuardRequestList method")
//			},
//			GuardUpdateRequestFunc: func(ctx context.Context, r *requests.Request) error {
//				panic("mock out the GuardUpdateRequest method")
//			},
//			MyFunc: func(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, error) {
//				panic("mock out the My method")
//			},
//			UpdateFunc: func(ctx context.Context, r *requests.UpdateRequest) error {
//				panic("mock out the Update method")
//			},
//			UploadImageFunc: func(ctx context.Context, r *requests.Image) (*requests.Image, error) {
//				panic("mock out the UploadImage method")
//			},
//		}
//
//		// use mockedRequestsService in code that requires RequestsService
//		// and then make assertions.
//
//	}
type RequestsServiceMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, r *requests.Request) (*requests.Request, error)

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(ctx context.Context, r *requests.Request) error

	// DeleteImageFunc mocks the DeleteImage method.
	DeleteImageFunc func(ctx context.Context, r *requests.Image) error

	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, r *requests.Request) (*requests.Request, error)

	// GuardRequestListFunc mocks the GuardRequestList method.
	GuardRequestListFunc func(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, int, error)

	// GuardUpdateRequestFunc mocks the GuardUpdateRequest method.
	GuardUpdateRequestFunc func(ctx context.Context, r *requests.Request) error

	// MyFunc mocks the My method.
	MyFunc func(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, error)

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, r *requests.UpdateRequest) error

	// UploadImageFunc mocks the UploadImage method.
	UploadImageFunc func(ctx context.Context, r *requests.Image) (*requests.Image, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.Request
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.Request
		}
		// DeleteImage holds details about calls to the DeleteImage method.
		DeleteImage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.Image
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.Request
		}
		// GuardRequestList holds details about calls to the GuardRequestList method.
		GuardRequestList []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.RequestListFilter
		}
		// GuardUpdateRequest holds details about calls to the GuardUpdateRequest method.
		GuardUpdateRequest []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.Request
		}
		// My holds details about calls to the My method.
		My []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.RequestListFilter
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.UpdateRequest
		}
		// UploadImage holds details about calls to the UploadImage method.
		UploadImage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *requests.Image
		}
	}
	lockCreate             sync.RWMutex
	lockDelete             sync.RWMutex
	lockDeleteImage        sync.RWMutex
	lockGet                sync.RWMutex
	lockGuardRequestList   sync.RWMutex
	lockGuardUpdateRequest sync.RWMutex
	lockMy                 sync.RWMutex
	lockUpdate             sync.RWMutex
	lockUploadImage        sync.RWMutex
}

// Create calls CreateFunc.
func (mock *RequestsServiceMock) Create(ctx context.Context, r *requests.Request) (*requests.Request, error) {
	if mock.CreateFunc == nil {
		panic("RequestsServiceMock.CreateFunc: method is nil but RequestsService.Create was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.Request
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, r)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedRequestsService.CreateCalls())
func (mock *RequestsServiceMock) CreateCalls() []struct {
	Ctx context.Context
	R   *requests.Request
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.Request
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *RequestsServiceMock) Delete(ctx context.Context, r *requests.Request) error {
	if mock.DeleteFunc == nil {
		panic("RequestsServiceMock.DeleteFunc: method is nil but RequestsService.Delete was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.Request
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(ctx, r)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedRequestsService.DeleteCalls())
func (mock *RequestsServiceMock) DeleteCalls() []struct {
	Ctx context.Context
	R   *requests.Request
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.Request
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// DeleteImage calls DeleteImageFunc.
func (mock *RequestsServiceMock) DeleteImage(ctx context.Context, r *requests.Image) error {
	if mock.DeleteImageFunc == nil {
		panic("RequestsServiceMock.DeleteImageFunc: method is nil but RequestsService.DeleteImage was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.Image
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockDeleteImage.Lock()
	mock.calls.DeleteImage = append(mock.calls.DeleteImage, callInfo)
	mock.lockDeleteImage.Unlock()
	return mock.DeleteImageFunc(ctx, r)
}

// DeleteImageCalls gets all the calls that were made to DeleteImage.
// Check the length with:
//
//	len(mockedRequestsService.DeleteImageCalls())
func (mock *RequestsServiceMock) DeleteImageCalls() []struct {
	Ctx context.Context
	R   *requests.Image
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.Image
	}
	mock.lockDeleteImage.RLock()
	calls = mock.calls.DeleteImage
	mock.lockDeleteImage.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *RequestsServiceMock) Get(ctx context.Context, r *requests.Request) (*requests.Request, error) {
	if mock.GetFunc == nil {
		panic("RequestsServiceMock.GetFunc: method is nil but RequestsService.Get was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.Request
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, r)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedRequestsService.GetCalls())
func (mock *RequestsServiceMock) GetCalls() []struct {
	Ctx context.Context
	R   *requests.Request
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.Request
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// GuardRequestList calls GuardRequestListFunc.
func (mock *RequestsServiceMock) GuardRequestList(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, int, error) {
	if mock.GuardRequestListFunc == nil {
		panic("RequestsServiceMock.GuardRequestListFunc: method is nil but RequestsService.GuardRequestList was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.RequestListFilter
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockGuardRequestList.Lock()
	mock.calls.GuardRequestList = append(mock.calls.GuardRequestList, callInfo)
	mock.lockGuardRequestList.Unlock()
	return mock.GuardRequestListFunc(ctx, r)
}

// GuardRequestListCalls gets all the calls that were made to GuardRequestList.
// Check the length with:
//
//	len(mockedRequestsService.GuardRequestListCalls())
func (mock *RequestsServiceMock) GuardRequestListCalls() []struct {
	Ctx context.Context
	R   *requests.RequestListFilter
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.RequestListFilter
	}
	mock.lockGuardRequestList.RLock()
	calls = mock.calls.GuardRequestList
	mock.lockGuardRequestList.RUnlock()
	return calls
}

// GuardUpdateRequest calls GuardUpdateRequestFunc.
func (mock *RequestsServiceMock) GuardUpdateRequest(ctx context.Context, r *requests.Request) error {
	if mock.GuardUpdateRequestFunc == nil {
		panic("RequestsServiceMock.GuardUpdateRequestFunc: method is nil but RequestsService.GuardUpdateRequest was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.Request
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockGuardUpdateRequest.Lock()
	mock.calls.GuardUpdateRequest = append(mock.calls.GuardUpdateRequest, callInfo)
	mock.lockGuardUpdateRequest.Unlock()
	return mock.GuardUpdateRequestFunc(ctx, r)
}

// GuardUpdateRequestCalls gets all the calls that were made to GuardUpdateRequest.
// Check the length with:
//
//	len(mockedRequestsService.GuardUpdateRequestCalls())
func (mock *RequestsServiceMock) GuardUpdateRequestCalls() []struct {
	Ctx context.Context
	R   *requests.Request
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.Request
	}
	mock.lockGuardUpdateRequest.RLock()
	calls = mock.calls.GuardUpdateRequest
	mock.lockGuardUpdateRequest.RUnlock()
	return calls
}

// My calls MyFunc.
func (mock *RequestsServiceMock) My(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, error) {
	if mock.MyFunc == nil {
		panic("RequestsServiceMock.MyFunc: method is nil but RequestsService.My was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.RequestListFilter
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockMy.Lock()
	mock.calls.My = append(mock.calls.My, callInfo)
	mock.lockMy.Unlock()
	return mock.MyFunc(ctx, r)
}

// MyCalls gets all the calls that were made to My.
// Check the length with:
//
//	len(mockedRequestsService.MyCalls())
func (mock *RequestsServiceMock) MyCalls() []struct {
	Ctx context.Context
	R   *requests.RequestListFilter
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.RequestListFilter
	}
	mock.lockMy.RLock()
	calls = mock.calls.My
	mock.lockMy.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *RequestsServiceMock) Update(ctx context.Context, r *requests.UpdateRequest) error {
	if mock.UpdateFunc == nil {
		panic("RequestsServiceMock.UpdateFunc: method is nil but RequestsService.Update was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.UpdateRequest
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(ctx, r)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedRequestsService.UpdateCalls())
func (mock *RequestsServiceMock) UpdateCalls() []struct {
	Ctx context.Context
	R   *requests.UpdateRequest
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.UpdateRequest
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}

// UploadImage calls UploadImageFunc.
func (mock *RequestsServiceMock) UploadImage(ctx context.Context, r *requests.Image) (*requests.Image, error) {
	if mock.UploadImageFunc == nil {
		panic("RequestsServiceMock.UploadImageFunc: method is nil but RequestsService.UploadImage was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *requests.Image
	}{
		Ctx: ctx,
		R:   r,
	}
	mock.lockUploadImage.Lock()
	mock.calls.UploadImage = append(mock.calls.UploadImage, callInfo)
	mock.lockUploadImage.Unlock()
	return mock.UploadImageFunc(ctx, r)
}

// UploadImageCalls gets all the calls that were made to UploadImage.
// Check the length with:
//
//	len(mockedRequestsService.UploadImageCalls())
func (mock *RequestsServiceMock) UploadImageCalls() []struct {
	Ctx context.Context
	R   *requests.Image
} {
	var calls []struct {
		Ctx context.Context
		R   *requests.Image
	}
	mock.lockUploadImage.RLock()
	calls = mock.calls.UploadImage
	mock.lockUploadImage.RUnlock()
	return calls
}
