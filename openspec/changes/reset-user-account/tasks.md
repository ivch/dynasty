## 1. Service Layer

- [x] 1.1 Implement `AdminResetApartment(ctx context.Context, adminID, buildingID, apartmentNumber uint) error` on `users.Service`: fetch admin user, verify `Role == 1`, call `FindUserByApartment` (return `UserNotFound` if nil), build placeholder, call repo transactional reset

## 2. Error Codes

- [x] 2.1 Add `InsufficientPermissionsCode` error to `common/errs/errs.go` with EN/RU/UA translations
- [x] 2.2 Add `NoRegCodesAvailableCode` error to `common/errs/errs.go` with EN/RU/UA translations

## 3. Repository Layer

- [x] 3.1 Add `AdminResetApartment(targetID uint, placeholder *User) error` to `UserRepository` interface in `server/handlers/users/service.go`
- [x] 3.2 Implement the repo method in `server/handlers/users/repo/repo.go` as a DB transaction: get unused reg_code (return `NoRegCodesAvailable` if none), hard-delete target user by ID, insert placeholder user, mark reg_code as used

## 4. Transport Layer

- [x] 4.1 Add `AdminResetApartment(ctx context.Context, adminID, buildingID, apartmentNumber uint) error` to the `UsersService` interface in `server/handlers/users/transport/http.go`
- [x] 4.2 Add `adminResetApartmentRequest` DTO in `server/handlers/users/transport/dto.go` with `building_id` and `apartment_number` fields
- [x] 4.3 Implement `AdminResetApartment` HTTP handler in `server/handlers/users/transport/http.go`: decode body, read `X-Auth-User` header for admin ID, validate required fields, call service, return `{}` on success
- [x] 4.4 Register route `POST /v1/admin/apartment/reset` in `attachRoutes()`

## 5. Mocks and Tests

- [x] 5.1 Regenerate mocks with `make gen` after interface changes
- [x] 5.2 Write service tests for `AdminResetApartment`: success, non-admin rejected, admin not found, apartment not found
- [x] 5.3 Write transport tests for `AdminResetApartment` handler: success, missing fields, non-admin, missing admin header, no reg codes

## 6. Validation

- [x] 6.1 Run `make lint` and fix any lint errors
- [x] 6.2 Run `make test` and confirm all tests pass
