## ADDED Requirements

### Requirement: Admin can reset an apartment slot
The system SHALL provide an endpoint `POST /v1/admin/apartment/reset` that allows a user with role 1 (admin) to fully wipe the master user for a given apartment and restore the slot to a registerable placeholder state. The request body SHALL contain `building_id` and `apartment_number`. The requesting user's ID SHALL be read from the `X-Auth-User` header. All database operations SHALL execute within a single transaction.

#### Scenario: Successful reset
- **WHEN** an admin sends `POST /v1/admin/apartment/reset` with a valid `building_id` and `apartment_number` that has an existing master user
- **THEN** the master user record is hard-deleted from the database
- **AND** all associated records are deleted via cascade (sessions, requests, password recovery records, family members)
- **AND** a new placeholder user is created for the same apartment with role 5 (PredefinedUserRole), `active = false`, an empty email, the original apartment number, and a phone set to the concatenation of `building_id`, `entry_id`, and `apartment_number` (e.g. `"11123"` for building 1, entry 1, apartment 123)
- **AND** the placeholder's `reg_code` is set to an unused code from the `reg_codes` table, which is then marked as used
- **AND** the response is HTTP 200 with an empty JSON object `{}`

#### Scenario: No master user found for the apartment
- **WHEN** an admin sends `POST /v1/admin/apartment/reset` with a `building_id` and `apartment_number` that has no master user record
- **THEN** no database changes are made
- **AND** the system returns HTTP 500 with the `UserNotFound` error

#### Scenario: No unused reg_code available
- **WHEN** an admin sends `POST /v1/admin/apartment/reset` but the `reg_codes` table has no unused codes
- **THEN** the transaction is rolled back (no user deleted, no placeholder created)
- **AND** the system returns HTTP 500 with the `NoRegCodesAvailable` error

#### Scenario: Non-admin user attempts reset
- **WHEN** a user with role other than 1 sends `POST /v1/admin/apartment/reset`
- **THEN** the system returns HTTP 403 with the `InsufficientPermissions` error

#### Scenario: Requesting user not found (invalid X-Auth-User)
- **WHEN** the `X-Auth-User` header contains a user ID that does not exist in the database
- **THEN** the system returns HTTP 500 with the `UserNotFound` error

#### Scenario: Request body missing required fields
- **WHEN** the request body is missing `building_id` or `apartment_number`
- **THEN** the system returns HTTP 400 with the `BadRequest` error
