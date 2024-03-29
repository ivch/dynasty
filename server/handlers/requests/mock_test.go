// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package requests

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

var (
	lockrequestsRepositoryMockAddImage              sync.RWMutex
	lockrequestsRepositoryMockCountForGuard         sync.RWMutex
	lockrequestsRepositoryMockCreate                sync.RWMutex
	lockrequestsRepositoryMockDelete                sync.RWMutex
	lockrequestsRepositoryMockDeleteImage           sync.RWMutex
	lockrequestsRepositoryMockGetRequestByIDAndUser sync.RWMutex
	lockrequestsRepositoryMockListByUser            sync.RWMutex
	lockrequestsRepositoryMockListForGuard          sync.RWMutex
	lockrequestsRepositoryMockUpdate                sync.RWMutex
	lockrequestsRepositoryMockUpdateForGuard        sync.RWMutex
)

// Ensure, that requestsRepositoryMock does implement requestsRepository.
// If this is not the case, regenerate this file with moq.
var _ requestsRepository = &requestsRepositoryMock{}

// requestsRepositoryMock is a mock implementation of requestsRepository.
//
//     func TestSomethingThatUsesrequestsRepository(t *testing.T) {
//
//         // make and configure a mocked requestsRepository
//         mockedrequestsRepository := &requestsRepositoryMock{
//             AddImageFunc: func(userID uint, requestID uint, filename string) error {
// 	               panic("mock out the AddImage method")
//             },
//             CountForGuardFunc: func(req *RequestListFilter) (int, error) {
// 	               panic("mock out the CountForGuard method")
//             },
//             CreateFunc: func(req *Request) error {
// 	               panic("mock out the Create method")
//             },
//             DeleteFunc: func(id uint, userID uint) error {
// 	               panic("mock out the Delete method")
//             },
//             DeleteImageFunc: func(userID uint, requestID uint, filename string) error {
// 	               panic("mock out the DeleteImage method")
//             },
//             GetRequestByIDAndUserFunc: func(id uint, userID uint) (*Request, error) {
// 	               panic("mock out the GetRequestByIDAndUser method")
//             },
//             ListByUserFunc: func(r *RequestListFilter) ([]*Request, error) {
// 	               panic("mock out the ListByUser method")
//             },
//             ListForGuardFunc: func(req *RequestListFilter) ([]*Request, error) {
// 	               panic("mock out the ListForGuard method")
//             },
//             UpdateFunc: func(update *UpdateRequest) error {
// 	               panic("mock out the Update method")
//             },
//             UpdateForGuardFunc: func(id uint, status string) error {
// 	               panic("mock out the UpdateForGuard method")
//             },
//         }
//
//         // use mockedrequestsRepository in code that requires requestsRepository
//         // and then make assertions.
//
//     }
type requestsRepositoryMock struct {
	// AddImageFunc mocks the AddImage method.
	AddImageFunc func(userID uint, requestID uint, filename string) error

	// CountForGuardFunc mocks the CountForGuard method.
	CountForGuardFunc func(req *RequestListFilter) (int, error)

	// CreateFunc mocks the Create method.
	CreateFunc func(req *Request) error

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(id uint, userID uint) error

	// DeleteImageFunc mocks the DeleteImage method.
	DeleteImageFunc func(userID uint, requestID uint, filename string) error

	// GetRequestByIDAndUserFunc mocks the GetRequestByIDAndUser method.
	GetRequestByIDAndUserFunc func(id uint, userID uint) (*Request, error)

	// ListByUserFunc mocks the ListByUser method.
	ListByUserFunc func(r *RequestListFilter) ([]*Request, error)

	// ListForGuardFunc mocks the ListForGuard method.
	ListForGuardFunc func(req *RequestListFilter) ([]*Request, error)

	// UpdateFunc mocks the Update method.
	UpdateFunc func(update *UpdateRequest) error

	// UpdateForGuardFunc mocks the UpdateForGuard method.
	UpdateForGuardFunc func(id uint, status string) error

	// calls tracks calls to the methods.
	calls struct {
		// AddImage holds details about calls to the AddImage method.
		AddImage []struct {
			// UserID is the userID argument value.
			UserID uint
			// RequestID is the requestID argument value.
			RequestID uint
			// Filename is the filename argument value.
			Filename string
		}
		// CountForGuard holds details about calls to the CountForGuard method.
		CountForGuard []struct {
			// Req is the req argument value.
			Req *RequestListFilter
		}
		// Create holds details about calls to the Create method.
		Create []struct {
			// Req is the req argument value.
			Req *Request
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// ID is the id argument value.
			ID uint
			// UserID is the userID argument value.
			UserID uint
		}
		// DeleteImage holds details about calls to the DeleteImage method.
		DeleteImage []struct {
			// UserID is the userID argument value.
			UserID uint
			// RequestID is the requestID argument value.
			RequestID uint
			// Filename is the filename argument value.
			Filename string
		}
		// GetRequestByIDAndUser holds details about calls to the GetRequestByIDAndUser method.
		GetRequestByIDAndUser []struct {
			// ID is the id argument value.
			ID uint
			// UserID is the userID argument value.
			UserID uint
		}
		// ListByUser holds details about calls to the ListByUser method.
		ListByUser []struct {
			// R is the r argument value.
			R *RequestListFilter
		}
		// ListForGuard holds details about calls to the ListForGuard method.
		ListForGuard []struct {
			// Req is the req argument value.
			Req *RequestListFilter
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Update is the update argument value.
			Update *UpdateRequest
		}
		// UpdateForGuard holds details about calls to the UpdateForGuard method.
		UpdateForGuard []struct {
			// ID is the id argument value.
			ID uint
			// Status is the status argument value.
			Status string
		}
	}
}

