// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package requests

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

// Ensure, that requestsRepositoryMock does implement requestsRepository.
// If this is not the case, regenerate this file with moq.
var _ requestsRepository = &requestsRepositoryMock{}

// requestsRepositoryMock is a mock implementation of requestsRepository.
//
//	func TestSomethingThatUsesrequestsRepository(t *testing.T) {
//
//		// make and configure a mocked requestsRepository
//		mockedrequestsRepository := &requestsRepositoryMock{
//			AddImageFunc: func(userID uint, requestID uint, filename string) error {
//				panic("mock out the AddImage method")
//			},
//			CountForGuardFunc: func(req *RequestListFilter) (int, error) {
//				panic("mock out the CountForGuard method")
//			},
//			CreateFunc: func(req *Request) error {
//				panic("mock out the Create method")
//			},
//			DeleteFunc: func(id uint, userID uint) error {
//				panic("mock out the Delete method")
//			},
//			DeleteImageFunc: func(userID uint, requestID uint, filename string) error {
//				panic("mock out the DeleteImage method")
//			},
//			GetRequestByIDAndUserFunc: func(id uint, userID uint) (*Request, error) {
//				panic("mock out the GetRequestByIDAndUser method")
//			},
//			ListByUserFunc: func(r *RequestListFilter) ([]*Request, error) {
//				panic("mock out the ListByUser method")
//			},
//			ListForGuardFunc: func(req *RequestListFilter) ([]*Request, error) {
//				panic("mock out the ListForGuard method")
//			},
//			UpdateFunc: func(update *UpdateRequest) error {
//				panic("mock out the Update method")
//			},
//			UpdateForGuardFunc: func(id uint, status string) error {
//				panic("mock out the UpdateForGuard method")
//			},
//		}
//
//		// use mockedrequestsRepository in code that requires requestsRepository
//		// and then make assertions.
//
//	}
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
	lockAddImage              sync.RWMutex
	lockCountForGuard         sync.RWMutex
	lockCreate                sync.RWMutex
	lockDelete                sync.RWMutex
	lockDeleteImage           sync.RWMutex
	lockGetRequestByIDAndUser sync.RWMutex
	lockListByUser            sync.RWMutex
	lockListForGuard          sync.RWMutex
	lockUpdate                sync.RWMutex
	lockUpdateForGuard        sync.RWMutex
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
	mock.lockAddImage.Lock()
	mock.calls.AddImage = append(mock.calls.AddImage, callInfo)
	mock.lockAddImage.Unlock()
	return mock.AddImageFunc(userID, requestID, filename)
}

// AddImageCalls gets all the calls that were made to AddImage.
// Check the length with:
//
//	len(mockedrequestsRepository.AddImageCalls())
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
	mock.lockAddImage.RLock()
	calls = mock.calls.AddImage
	mock.lockAddImage.RUnlock()
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
	mock.lockCountForGuard.Lock()
	mock.calls.CountForGuard = append(mock.calls.CountForGuard, callInfo)
	mock.lockCountForGuard.Unlock()
	return mock.CountForGuardFunc(req)
}

