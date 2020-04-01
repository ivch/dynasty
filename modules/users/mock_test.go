// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package users

import (
	"context"
	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
	"sync"
)

var (
	lockuserRepositoryMockCreateUser          sync.RWMutex
	lockuserRepositoryMockDeleteUser          sync.RWMutex
	lockuserRepositoryMockFindUserByApartment sync.RWMutex
	lockuserRepositoryMockGetFamilyMembers    sync.RWMutex
	lockuserRepositoryMockGetRegCode          sync.RWMutex
	lockuserRepositoryMockGetUserByID         sync.RWMutex
	lockuserRepositoryMockGetUserByPhone      sync.RWMutex
	lockuserRepositoryMockUpdateUser          sync.RWMutex
	lockuserRepositoryMockUseRegCode          sync.RWMutex
	lockuserRepositoryMockValidateRegCode     sync.RWMutex
)

// Ensure, that userRepositoryMock does implement userRepository.
// If this is not the case, regenerate this file with moq.
var _ userRepository = &userRepositoryMock{}

// userRepositoryMock is a mock implementation of userRepository.
//
//     func TestSomethingThatUsesuserRepository(t *testing.T) {
//
//         // make and configure a mocked userRepository
//         mockeduserRepository := &userRepositoryMock{
//             CreateUserFunc: func(user *entities.User) error {
// 	               panic("mock out the CreateUser method")
//             },
//             DeleteUserFunc: func(u *entities.User) error {
// 	               panic("mock out the DeleteUser method")
//             },
//             FindUserByApartmentFunc: func(building uint, apt uint) (*entities.User, error) {
// 	               panic("mock out the FindUserByApartment method")
//             },
//             GetFamilyMembersFunc: func(ownerID uint) ([]*entities.User, error) {
// 	               panic("mock out the GetFamilyMembers method")
//             },
//             GetRegCodeFunc: func() (string, error) {
// 	               panic("mock out the GetRegCode method")
//             },
//             GetUserByIDFunc: func(id uint) (*entities.User, error) {
// 	               panic("mock out the GetUserByID method")
//             },
//             GetUserByPhoneFunc: func(phone string) (*entities.User, error) {
// 	               panic("mock out the GetUserByPhone method")
//             },
//             UpdateUserFunc: func(u *entities.User) error {
// 	               panic("mock out the UpdateUser method")
//             },
//             UseRegCodeFunc: func(code string) error {
// 	               panic("mock out the UseRegCode method")
//             },
//             ValidateRegCodeFunc: func(code string) error {
// 	               panic("mock out the ValidateRegCode method")
//             },
//         }
//
//         // use mockeduserRepository in code that requires userRepository
//         // and then make assertions.
//
//     }
type userRepositoryMock struct {
	// CreateUserFunc mocks the CreateUser method.
	CreateUserFunc func(user *entities.User) error

	// DeleteUserFunc mocks the DeleteUser method.
	DeleteUserFunc func(u *entities.User) error

	// FindUserByApartmentFunc mocks the FindUserByApartment method.
	FindUserByApartmentFunc func(building uint, apt uint) (*entities.User, error)

	// GetFamilyMembersFunc mocks the GetFamilyMembers method.
	GetFamilyMembersFunc func(ownerID uint) ([]*entities.User, error)

	// GetRegCodeFunc mocks the GetRegCode method.
	GetRegCodeFunc func() (string, error)

	// GetUserByIDFunc mocks the GetUserByID method.
	GetUserByIDFunc func(id uint) (*entities.User, error)

	// GetUserByPhoneFunc mocks the GetUserByPhone method.
	GetUserByPhoneFunc func(phone string) (*entities.User, error)

	// UpdateUserFunc mocks the UpdateUser method.
	UpdateUserFunc func(u *entities.User) error

	// UseRegCodeFunc mocks the UseRegCode method.
	UseRegCodeFunc func(code string) error

	// ValidateRegCodeFunc mocks the ValidateRegCode method.
	ValidateRegCodeFunc func(code string) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateUser holds details about calls to the CreateUser method.
		CreateUser []struct {
			// User is the user argument value.
			User *entities.User
		}
		// DeleteUser holds details about calls to the DeleteUser method.
		DeleteUser []struct {
			// U is the u argument value.
			U *entities.User
		}
		// FindUserByApartment holds details about calls to the FindUserByApartment method.
		FindUserByApartment []struct {
			// Building is the building argument value.
			Building uint
			// Apt is the apt argument value.
			Apt uint
		}
		// GetFamilyMembers holds details about calls to the GetFamilyMembers method.
		GetFamilyMembers []struct {
			// OwnerID is the ownerID argument value.
			OwnerID uint
		}
		// GetRegCode holds details about calls to the GetRegCode method.
		GetRegCode []struct {
		}
		// GetUserByID holds details about calls to the GetUserByID method.
		GetUserByID []struct {
			// ID is the id argument value.
			ID uint
		}
		// GetUserByPhone holds details about calls to the GetUserByPhone method.
		GetUserByPhone []struct {
			// Phone is the phone argument value.
			Phone string
		}
		// UpdateUser holds details about calls to the UpdateUser method.
		UpdateUser []struct {
			// U is the u argument value.
			U *entities.User
		}
		// UseRegCode holds details about calls to the UseRegCode method.
		UseRegCode []struct {
			// Code is the code argument value.
			Code string
		}
		// ValidateRegCode holds details about calls to the ValidateRegCode method.
		ValidateRegCode []struct {
			// Code is the code argument value.
			Code string
		}
	}
}