// AddImage calls AddImageFunc.
func (mock *requestsRepositoryMock) AddImage(userID uint, requestID uint, filename string) error {
	if mock.AddImageFunc == nil {
		panic("requestsRepositoryMock.AddImageFunc: method is nil but requestsRepository.AddImage was just called")
	}
	callInfo := struct {
		UserID    uint
		RequestID uint
		Filename  string
	}{
		UserID:    userID,
		RequestID: requestID,
		Filename:  filename,
	}
	lockrequestsRepositoryMockAddImage.Lock()
	mock.calls.AddImage = append(mock.calls.AddImage, callInfo)
	lockrequestsRepositoryMockAddImage.Unlock()
	return mock.AddImageFunc(userID, requestID, filename)
}

// AddImageCalls gets all the calls that were made to AddImage.
// Check the length with:
//     len(mockedrequestsRepository.AddImageCalls())
func (mock *requestsRepositoryMock) AddImageCalls() []struct {
	UserID    uint
	RequestID uint
	Filename  string
} {
	var calls []struct {
		UserID    uint
		RequestID uint
		Filename  string
	}
	lockrequestsRepositoryMockAddImage.RLock()
	calls = mock.calls.AddImage
	lockrequestsRepositoryMockAddImage.RUnlock()
	return calls
}

// CountForGuard calls CountForGuardFunc.
func (mock *requestsRepositoryMock) CountForGuard(req *RequestListFilter) (int, error) {
	if mock.CountForGuardFunc == nil {
		panic("requestsRepositoryMock.CountForGuardFunc: method is nil but requestsRepository.CountForGuard was just called")
	}
	callInfo := struct {
		Req *RequestListFilter
	}{
		Req: req,
	}
	lockrequestsRepositoryMockCountForGuard.Lock()
	mock.calls.CountForGuard = append(mock.calls.CountForGuard, callInfo)
	lockrequestsRepositoryMockCountForGuard.Unlock()
	return mock.CountForGuardFunc(req)
}

// CountForGuardCalls gets all the calls that were made to CountForGuard.
// Check the length with:
//     len(mockedrequestsRepository.CountForGuardCalls())
func (mock *requestsRepositoryMock) CountForGuardCalls() []struct {
	Req *RequestListFilter
} {
	var calls []struct {
		Req *RequestListFilter
	}
	lockrequestsRepositoryMockCountForGuard.RLock()
	calls = mock.calls.CountForGuard
	lockrequestsRepositoryMockCountForGuard.RUnlock()
	return calls
}