// CountForGuardCalls gets all the calls that were made to CountForGuard.
// Check the length with:
//
//	len(mockedrequestsRepository.CountForGuardCalls())
func (mock *requestsRepositoryMock) CountForGuardCalls() []struct {
	Req *RequestListFilter
} {
	var calls []struct {
		Req *RequestListFilter
	}
	mock.lockCountForGuard.RLock()
	calls = mock.calls.CountForGuard
	mock.lockCountForGuard.RUnlock()
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
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(req)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedrequestsRepository.CreateCalls())
func (mock *requestsRepositoryMock) CreateCalls() []struct {
	Req *Request
} {
	var calls []struct {
		Req *Request
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
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
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(id, userID)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedrequestsRepository.DeleteCalls())
func (mock *requestsRepositoryMock) DeleteCalls() []struct {
	ID     uint
	UserID uint
} {
	var calls []struct {
		ID     uint
		UserID uint
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
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
	mock.lockDeleteImage.Lock()
	mock.calls.DeleteImage = append(mock.calls.DeleteImage, callInfo)
	mock.lockDeleteImage.Unlock()
	return mock.DeleteImageFunc(userID, requestID, filename)
}

// DeleteImageCalls gets all the calls that were made to DeleteImage.
// Check the length with:
//
//	len(mockedrequestsRepository.DeleteImageCalls())
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
	mock.lockDeleteImage.RLock()
	calls = mock.calls.DeleteImage
	mock.lockDeleteImage.RUnlock()
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
	mock.lockGetRequestByIDAndUser.Lock()
	mock.calls.GetRequestByIDAndUser = append(mock.calls.GetRequestByIDAndUser, callInfo)
	mock.lockGetRequestByIDAndUser.Unlock()
	return mock.GetRequestByIDAndUserFunc(id, userID)
}

// GetRequestByIDAndUserCalls gets all the calls that were made to GetRequestByIDAndUser.
// Check the length with:
//
//	len(mockedrequestsRepository.GetRequestByIDAndUserCalls())
func (mock *requestsRepositoryMock) GetRequestByIDAndUserCalls() []struct {
	ID     uint
	UserID uint
} {
	var calls []struct {
		ID     uint
		UserID uint
	}
	mock.lockGetRequestByIDAndUser.RLock()
	calls = mock.calls.GetRequestByIDAndUser
	mock.lockGetRequestByIDAndUser.RUnlock()
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
	mock.lockListByUser.Lock()
	mock.calls.ListByUser = append(mock.calls.ListByUser, callInfo)
	mock.lockListByUser.Unlock()
	return mock.ListByUserFunc(r)
}

// ListByUserCalls gets all the calls that were made to ListByUser.
// Check the length with:
//
//	len(mockedrequestsRepository.ListByUserCalls())
func (mock *requestsRepositoryMock) ListByUserCalls() []struct {
	R *RequestListFilter
} {
	var calls []struct {
		R *RequestListFilter
	}
	mock.lockListByUser.RLock()
	calls = mock.calls.ListByUser
	mock.lockListByUser.RUnlock()
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
	mock.lockListForGuard.Lock()
	mock.calls.ListForGuard = append(mock.calls.ListForGuard, callInfo)
	mock.lockListForGuard.Unlock()
	return mock.ListForGuardFunc(req)
}

// ListForGuardCalls gets all the calls that were made to ListForGuard.
// Check the length with:
//
//	len(mockedrequestsRepository.ListForGuardCalls())
func (mock *requestsRepositoryMock) ListForGuardCalls() []struct {
	Req *RequestListFilter
} {
	var calls []struct {
		Req *RequestListFilter
	}
	mock.lockListForGuard.RLock()
	calls = mock.calls.ListForGuard
	mock.lockListForGuard.RUnlock()
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
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(update)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedrequestsRepository.UpdateCalls())
func (mock *requestsRepositoryMock) UpdateCalls() []struct {
	Update *UpdateRequest
} {
	var calls []struct {
		Update *UpdateRequest
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
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
	mock.lockUpdateForGuard.Lock()
	mock.calls.UpdateForGuard = append(mock.calls.UpdateForGuard, callInfo)
	mock.lockUpdateForGuard.Unlock()
	return mock.UpdateForGuardFunc(id, status)
}

// UpdateForGuardCalls gets all the calls that were made to UpdateForGuard.
// Check the length with:
//
//	len(mockedrequestsRepository.UpdateForGuardCalls())
func (mock *requestsRepositoryMock) UpdateForGuardCalls() []struct {
	ID     uint
	Status string
} {
	var calls []struct {
		ID     uint
		Status string
	}
	mock.lockUpdateForGuard.RLock()
	calls = mock.calls.UpdateForGuard
	mock.lockUpdateForGuard.RUnlock()
	return calls
}

// Ensure, that s3ClientMock does implement s3Client.
// If this is not the case, regenerate this file with moq.
var _ s3Client = &s3ClientMock{}

// s3ClientMock is a mock implementation of s3Client.
//
//	func TestSomethingThatUsess3Client(t *testing.T) {
//
//		// make and configure a mocked s3Client
//		mockeds3Client := &s3ClientMock{
//			DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
//				panic("mock out the DeleteObject method")
//			},
//			PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
//				panic("mock out the PutObject method")
//			},
//		}
//
//		// use mockeds3Client in code that requires s3Client
//		// and then make assertions.
//
//	}
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
	lockDeleteObject sync.RWMutex
	lockPutObject    sync.RWMutex
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
	mock.lockDeleteObject.Lock()
	mock.calls.DeleteObject = append(mock.calls.DeleteObject, callInfo)
	mock.lockDeleteObject.Unlock()
	return mock.DeleteObjectFunc(input)
}

// DeleteObjectCalls gets all the calls that were made to DeleteObject.
// Check the length with:
//
//	len(mockeds3Client.DeleteObjectCalls())
func (mock *s3ClientMock) DeleteObjectCalls() []struct {
	Input *s3.DeleteObjectInput
} {
	var calls []struct {
		Input *s3.DeleteObjectInput
	}
	mock.lockDeleteObject.RLock()
	calls = mock.calls.DeleteObject
	mock.lockDeleteObject.RUnlock()
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
	mock.lockPutObject.Lock()
	mock.calls.PutObject = append(mock.calls.PutObject, callInfo)
	mock.lockPutObject.Unlock()
	return mock.PutObjectFunc(input)
}

// PutObjectCalls gets all the calls that were made to PutObject.
// Check the length with:
//
//	len(mockeds3Client.PutObjectCalls())
func (mock *s3ClientMock) PutObjectCalls() []struct {
	Input *s3.PutObjectInput
} {
	var calls []struct {
		Input *s3.PutObjectInput
	}
	mock.lockPutObject.RLock()
	calls = mock.calls.PutObject
	mock.lockPutObject.RUnlock()
	return calls
}