// CreateUser calls CreateUserFunc.
func (mock *userRepositoryMock) CreateUser(user *entities.User) error {
	if mock.CreateUserFunc == nil {
		panic("userRepositoryMock.CreateUserFunc: method is nil but userRepository.CreateUser was just called")
	}
	callInfo := struct {
		User *entities.User
	}{
		User: user,
	}
	lockuserRepositoryMockCreateUser.Lock()
	mock.calls.CreateUser = append(mock.calls.CreateUser, callInfo)
	lockuserRepositoryMockCreateUser.Unlock()
	return mock.CreateUserFunc(user)
}

// CreateUserCalls gets all the calls that were made to CreateUser.
// Check the length with:
//     len(mockeduserRepository.CreateUserCalls())
func (mock *userRepositoryMock) CreateUserCalls() []struct {
	User *entities.User
} {
	var calls []struct {
		User *entities.User
	}
	lockuserRepositoryMockCreateUser.RLock()
	calls = mock.calls.CreateUser
	lockuserRepositoryMockCreateUser.RUnlock()
	return calls
}

// DeleteUser calls DeleteUserFunc.
func (mock *userRepositoryMock) DeleteUser(u *entities.User) error {
	if mock.DeleteUserFunc == nil {
		panic("userRepositoryMock.DeleteUserFunc: method is nil but userRepository.DeleteUser was just called")
	}
	callInfo := struct {
		U *entities.User
	}{
		U: u,
	}
	lockuserRepositoryMockDeleteUser.Lock()
	mock.calls.DeleteUser = append(mock.calls.DeleteUser, callInfo)
	lockuserRepositoryMockDeleteUser.Unlock()
	return mock.DeleteUserFunc(u)
}

// DeleteUserCalls gets all the calls that were made to DeleteUser.
// Check the length with:
//     len(mockeduserRepository.DeleteUserCalls())
func (mock *userRepositoryMock) DeleteUserCalls() []struct {
	U *entities.User
} {
	var calls []struct {
		U *entities.User
	}
	lockuserRepositoryMockDeleteUser.RLock()
	calls = mock.calls.DeleteUser
	lockuserRepositoryMockDeleteUser.RUnlock()
	return calls
}