// Create calls CreateFunc.
func (mock *requestsRepositoryMock) Create(req *Request) error {
	if mock.CreateFunc == nil {
		panic("requestsRepositoryMock.CreateFunc: method is nil but requestsRepository.Create was just called")
	}
	callInfo := struct {
		Req *Request
	}{
		Req: req,
	}
	lockrequestsRepositoryMockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	lockrequestsRepositoryMockCreate.Unlock()
	return mock.CreateFunc(req)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//     len(mockedrequestsRepository.CreateCalls())
func (mock *requestsRepositoryMock) CreateCalls() []struct {
	Req *Request
} {
	var calls []struct {
		Req *Request
	}
	lockrequestsRepositoryMockCreate.RLock()
	calls = mock.calls.Create
	lockrequestsRepositoryMockCreate.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *requestsRepositoryMock) Delete(id uint, userID uint) error {
	if mock.DeleteFunc == nil {
		panic("requestsRepositoryMock.DeleteFunc: method is nil but requestsRepository.Delete was just called")
	}
	callInfo := struct {
		ID     uint
		UserID uint
	}{
		ID:     id,
		UserID: userID,
	}
	lockrequestsRepositoryMockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	lockrequestsRepositoryMockDelete.Unlock()
	return mock.DeleteFunc(id, userID)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//     len(mockedrequestsRepository.DeleteCalls())
func (mock *requestsRepositoryMock) DeleteCalls() []struct {
	ID     uint
	UserID uint
} {
	var calls []struct {
		ID     uint
		UserID uint
	}
	lockrequestsRepositoryMockDelete.RLock()
	calls = mock.calls.Delete
	lockrequestsRepositoryMockDelete.RUnlock()
	return calls
}

// DeleteImage calls DeleteImageFunc.
func (mock *requestsRepositoryMock) DeleteImage(userID uint, requestID uint, filename string) error {
	if mock.DeleteImageFunc == nil {
		panic("requestsRepositoryMock.DeleteImageFunc: method is nil but requestsRepository.DeleteImage was just called")
	}
	callInfo := struct {
		UserID    uint
		RequestID uint
		Filename  string
	}{
		UserID:    userID,
		RequestID: requestID,
		Filename:  filename,
	}
	lockrequestsRepositoryMockDeleteImage.Lock()
	mock.calls.DeleteImage = append(mock.calls.DeleteImage, callInfo)
	lockrequestsRepositoryMockDeleteImage.Unlock()
	return mock.DeleteImageFunc(userID, requestID, filename)
}

// DeleteImageCalls gets all the calls that were made to DeleteImage.
// Check the length with:
//     len(mockedrequestsRepository.DeleteImageCalls())
func (mock *requestsRepositoryMock) DeleteImageCalls() []struct {
	UserID    uint
	RequestID uint
	Filename  string
} {
	var calls []struct {
		UserID    uint
		RequestID uint
		Filename  string
	}
	lockrequestsRepositoryMockDeleteImage.RLock()
	calls = mock.calls.DeleteImage
	lockrequestsRepositoryMockDeleteImage.RUnlock()
	return calls
}

// GetRequestByIDAndUser calls GetRequestByIDAndUserFunc.
func (mock *requestsRepositoryMock) GetRequestByIDAndUser(id uint, userID uint) (*Request, error) {
	if mock.GetRequestByIDAndUserFunc == nil {
		panic("requestsRepositoryMock.GetRequestByIDAndUserFunc: method is nil but requestsRepository.GetRequestByIDAndUser was just called")
	}
	callInfo := struct {
		ID     uint
		UserID uint
	}{
		ID:     id,
		UserID: userID,
	}
	lockrequestsRepositoryMockGetRequestByIDAndUser.Lock()
	mock.calls.GetRequestByIDAndUser = append(mock.calls.GetRequestByIDAndUser, callInfo)
	lockrequestsRepositoryMockGetRequestByIDAndUser.Unlock()
	return mock.GetRequestByIDAndUserFunc(id, userID)
}

// GetRequestByIDAndUserCalls gets all the calls that were made to GetRequestByIDAndUser.
// Check the length with:
//     len(mockedrequestsRepository.GetRequestByIDAndUserCalls())
func (mock *requestsRepositoryMock) GetRequestByIDAndUserCalls() []struct {
	ID     uint
	UserID uint
} {
	var calls []struct {
		ID     uint
		UserID uint
	}
	lockrequestsRepositoryMockGetRequestByIDAndUser.RLock()
	calls = mock.calls.GetRequestByIDAndUser
	lockrequestsRepositoryMockGetRequestByIDAndUser.RUnlock()
	return calls
}

// ListByUser calls ListByUserFunc.
func (mock *requestsRepositoryMock) ListByUser(r *RequestListFilter) ([]*Request, error) {
	if mock.ListByUserFunc == nil {
		panic("requestsRepositoryMock.ListByUserFunc: method is nil but requestsRepository.ListByUser was just called")
	}
	callInfo := struct {
		R *RequestListFilter
	}{
		R: r,
	}
	lockrequestsRepositoryMockListByUser.Lock()
	mock.calls.ListByUser = append(mock.calls.ListByUser, callInfo)
	lockrequestsRepositoryMockListByUser.Unlock()
	return mock.ListByUserFunc(r)
}

// ListByUserCalls gets all the calls that were made to ListByUser.
// Check the length with:
//     len(mockedrequestsRepository.ListByUserCalls())
func (mock *requestsRepositoryMock) ListByUserCalls() []struct {
	R *RequestListFilter
} {
	var calls []struct {
		R *RequestListFilter
	}
	lockrequestsRepositoryMockListByUser.RLock()
	calls = mock.calls.ListByUser
	lockrequestsRepositoryMockListByUser.RUnlock()
	return calls
}

// ListForGuard calls ListForGuardFunc.
func (mock *requestsRepositoryMock) ListForGuard(req *RequestListFilter) ([]*Request, error) {
	if mock.ListForGuardFunc == nil {
		panic("requestsRepositoryMock.ListForGuardFunc: method is nil but requestsRepository.ListForGuard was just called")
	}
	callInfo := struct {
		Req *RequestListFilter
	}{
		Req: req,
	}
	lockrequestsRepositoryMockListForGuard.Lock()
	mock.calls.ListForGuard = append(mock.calls.ListForGuard, callInfo)
	lockrequestsRepositoryMockListForGuard.Unlock()
	return mock.ListForGuardFunc(req)
}

// ListForGuardCalls gets all the calls that were made to ListForGuard.
// Check the length with:
//     len(mockedrequestsRepository.ListForGuardCalls())
func (mock *requestsRepositoryMock) ListForGuardCalls() []struct {
	Req *RequestListFilter
} {
	var calls []struct {
		Req *RequestListFilter
	}
	lockrequestsRepositoryMockListForGuard.RLock()
	calls = mock.calls.ListForGuard
	lockrequestsRepositoryMockListForGuard.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *requestsRepositoryMock) Update(update *UpdateRequest) error {
	if mock.UpdateFunc == nil {
		panic("requestsRepositoryMock.UpdateFunc: method is nil but requestsRepository.Update was just called")
	}
	callInfo := struct {
		Update *UpdateRequest
	}{
		Update: update,
	}
	lockrequestsRepositoryMockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	lockrequestsRepositoryMockUpdate.Unlock()
	return mock.UpdateFunc(update)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//     len(mockedrequestsRepository.UpdateCalls())
func (mock *requestsRepositoryMock) UpdateCalls() []struct {
	Update *UpdateRequest
} {
	var calls []struct {
		Update *UpdateRequest
	}
	lockrequestsRepositoryMockUpdate.RLock()
	calls = mock.calls.Update
	lockrequestsRepositoryMockUpdate.RUnlock()
	return calls
}

// UpdateForGuard calls UpdateForGuardFunc.
func (mock *requestsRepositoryMock) UpdateForGuard(id uint, status string) error {
	if mock.UpdateForGuardFunc == nil {
		panic("requestsRepositoryMock.UpdateForGuardFunc: method is nil but requestsRepository.UpdateForGuard was just called")
	}
	callInfo := struct {
		ID     uint
		Status string
	}{
		ID:     id,
		Status: status,
	}
	lockrequestsRepositoryMockUpdateForGuard.Lock()
	mock.calls.UpdateForGuard = append(mock.calls.UpdateForGuard, callInfo)
	lockrequestsRepositoryMockUpdateForGuard.Unlock()
	return mock.UpdateForGuardFunc(id, status)
}

// UpdateForGuardCalls gets all the calls that were made to UpdateForGuard.
// Check the length with:
//     len(mockedrequestsRepository.UpdateForGuardCalls())
func (mock *requestsRepositoryMock) UpdateForGuardCalls() []struct {
	ID     uint
	Status string
} {
	var calls []struct {
		ID     uint
		Status string
	}
	lockrequestsRepositoryMockUpdateForGuard.RLock()
	calls = mock.calls.UpdateForGuard
	lockrequestsRepositoryMockUpdateForGuard.RUnlock()
	return calls
}

var (
	locks3ClientMockDeleteObject sync.RWMutex
	locks3ClientMockPutObject    sync.RWMutex
)

// Ensure, that s3ClientMock does implement s3Client.
// If this is not the case, regenerate this file with moq.
var _ s3Client = &s3ClientMock{}

// s3ClientMock is a mock implementation of s3Client.
//
//     func TestSomethingThatUsess3Client(t *testing.T) {
//
//         // make and configure a mocked s3Client
//         mockeds3Client := &s3ClientMock{
//             DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
// 	               panic("mock out the DeleteObject method")
//             },
//             PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
// 	               panic("mock out the PutObject method")
//             },
//         }
//
//         // use mockeds3Client in code that requires s3Client
//         // and then make assertions.
//
//     }
type s3ClientMock struct {
	// DeleteObjectFunc mocks the DeleteObject method.
	DeleteObjectFunc func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)

	// PutObjectFunc mocks the PutObject method.
	PutObjectFunc func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// DeleteObject holds details about calls to the DeleteObject method.
		DeleteObject []struct {
			// Input is the input argument value.
			Input *s3.DeleteObjectInput
		}
		// PutObject holds details about calls to the PutObject method.
		PutObject []struct {
			// Input is the input argument value.
			Input *s3.PutObjectInput
		}
	}
}

// DeleteObject calls DeleteObjectFunc.
func (mock *s3ClientMock) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	if mock.DeleteObjectFunc == nil {
		panic("s3ClientMock.DeleteObjectFunc: method is nil but s3Client.DeleteObject was just called")
	}
	callInfo := struct {
		Input *s3.DeleteObjectInput
	}{
		Input: input,
	}
	locks3ClientMockDeleteObject.Lock()
	mock.calls.DeleteObject = append(mock.calls.DeleteObject, callInfo)
	locks3ClientMockDeleteObject.Unlock()
	return mock.DeleteObjectFunc(input)
}

