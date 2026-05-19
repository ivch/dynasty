## Context

The `users` service manages user accounts but has no admin-facing API to return an apartment slot to its pre-registration state. When a master account needs to be wiped (compromised account, tenant change, etc.), admins currently have no supported path.

The apartment slot pattern already exists: `PredefinedUserRole = 5` and `registerPredefinedUser` are in the codebase; this feature creates the admin-side counterpart that resets a slot back to that state.

Role 1 (`admin`) is defined in the `user_roles` table. The requesting user's ID arrives via the `X-Auth-User` header (set by the `IDCtx` middleware). There is no general-purpose role-check middleware; role verification is performed in the service layer.

## Goals / Non-Goals

**Goals:**
- New `POST /v1/admin/apartment/reset` endpoint, callable only by admin-role users, identified by `building_id` + `apartment_number` in the request body
- Looks up the master user for the apartment via the existing `FindUserByApartment` repo method
- Hard-deletes the master user record; DB cascade handles sessions, requests, password recovery, and family members automatically
- Creates a placeholder user (role 5, `active = false`) for the apartment slot with a synthetic phone `{building_id}{entry_id}{apartment_number}` (e.g. `"11123"`), empty email, real apartment number, and a fresh reg_code from the `reg_codes` table
- All steps wrapped in a single DB transaction; rolls back if no unused reg_code is available

**Non-Goals:**
- User lookup by internal user ID (apartment coordinates are the sole identifier)
- Optional temporary password on reset
- Bulk reset of multiple accounts
- Audit log / admin action history
- Fixing the existing password-change session-invalidation TODOs in `service.go` (separate concern)

## Decisions

### D1: Apartment coordinates as identifier, not user ID

**Decision**: The endpoint accepts `building_id` + `apartment_number` in the body; `FindUserByApartment` locates the master user. No `{id}` URL parameter.

**Rationale**: Admins operate on apartment slots, not on internal user IDs they may not know. `FindUserByApartment` already exists and enforces `parent_id IS NULL` (master accounts only). The URL `/v1/admin/apartment/reset` reflects what is actually being addressed.

**Alternative considered**: `POST /v1/admin/user/{id}/reset` with body as safety confirmation. Rejected — requires admin to know internal IDs; adds a cross-check with no extra safety benefit since building+apt already uniquely identifies the master record.

### D2: Extend `users.Service`, not a new `admin` handler package

**Decision**: Add `AdminResetApartment` to the existing `users.Service` and transport.

**Rationale**: The operation is a user/apartment mutation. The neomonolith pattern groups by domain (users), not by actor (admin). A new package for a single endpoint is unnecessary indirection.

**Alternative considered**: Separate `server/handlers/admin/` package. Rejected — adds a package and an extra interface for no architectural gain at this scope.

### D3: Transaction lives in the repo layer

**Decision**: A new `AdminResetApartment(buildingID, entryID, apartmentNumber uint) error` repo method wraps all DB steps in a single `db.Transaction(...)`, consistent with the existing `ResetPassword` pattern.

**Rationale**: Keeps transaction boundaries at the data layer. The service layer calls one repo method; it either succeeds atomically or returns an error.

### D4: No `sessionDeleter` interface needed

**Decision**: Session cleanup is handled entirely by the existing `ON DELETE CASCADE` constraint on `sessions.user_id → users.id`.

**Rationale**: When the master user is hard-deleted, PostgreSQL cascades the delete to sessions, requests, password_recovery, and family member rows automatically. An explicit `DeleteSessionByUserID` call would be redundant and would require a cross-package interface.

### D5: Role check in the service layer

**Decision**: `AdminResetApartment` on the service fetches the requesting admin's user record and verifies `Role == 1` before proceeding.

**Rationale**: No general-purpose role-check middleware exists; adding one is a larger cross-cutting change. Consistent with how `DeleteFamilyMember` checks ownership in the service layer.

**Alternative considered**: New `RoleCheck` middleware in `server/middlewares/`. Deferred — appropriate when multiple endpoints share the same role gate.

## Risks / Trade-offs

- **Apartment not found**: If no master user exists for the given building+apt, `FindUserByApartment` returns `nil`. The service should return `UserNotFound` and not proceed. → Handle in service before calling repo.
- **reg_code exhaustion**: `GetRegCode` returns an error when no unused codes exist. The repo transaction rolls back cleanly. → Surfaced as `NoRegCodesAvailable` error to the caller.
- **Synthetic phone collisions**: Phone `"11123"` is deterministic per building/entry/apt combination. If `phone` had a UNIQUE index this would be safe by construction; currently there is only a regular index. A future reset of the same apartment (after a second tenant) would produce the same synthetic phone, which is fine since the previous placeholder was deleted. → No issue.
- **Admin can reset another admin's apartment**: No guard. → Out of scope; noted for future role-hierarchy enforcement.

## Migration Plan

1. Deploy; new endpoint only, no schema changes.
2. Rollback: remove the new route from the router — no state to undo.