// FindUserByApartment calls FindUserByApartmentFunc.
func (mock *userRepositoryMock) FindUserByApartment(building uint, apt uint) (*entities.User, error) {
	if mock.FindUserByApartmentFunc == nil {
		panic("userRepositoryMock.FindUserByApartmentFunc: method is nil but userRepository.FindUserByApartment was just called")
	}
	callInfo := struct {
		Building uint
		Apt      uint
	}{
		Building: building,
		Apt:      apt,
	}
	lockuserRepositoryMockFindUserByApartment.Lock()
	mock.calls.FindUserByApartment = append(mock.calls.FindUserByApartment, callInfo)
	lockuserRepositoryMockFindUserByApartment.Unlock()
	return mock.FindUserByApartmentFunc(building, apt)
}

// FindUserByApartmentCalls gets all the calls that were made to FindUserByApartment.
// Check the length with:
//     len(mockeduserRepository.FindUserByApartmentCalls())
func (mock *userRepositoryMock) FindUserByApartmentCalls() []struct {
	Building uint
	Apt      uint
} {
	var calls []struct {
		Building uint
		Apt      uint
	}
	lockuserRepositoryMockFindUserByApartment.RLock()
	calls = mock.calls.FindUserByApartment
	lockuserRepositoryMockFindUserByApartment.RUnlock()
	return calls
}

// GetFamilyMembers calls GetFamilyMembersFunc.
func (mock *userRepositoryMock) GetFamilyMembers(ownerID uint) ([]*entities.User, error) {
	if mock.GetFamilyMembersFunc == nil {
		panic("userRepositoryMock.GetFamilyMembersFunc: method is nil but userRepository.GetFamilyMembers was just called")
	}
	callInfo := struct {
		OwnerID uint
	}{
		OwnerID: ownerID,
	}
	lockuserRepositoryMockGetFamilyMembers.Lock()
	mock.calls.GetFamilyMembers = append(mock.calls.GetFamilyMembers, callInfo)
	lockuserRepositoryMockGetFamilyMembers.Unlock()
	return mock.GetFamilyMembersFunc(ownerID)
}

// GetFamilyMembersCalls gets all the calls that were made to GetFamilyMembers.
// Check the length with:
//     len(mockeduserRepository.GetFamilyMembersCalls())
func (mock *userRepositoryMock) GetFamilyMembersCalls() []struct {
	OwnerID uint
} {
	var calls []struct {
		OwnerID uint
	}
	lockuserRepositoryMockGetFamilyMembers.RLock()
	calls = mock.calls.GetFamilyMembers
	lockuserRepositoryMockGetFamilyMembers.RUnlock()
	return calls
}

// GetRegCode calls GetRegCodeFunc.
func (mock *userRepositoryMock) GetRegCode() (string, error) {
	if mock.GetRegCodeFunc == nil {
		panic("userRepositoryMock.GetRegCodeFunc: method is nil but userRepository.GetRegCode was just called")
	}
	callInfo := struct {
	}{}
	lockuserRepositoryMockGetRegCode.Lock()
	mock.calls.GetRegCode = append(mock.calls.GetRegCode, callInfo)
	lockuserRepositoryMockGetRegCode.Unlock()
	return mock.GetRegCodeFunc()
}

// GetRegCodeCalls gets all the calls that were made to GetRegCode.
// Check the length with:
//     len(mockeduserRepository.GetRegCodeCalls())
func (mock *userRepositoryMock) GetRegCodeCalls() []struct {
} {
	var calls []struct {
	}
	lockuserRepositoryMockGetRegCode.RLock()
	calls = mock.calls.GetRegCode
	lockuserRepositoryMockGetRegCode.RUnlock()
	return calls
}

// GetUserByID calls GetUserByIDFunc.
func (mock *userRepositoryMock) GetUserByID(id uint) (*entities.User, error) {
	if mock.GetUserByIDFunc == nil {
		panic("userRepositoryMock.GetUserByIDFunc: method is nil but userRepository.GetUserByID was just called")
	}
	callInfo := struct {
		ID uint
	}{
		ID: id,
	}
	lockuserRepositoryMockGetUserByID.Lock()
	mock.calls.GetUserByID = append(mock.calls.GetUserByID, callInfo)
	lockuserRepositoryMockGetUserByID.Unlock()
	return mock.GetUserByIDFunc(id)
}

