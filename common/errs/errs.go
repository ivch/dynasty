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
	badRecoveryCode
	emptyPasswordCode
	recoveryCodeOutdatedCode
	aptNumberIsTooBigCode
	emailAlreadyExistsCode
)

type SvcError struct {
	Code int
	Err  error
	Ru   string
	Ua   string
}

func New(code int, en, ru, ua string) SvcError {
	return SvcError{
		Code: code,
		Err:  errors.New(en),
		Ru:   ru,
		Ua:   ua,
	}
}

func (s SvcError) Error() string {
	return s.Err.Error()
}

var (
	Generic                       = New(genericCode, "something went wrong", "что-то пошло не так :(", "щось пішло не так :(")
	EmptyUserID                   = New(emptyUserIDCode, "empty user id", "не указан ID", "не вказано ID")
	BadUserID                     = New(badUserIDCode, "bad user id", "неправильный ID", "невірний iD")
	BadRequest                    = New(badRequestCode, "failed to decode request", "неправильный запрос", "невірний запит")
	PasswordTooShort              = New(passwordTooShortCode, "password should be at least 6 characters", "пароль должен быть не короче 6 символов", "пароль має бути більше 6 символів")
	PhoneWrongLength              = New(phoneWrongLengthCode, "phone should min 12, max 13 character", "телефон может быть 12 или 13 символов", "телефон може бути 12 або 13 символів")
	PhoneWrongChars               = New(phoneWrongCharsCode, "phone should contain only numeric characters", "телефон должен содержать только цифры", "телефон має складатися лише з цифр")
	FNameLength                   = New(fNameLengthCode, "first name is required", "имя обязательно", "ім'я обов'язкове")
	LNameLength                   = New(lNameLengthCode, "last name is required", "фамилия обязательна", "прізвище обов'язкове")
	BuildingEmpty                 = New(buildingEmptyCode, "building is required", "номер дома обязателен", "номер дому обов'язковий")
	EntryEmpty                    = New(entryEmptyCode, "entry is required", "номер секции обязательна", "номер секції обов'язковий")
	ApartmentEmpty                = New(apartmentEmptyCode, "apartment is required", "номер квартиры обязателен", "номер помешкання обов'язковий")
	WrongApartment                = New(wrongApartmentCode, "wrong apartment number", "неверный номер квартиры", "невірний номер помешкання")
	EmailEmpty                    = New(emailEmptyCode, "email is required", "email обязателен", "email обов'язковий")
	EmailInvalid                  = New(emailInvalidCode, "email is invalid", "неправильный email", "невірний email")
	RegCodeLength                 = New(regCodeLengthCode, "reg code should be at least 5 characters", "код регистрации должен быть не меньше 5 символов", "код реєстрації має бути більше 5 символів")
	RegCodeInvalid                = New(regCodeInvalidCode, "provided invalid reg code", "указан неправильный код", "вказано невірний код")
	UserNotFound                  = New(userNotFoundCode, "user not found", "пользователь не найден", "користувача не знайдено")
	MasterAccountExists           = New(masterAccountExistsCode, "master account for this apartment already exists", "основной аккаунт для этой квартиры уже существует", "основний аккаунт для цього помешкання вже існує")
	RegCodeWrong                  = New(regCodeWrongCode, "wrong reg code provided", "указан неверный код", "вказано невірний код")
	FamilyMembersLimitExceeded    = New(familyMembersLimitExceededCode, "family members limit exceeded", "достигнут максимум членов семьи", "досягнуто максимум членів родини")
	FamilyMemberAlreadyRegistered = New(familyMemberAlreadyRegisteredCode, "family member already registered", "член семьи уже зарегистрирован", "члена родини вже зареєстровано")
	FamilyMemberWrongAddress      = New(familyMemberWrongAddressCode, "family member provided wrong address", "указан неправильный семейный адрес", "вказано невірну адресу родини")
	FamilyMemberPhoneExists       = New(familyMemberPhoneExistsCode, "provided phone number already exists", "указанный номер телефона уже используется", "вказаний номер телефону вже використовується")
	FamilyMemberWrongOwner        = New(familyMemberWrongOwnerCode, "provided family member has incorrect parent id", "указан неправильный семейний аккаунт", "вказано невірний родинний аккаунт")
	FamilyMemberBadID             = New(familyMemberBadIDCode, "family member bad id", "неправильный ID члена семьи", "невірний ID члена родини")
	FamilyMemberParentMismatch    = New(familyMemberParentMismatchCode, "family member parent mismatch", "неправильный семейный аккаунт", "невірний родинний аккаунт")
	InvalidCredentials            = New(invalidCredentialsCode, "invalid credentials", "неправильный логин или пароль", "неправильні логін чи пароль")
	UserIsInactive                = New(usersIsInactiveCode, "users is inactive", "пользователь отключен", "користувача вимкнено")
	NoSessionToRefresh            = New(noSessionToRefreshCode, "no session to refresh", "сессия не найдена", "сесію не знайдено")
	TokenExpired                  = New(tokenExpiredCode, "token expired", "токен истек", "токен закінчився")
	FailedParsingToken            = New(failedParsingTokenCode, "failed to parse token", "невозможно прочитать токен", "неможливо прочитати токен")
	FailedParsingTokenClaims      = New(failedParsingTokenClaimsCode, "failed to parse token claims", "ошибка чтения токена", "помилка зчитування токену")
	TokenIsInvalid                = New(tokenIsInvalidCode, "token is invalid", "неправильный токен", "невірний токен")
	NoAuthHeader                  = New(noAuthHeaderCode, "no auth header", "нет авторизации", "немає авторизації")
	Unauthorized                  = New(unauthorizedCode, "user is unauthorized", "доступ запрещен", "доступ заборонено")
	WrongRequestType              = New(wrongRequestTypeCode, "wrong request type", "неправильный тип заявки", "неправильний тип заяви")
	WrongRequestStatus            = New(wrongRequestStatusCode, "wrong request status", "неправильный статус заявки", "неправильний статус заяви")
	WrongRequestDate              = New(wrongRequestDateCode, "wrong request date", "неправильная дата заявки", "неправильна дата заяви")
	WrongRequestPlace             = New(wrongRequestPlaceCode, "wrong request place", "неправильное место заявки", "неправильне місце заяви")
	EmptyOffset                   = New(emptyOffsetCode, "empty offset", "empty offset", "empty offset")
	BadOffset                     = New(badOffsetCode, "bad offset", "bad offset", "bad offset")
	EmptyLimit                    = New(emptyLimitCode, "empty limit", "empty limit", "empty limit")
	BadLimit                      = New(badLimitCode, "bad limit", "bad limit", "bad limit")
	LimitTooSmall                 = New(limitTooSmallCode, "limit should be grater then 0", "limit should be grater then 0", "limit should be grater then 0")
	LimitTooBig                   = New(limitTooBigCode, "limit should less or equal 200", "limit should less or equal 200", "limit should less or equal 200")
	NoFile                        = New(noFileCode, "error retrieving the file", "ошибка чтения файла", "помилка зчитування файлу")
	FileWrongType                 = New(fileWrongTypeCode, "wrong filetype", "неверный тип файла", "невірний тип файлу")
	FileIsTooBig                  = New(fileIsTooBigCode, "too big file", "слишком большой файл", "файл завеликий")
	TooMuchFiles                  = New(tooMuchFilesCode, "only allowed 3 images per request", "не более 3х фото на заявку", "не більше 3х фото на заяву")
	RequestPerDayLimitExceeded    = New(requestPerDayLimitExceededCode, "request per day limit exceeded", "достигнут лимит заявок за день", "досягнуто денний ліміт заяв")
	PasswordConfirmMismatch       = New(passwordConfirmMismatchCode, "password confirmation mismatch", "неправильное подтверждение пароля", "неправильне підтвердження пароля")
	PasswordRecoveryLimit         = New(passwordRecoveryLimitCode, "password recoveries limit exceeded", "достигнут лимит попыток восстановлениия пароля", "досягнуто ліміт спроб відновлення паролю")
	BadRecoveryCode               = New(badRecoveryCode, "bad recovery code", "неправильный код восстановления пароля", "невірний код відновлення паролю")
	RecoveryCodeOutdated          = New(recoveryCodeOutdatedCode, "recovery code outdated", "код восстановления пароля истек", "термін дії коду відновлення закінчився")
	EmptyPassword                 = New(emptyPasswordCode, "empty password", "пустой пароль", "пустий пароль")
	AptNumberIsTooBig             = New(aptNumberIsTooBigCode, "apartment number is too big", "слишком большой номер квартиры", "завеликий номер помешкання")
	EmailAlreadyExists            = New(emailAlreadyExistsCode, "provided email already in use", "указанный email уже используется", "вказаний email вже викорістовується")

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
		BadRecoveryCode:               badRecoveryCode,
		EmptyPassword:                 emptyPasswordCode,
		RecoveryCodeOutdated:          recoveryCodeOutdatedCode,
		AptNumberIsTooBig:             aptNumberIsTooBigCode,
		EmailAlreadyExists:            emailAlreadyExistsCode,
	}
)

func Code(err error) uint {
	if _, ok := codes[err]; ok {
		return codes[err]
	}

	return genericCode
}
