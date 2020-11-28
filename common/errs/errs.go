package errs

import (
	"errors"
)

const (
	genericCode    = 0
	badRequestCode = iota + 100
	emptyUserIDCode
	badUserIDCode
	passwordTooShortCode
	phoneWrongLengthCode
	phoneWrongCharsCode
	fNameLengthCode
	lNameLengthCode
	buildingEmptyCode
	entryEmptyCode
	apartmentEmptyCode
	wrongApartmentCode
	emailEmptyCode
	emailInvalidCode
	regCodeLengthCode
	regCodeInvalidCode
	userNotFoundCode
	masterAccountExistsCode
	regCodeWrongCode
	familyMemberAlreadyRegisteredCode
	familyMemberWrongAddressCode
	familyMemberPhoneExistsCode
	familyMembersLimitExceededCode
	familyMemberWrongOwnerCode
	familyMemberBadIDCode
	familyMemberParentMismatchCode
	invalidCredentialsCode
	usersIsInactiveCode
	noSessionToRefreshCode
	tokenExpiredCode
	failedParsingTokenCode
	failedParsingTokenClaimsCode
	tokenIsInvalidCode
	noAuthHeaderCode
	unauthorizedCode
	wrongRequestTypeCode
	wrongRequestStatusCode
	wrongRequestDateCode
	wrongRequestPlaceCode
	emptyOffsetCode
	badOffsetCode
	emptyLimitCode
	badLimitCode
	limitTooSmallCode
	limitTooBigCode
	noFileCode
	fileWrongTypeCode
	fileIsTooBigCode
	tooMuchFilesCode
	requestPerDayLimitExceededCode
	passwordConfirmMismatchCode
	passwordRecoveryLimitCode
)