// GetUserByIDCalls gets all the calls that were made to GetUserByID.
// Check the length with:
//     len(mockeduserRepository.GetUserByIDCalls())
func (mock *userRepositoryMock) GetUserByIDCalls() []struct {
	ID uint
} {
	var calls []struct {
		ID uint
	}
	lockuserRepositoryMockGetUserByID.RLock()
	calls = mock.calls.GetUserByID
	lockuserRepositoryMockGetUserByID.RUnlock()
	return calls
}

// GetUserByPhone calls GetUserByPhoneFunc.
func (mock *userRepositoryMock) GetUserByPhone(phone string) (*entities.User, error) {
	if mock.GetUserByPhoneFunc == nil {
		panic("userRepositoryMock.GetUserByPhoneFunc: method is nil but userRepository.GetUserByPhone was just called")
	}
	callInfo := struct {
		Phone string
	}{
		Phone: phone,
	}
	lockuserRepositoryMockGetUserByPhone.Lock()
	mock.calls.GetUserByPhone = append(mock.calls.GetUserByPhone, callInfo)
	lockuserRepositoryMockGetUserByPhone.Unlock()
	return mock.GetUserByPhoneFunc(phone)
}

// GetUserByPhoneCalls gets all the calls that were made to GetUserByPhone.
// Check the length with:
//     len(mockeduserRepository.GetUserByPhoneCalls())
func (mock *userRepositoryMock) GetUserByPhoneCalls() []struct {
	Phone string
} {
	var calls []struct {
		Phone string
	}
	lockuserRepositoryMockGetUserByPhone.RLock()
	calls = mock.calls.GetUserByPhone
	lockuserRepositoryMockGetUserByPhone.RUnlock()
	return calls
}

// UpdateUser calls UpdateUserFunc.
func (mock *userRepositoryMock) UpdateUser(u *entities.User) error {
	if mock.UpdateUserFunc == nil {
		panic("userRepositoryMock.UpdateUserFunc: method is nil but userRepository.UpdateUser was just called")
	}
	callInfo := struct {
		U *entities.User
	}{
		U: u,
	}
	lockuserRepositoryMockUpdateUser.Lock()
	mock.calls.UpdateUser = append(mock.calls.UpdateUser, callInfo)
	lockuserRepositoryMockUpdateUser.Unlock()
	return mock.UpdateUserFunc(u)
}

// UpdateUserCalls gets all the calls that were made to UpdateUser.
// Check the length with:
//     len(mockeduserRepository.UpdateUserCalls())
func (mock *userRepositoryMock) UpdateUserCalls() []struct {
	U *entities.User
} {
	var calls []struct {
		U *entities.User
	}
	lockuserRepositoryMockUpdateUser.RLock()
	calls = mock.calls.UpdateUser
	lockuserRepositoryMockUpdateUser.RUnlock()
	return calls
}

// UseRegCode calls UseRegCodeFunc.
func (mock *userRepositoryMock) UseRegCode(code string) error {
	if mock.UseRegCodeFunc == nil {
		panic("userRepositoryMock.UseRegCodeFunc: method is nil but userRepository.UseRegCode was just called")
	}
	callInfo := struct {
		Code string
	}{
		Code: code,
	}
	lockuserRepositoryMockUseRegCode.Lock()
	mock.calls.UseRegCode = append(mock.calls.UseRegCode, callInfo)
	lockuserRepositoryMockUseRegCode.Unlock()
	return mock.UseRegCodeFunc(code)
}

