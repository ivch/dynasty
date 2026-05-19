## Why

Administrators currently have no way to return an apartment slot to its pre-registration state. When a resident account needs to be fully wiped (compromised account, tenant change, etc.) there is no admin endpoint to delete the user and restore the apartment slot so a new resident can register.

## What Changes

- Add a new admin HTTP endpoint `POST /v1/admin/apartment/reset` with body `{ building_id, apartment_number }` that:
  - Looks up the master user for that apartment via `FindUserByApartment`
  - Hard-deletes the target user record (cascades to sessions, requests, password recovery, and family members)
  - Creates a new placeholder user for the same apartment with `role = 5` (PredefinedUserRole), a synthetic phone `{building_id}{entry_id}{apartment_number}` (e.g. `"11123"`), empty email, and a fresh reg_code
  - All steps execute in a single DB transaction; rolls back if no unused reg_code is available
- Protect the endpoint with role-based access (admin role only)
- Add new error codes for insufficient permissions and reg_code exhaustion

## Capabilities

### New Capabilities

- `admin-reset-user`: Admin endpoint to fully wipe a user account and recreate the apartment slot as a registerable placeholder

### Modified Capabilities

- none

## Impact

- **New route**: `POST /v1/admin/apartment/reset` in `server/handlers/users/transport/`
- **Service layer**: new `AdminResetUser` method on `users.Service`
- **Repo layer**: new `AdminResetUser` transactional repo method; no new tables; uses existing `users` and `reg_codes` tables
- **Error codes**: 2 new codes in `common/errs/errs.go` (insufficient permissions, no reg codes available)
- **No auth/session dependency**: cascade deletes handle session cleanup automatically