var (
	Generic                       = errors.New("something went wrong")
	EmptyUserID                   = errors.New("empty user id")
	BadUserID                     = errors.New("bad user id")
	BadRequest                    = errors.New("failed to decode request")
	PasswordTooShort              = errors.New("password should be at least 6 characters")
	PhoneWrongLength              = errors.New("phone should min 12, max 13 character")
	PhoneWrongChars               = errors.New("phone should contain only numeric characters")
	FNameLength                   = errors.New("first name is required")
	LNameLength                   = errors.New("last name is required")
	BuildingEmpty                 = errors.New("building is required")
	EntryEmpty                    = errors.New("entry is required")
	ApartmentEmpty                = errors.New("apartment is required")
	WrongApartment                = errors.New("wrong apartment number")
	EmailEmpty                    = errors.New("email is required")
	EmailInvalid                  = errors.New("email is invalid")
	RegCodeLength                 = errors.New("reg code should be at least 5 characters")
	RegCodeInvalid                = errors.New("provided wrong reg code")
	UserNotFound                  = errors.New("user not found")
	MasterAccountExists           = errors.New("master account for this apartment already exists")
	RegCodeWrong                  = errors.New("wrong reg code provided")
	FamilyMembersLimitExceeded    = errors.New("family members limit exceeded")
	FamilyMemberAlreadyRegistered = errors.New("family member already registered")
	FamilyMemberWrongAddress      = errors.New("family member provided wrong address")
	FamilyMemberPhoneExists       = errors.New("provided phone number already exists")
	FamilyMemberWrongOwner        = errors.New("provided family member has incorrect parent id")
	FamilyMemberBadID             = errors.New("family member bad id")
	FamilyMemberParentMismatch    = errors.New("family member parent mismatch")
	InvalidCredentials            = errors.New("invalid credentials")
	UserIsInactive                = errors.New("users is inactive")
	NoSessionToRefresh            = errors.New("no session to refresh")
	TokenExpired                  = errors.New("token expired")
	FailedParsingToken            = errors.New("failed to parse token")
	FailedParsingTokenClaims      = errors.New("failed to parse token claims")
	TokenIsInvalid                = errors.New("token is invalid")
	NoAuthHeader                  = errors.New("no auth header")
	Unauthorized                  = errors.New("user is unauthorized")
	WrongRequestType              = errors.New("wrong request type")
	WrongRequestStatus            = errors.New("wrong request status")
	WrongRequestDate              = errors.New("wrong request date")
	WrongRequestPlace             = errors.New("wrong request place")
	EmptyOffset                   = errors.New("empty offset")
	BadOffset                     = errors.New("bad offset")
	EmptyLimit                    = errors.New("empty limit")
	BadLimit                      = errors.New("bad limit")
	LimitTooSmall                 = errors.New("limit should be grater then 0")
	LimitTooBig                   = errors.New("limit should less or equal 200")
	NoFile                        = errors.New("error retrieving the file")
	FileWrongType                 = errors.New("wrong filetype")
	FileIsTooBig                  = errors.New("too big file")
	TooMuchFiles                  = errors.New("only allowed 3 images per request")
	RequestPerDayLimitExceeded    = errors.New("request per day limit exceeded")
	PasswordConfirmMismatch       = errors.New("password confirmation mismatch")
	PasswordRecoveryLimit         = errors.New("password recoveries limit exceeded")

	codes = map[error]uint{
		Generic:                       genericCode,
		BadRequest:                    badRequestCode,
		EmptyUserID:                   emptyUserIDCode,
		BadUserID:                     badUserIDCode,
		PasswordTooShort:              passwordTooShortCode,
		PhoneWrongLength:              phoneWrongLengthCode,
		PhoneWrongChars:               phoneWrongCharsCode,
		FNameLength:                   fNameLengthCode,
		LNameLength:                   lNameLengthCode,
		BuildingEmpty:                 buildingEmptyCode,
		EntryEmpty:                    entryEmptyCode,
		ApartmentEmpty:                apartmentEmptyCode,
		EmailEmpty:                    emailEmptyCode,
		EmailInvalid:                  emailInvalidCode,
		RegCodeLength:                 regCodeLengthCode,
		RegCodeInvalid:                regCodeInvalidCode,
		UserNotFound:                  userNotFoundCode,
		MasterAccountExists:           masterAccountExistsCode,
		RegCodeWrong:                  regCodeWrongCode,
		FamilyMemberAlreadyRegistered: familyMemberAlreadyRegisteredCode,
		FamilyMemberWrongAddress:      familyMemberWrongAddressCode,
		FamilyMemberPhoneExists:       familyMemberPhoneExistsCode,
		FamilyMembersLimitExceeded:    familyMembersLimitExceededCode,
		FamilyMemberWrongOwner:        familyMemberWrongOwnerCode,
		FamilyMemberBadID:             familyMemberBadIDCode,
		FamilyMemberParentMismatch:    familyMemberParentMismatchCode,
		InvalidCredentials:            invalidCredentialsCode,
		UserIsInactive:                usersIsInactiveCode,
		NoSessionToRefresh:            noSessionToRefreshCode,
		TokenExpired:                  tokenExpiredCode,
		FailedParsingToken:            failedParsingTokenCode,
		FailedParsingTokenClaims:      failedParsingTokenClaimsCode,
		TokenIsInvalid:                tokenIsInvalidCode,
		NoAuthHeader:                  noAuthHeaderCode,
		Unauthorized:                  unauthorizedCode,
		WrongRequestType:              wrongRequestTypeCode,
		WrongRequestStatus:            wrongRequestStatusCode,
		WrongRequestDate:              wrongRequestDateCode,
		WrongRequestPlace:             wrongRequestPlaceCode,
		EmptyOffset:                   emptyOffsetCode,
		BadOffset:                     badOffsetCode,
		EmptyLimit:                    emptyLimitCode,
		BadLimit:                      badLimitCode,
		LimitTooSmall:                 limitTooSmallCode,
		LimitTooBig:                   limitTooBigCode,
		NoFile:                        noFileCode,
		FileWrongType:                 fileWrongTypeCode,
		FileIsTooBig:                  fileIsTooBigCode,
		TooMuchFiles:                  tooMuchFilesCode,
		WrongApartment:                wrongApartmentCode,
		RequestPerDayLimitExceeded:    requestPerDayLimitExceededCode,
		PasswordConfirmMismatch:       passwordConfirmMismatchCode,
		PasswordRecoveryLimit:         passwordRecoveryLimitCode,
	}
)

func Code(err error) uint {
	if _, ok := codes[err]; ok {
		return codes[err]
	}

	return genericCode
}