// UseRegCodeCalls gets all the calls that were made to UseRegCode.
// Check the length with:
//     len(mockeduserRepository.UseRegCodeCalls())
func (mock *userRepositoryMock) UseRegCodeCalls() []struct {
	Code string
} {
	var calls []struct {
		Code string
	}
	lockuserRepositoryMockUseRegCode.RLock()
	calls = mock.calls.UseRegCode
	lockuserRepositoryMockUseRegCode.RUnlock()
	return calls
}

// ValidateRegCode calls ValidateRegCodeFunc.
func (mock *userRepositoryMock) ValidateRegCode(code string) error {
	if mock.ValidateRegCodeFunc == nil {
		panic("userRepositoryMock.ValidateRegCodeFunc: method is nil but userRepository.ValidateRegCode was just called")
	}
	callInfo := struct {
		Code string
	}{
		Code: code,
	}
	lockuserRepositoryMockValidateRegCode.Lock()
	mock.calls.ValidateRegCode = append(mock.calls.ValidateRegCode, callInfo)
	lockuserRepositoryMockValidateRegCode.Unlock()
	return mock.ValidateRegCodeFunc(code)
}

// ValidateRegCodeCalls gets all the calls that were made to ValidateRegCode.
// Check the length with:
//     len(mockeduserRepository.ValidateRegCodeCalls())
func (mock *userRepositoryMock) ValidateRegCodeCalls() []struct {
	Code string
} {
	var calls []struct {
		Code string
	}
	lockuserRepositoryMockValidateRegCode.RLock()
	calls = mock.calls.ValidateRegCode
	lockuserRepositoryMockValidateRegCode.RUnlock()
	return calls
}

var (
	lockServiceMockAddFamilyMember        sync.RWMutex
	lockServiceMockDeleteFamilyMember     sync.RWMutex
	lockServiceMockListFamilyMembers      sync.RWMutex
	lockServiceMockRegister               sync.RWMutex
	lockServiceMockUserByID               sync.RWMutex
	lockServiceMockUserByPhoneAndPassword sync.RWMutex
)

// Ensure, that ServiceMock does implement Service.
// If this is not the case, regenerate this file with moq.
var _ Service = &ServiceMock{}