// DeleteObjectCalls gets all the calls that were made to DeleteObject.
// Check the length with:
//     len(mockeds3Client.DeleteObjectCalls())
func (mock *s3ClientMock) DeleteObjectCalls() []struct {
	Input *s3.DeleteObjectInput
} {
	var calls []struct {
		Input *s3.DeleteObjectInput
	}
	locks3ClientMockDeleteObject.RLock()
	calls = mock.calls.DeleteObject
	locks3ClientMockDeleteObject.RUnlock()
	return calls
}

// PutObject calls PutObjectFunc.
func (mock *s3ClientMock) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if mock.PutObjectFunc == nil {
		panic("s3ClientMock.PutObjectFunc: method is nil but s3Client.PutObject was just called")
	}
	callInfo := struct {
		Input *s3.PutObjectInput
	}{
		Input: input,
	}
	locks3ClientMockPutObject.Lock()
	mock.calls.PutObject = append(mock.calls.PutObject, callInfo)
	locks3ClientMockPutObject.Unlock()
	return mock.PutObjectFunc(input)
}

// PutObjectCalls gets all the calls that were made to PutObject.
// Check the length with:
//     len(mockeds3Client.PutObjectCalls())
func (mock *s3ClientMock) PutObjectCalls() []struct {
	Input *s3.PutObjectInput
} {
	var calls []struct {
		Input *s3.PutObjectInput
	}
	locks3ClientMockPutObject.RLock()
	calls = mock.calls.PutObject
	locks3ClientMockPutObject.RUnlock()
	return calls
}
