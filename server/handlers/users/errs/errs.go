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
	emailEmptyCode
	emailInvalidCode
	regCodeLengthCode
	regCodeInvalidCode
	userNotFoundCode
	masterAccountExistsCode
	regCodeWrongCode
	familyMemberAlreadyRegisteredCode
	familyMemberWrongAddressCode
)

var (
	Generic             = errors.New("something went wrong")
	EmptyUserID         = errors.New("empty user id")
	BadUserID           = errors.New("bad user id")
	BadRequest          = errors.New("failed to decode request")
	PasswordTooShort    = errors.New("password should be at least 6 characters")
	PhoneWrongLength    = errors.New("phone should min 12, max 13 character")
	PhoneWrongChars     = errors.New("phone should contain only numeric characters")
	FNameLength         = errors.New("first name is required")
	LNameLength         = errors.New("last name is required")
	BuildingEmpty       = errors.New("building is required")
	EntryEmpty          = errors.New("entry is required")
	ApartmentEmpty      = errors.New("apartment is required")
	EmailEmpty          = errors.New("email is required")
	EmailInvalid        = errors.New("email is invalid")
	RegCodeLength       = errors.New("reg code should be at least 5 characters")
	RegCodeInvalid      = errors.New("provided wrong reg code")
	UserNotFound        = errors.New("user not found")
	MasterAccountExists = errors.New("master account for this apartment already exists")
	RegCodeWrong        = errors.New("wrong reg code provided")
	// errInvalidRequest                = errors.New("request validation error")
	// errFamilyMembersLimitExceeded    = errors.New("family members limit exceeded")
	FamilyMemberAlreadyRegistered = errors.New("family member already registered")
	FamilyMemberWrongAddress      = errors.New("family member provided wrong address")

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
	}
)

func Code(err error) uint {
	if _, ok := codes[err]; ok {
		return codes[err]
	}

	return genericCode
}