// ServiceMock is a mock implementation of Service.
//
//     func TestSomethingThatUsesService(t *testing.T) {
//
//         // make and configure a mocked Service
//         mockedService := &ServiceMock{
//             AddFamilyMemberFunc: func(ctx context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error) {
// 	               panic("mock out the AddFamilyMember method")
//             },
//             DeleteFamilyMemberFunc: func(ctx context.Context, r *dto.DeleteFamilyMemberRequest) error {
// 	               panic("mock out the DeleteFamilyMember method")
//             },
//             ListFamilyMembersFunc: func(ctx context.Context, id uint) (*dto.ListFamilyMembersResponse, error) {
// 	               panic("mock out the ListFamilyMembers method")
//             },
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
//         // use mockedService in code that requires Service
//         // and then make assertions.
//
//     }
type ServiceMock struct {
	// AddFamilyMemberFunc mocks the AddFamilyMember method.
	AddFamilyMemberFunc func(ctx context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error)

	// DeleteFamilyMemberFunc mocks the DeleteFamilyMember method.
	DeleteFamilyMemberFunc func(ctx context.Context, r *dto.DeleteFamilyMemberRequest) error

	// ListFamilyMembersFunc mocks the ListFamilyMembers method.
	ListFamilyMembersFunc func(ctx context.Context, id uint) (*dto.ListFamilyMembersResponse, error)

	// RegisterFunc mocks the Register method.
	RegisterFunc func(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)

	// UserByIDFunc mocks the UserByID method.
	UserByIDFunc func(ctx context.Context, id uint) (*dto.UserByIDResponse, error)

	// UserByPhoneAndPasswordFunc mocks the UserByPhoneAndPassword method.
	UserByPhoneAndPasswordFunc func(ctx context.Context, phone string, password string) (*dto.UserAuthResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddFamilyMember holds details about calls to the AddFamilyMember method.
		AddFamilyMember []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *dto.AddFamilyMemberRequest
		}
		// DeleteFamilyMember holds details about calls to the DeleteFamilyMember method.
		DeleteFamilyMember []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// R is the r argument value.
			R *dto.DeleteFamilyMemberRequest
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

// AddFamilyMember calls AddFamilyMemberFunc.
func (mock *ServiceMock) AddFamilyMember(ctx context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error) {
	if mock.AddFamilyMemberFunc == nil {
		panic("ServiceMock.AddFamilyMemberFunc: method is nil but Service.AddFamilyMember was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *dto.AddFamilyMemberRequest
	}{
		Ctx: ctx,
		R:   r,
	}
	lockServiceMockAddFamilyMember.Lock()
	mock.calls.AddFamilyMember = append(mock.calls.AddFamilyMember, callInfo)
	lockServiceMockAddFamilyMember.Unlock()
	return mock.AddFamilyMemberFunc(ctx, r)
}

// AddFamilyMemberCalls gets all the calls that were made to AddFamilyMember.
// Check the length with:
//     len(mockedService.AddFamilyMemberCalls())
func (mock *ServiceMock) AddFamilyMemberCalls() []struct {
	Ctx context.Context
	R   *dto.AddFamilyMemberRequest
} {
	var calls []struct {
		Ctx context.Context
		R   *dto.AddFamilyMemberRequest
	}
	lockServiceMockAddFamilyMember.RLock()
	calls = mock.calls.AddFamilyMember
	lockServiceMockAddFamilyMember.RUnlock()
	return calls
}

// DeleteFamilyMember calls DeleteFamilyMemberFunc.
func (mock *ServiceMock) DeleteFamilyMember(ctx context.Context, r *dto.DeleteFamilyMemberRequest) error {
	if mock.DeleteFamilyMemberFunc == nil {
		panic("ServiceMock.DeleteFamilyMemberFunc: method is nil but Service.DeleteFamilyMember was just called")
	}
	callInfo := struct {
		Ctx context.Context
		R   *dto.DeleteFamilyMemberRequest
	}{
		Ctx: ctx,
		R:   r,
	}
	lockServiceMockDeleteFamilyMember.Lock()
	mock.calls.DeleteFamilyMember = append(mock.calls.DeleteFamilyMember, callInfo)
	lockServiceMockDeleteFamilyMember.Unlock()
	return mock.DeleteFamilyMemberFunc(ctx, r)
}

// DeleteFamilyMemberCalls gets all the calls that were made to DeleteFamilyMember.
// Check the length with:
//     len(mockedService.DeleteFamilyMemberCalls())
func (mock *ServiceMock) DeleteFamilyMemberCalls() []struct {
	Ctx context.Context
	R   *dto.DeleteFamilyMemberRequest
} {
	var calls []struct {
		Ctx context.Context
		R   *dto.DeleteFamilyMemberRequest
	}
	lockServiceMockDeleteFamilyMember.RLock()
	calls = mock.calls.DeleteFamilyMember
	lockServiceMockDeleteFamilyMember.RUnlock()
	return calls
}

// ListFamilyMembers calls ListFamilyMembersFunc.
func (mock *ServiceMock) ListFamilyMembers(ctx context.Context, id uint) (*dto.ListFamilyMembersResponse, error) {
	if mock.ListFamilyMembersFunc == nil {
		panic("ServiceMock.ListFamilyMembersFunc: method is nil but Service.ListFamilyMembers was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  uint
	}{
		Ctx: ctx,
		ID:  id,
	}
	lockServiceMockListFamilyMembers.Lock()
	mock.calls.ListFamilyMembers = append(mock.calls.ListFamilyMembers, callInfo)
	lockServiceMockListFamilyMembers.Unlock()
	return mock.ListFamilyMembersFunc(ctx, id)
}

// ListFamilyMembersCalls gets all the calls that were made to ListFamilyMembers.
// Check the length with:
//     len(mockedService.ListFamilyMembersCalls())
func (mock *ServiceMock) ListFamilyMembersCalls() []struct {
	Ctx context.Context
	ID  uint
} {
	var calls []struct {
		Ctx context.Context
		ID  uint
	}
	lockServiceMockListFamilyMembers.RLock()
	calls = mock.calls.ListFamilyMembers
	lockServiceMockListFamilyMembers.RUnlock()
	return calls
}

// Register calls RegisterFunc.
func (mock *ServiceMock) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	if mock.RegisterFunc == nil {
		panic("ServiceMock.RegisterFunc: method is nil but Service.Register was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Req *dto.UserRegisterRequest
	}{
		Ctx: ctx,
		Req: req,
	}
	lockServiceMockRegister.Lock()
	mock.calls.Register = append(mock.calls.Register, callInfo)
	lockServiceMockRegister.Unlock()
	return mock.RegisterFunc(ctx, req)
}

// RegisterCalls gets all the calls that were made to Register.
// Check the length with:
//     len(mockedService.RegisterCalls())
func (mock *ServiceMock) RegisterCalls() []struct {
	Ctx context.Context
	Req *dto.UserRegisterRequest
} {
	var calls []struct {
		Ctx context.Context
		Req *dto.UserRegisterRequest
	}
	lockServiceMockRegister.RLock()
	calls = mock.calls.Register
	lockServiceMockRegister.RUnlock()
	return calls
}

// UserByID calls UserByIDFunc.
func (mock *ServiceMock) UserByID(ctx context.Context, id uint) (*dto.UserByIDResponse, error) {
	if mock.UserByIDFunc == nil {
		panic("ServiceMock.UserByIDFunc: method is nil but Service.UserByID was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  uint
	}{
		Ctx: ctx,
		ID:  id,
	}
	lockServiceMockUserByID.Lock()
	mock.calls.UserByID = append(mock.calls.UserByID, callInfo)
	lockServiceMockUserByID.Unlock()
	return mock.UserByIDFunc(ctx, id)
}

// UserByIDCalls gets all the calls that were made to UserByID.
// Check the length with:
//     len(mockedService.UserByIDCalls())
func (mock *ServiceMock) UserByIDCalls() []struct {
	Ctx context.Context
	ID  uint
} {
	var calls []struct {
		Ctx context.Context
		ID  uint
	}
	lockServiceMockUserByID.RLock()
	calls = mock.calls.UserByID
	lockServiceMockUserByID.RUnlock()
	return calls
}

// UserByPhoneAndPassword calls UserByPhoneAndPasswordFunc.
func (mock *ServiceMock) UserByPhoneAndPassword(ctx context.Context, phone string, password string) (*dto.UserAuthResponse, error) {
	if mock.UserByPhoneAndPasswordFunc == nil {
		panic("ServiceMock.UserByPhoneAndPasswordFunc: method is nil but Service.UserByPhoneAndPassword was just called")
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
	lockServiceMockUserByPhoneAndPassword.Lock()
	mock.calls.UserByPhoneAndPassword = append(mock.calls.UserByPhoneAndPassword, callInfo)
	lockServiceMockUserByPhoneAndPassword.Unlock()
	return mock.UserByPhoneAndPasswordFunc(ctx, phone, password)
}

// UserByPhoneAndPasswordCalls gets all the calls that were made to UserByPhoneAndPassword.
// Check the length with:
//     len(mockedService.UserByPhoneAndPasswordCalls())
func (mock *ServiceMock) UserByPhoneAndPasswordCalls() []struct {
	Ctx      context.Context
	Phone    string
	Password string
} {
	var calls []struct {
		Ctx      context.Context
		Phone    string
		Password string
	}
	lockServiceMockUserByPhoneAndPassword.RLock()
	calls = mock.calls.UserByPhoneAndPassword
	lockServiceMockUserByPhoneAndPassword.RUnlock()
	return calls
}
