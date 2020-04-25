# TODO List

## Registration
- [x] add section number
- [x] check building and entry id
- [x] add master account
- [x] add ability to create accounts under master account
- [x] check phone number length (38 problem)

## Mobile client
- [x] different error codes

## Security service
- [x] get requests list
- [x] mark request as completed
- [x] get neighbor info for request

## Guard UI
- [x] add FE configuration from BE

## Family members
- [x] cancel registration if there is a user for this building/app
- [x] limit family members to 10
- [x] on login check if user is active
- [ ] on token validation check if user is active
- [x] update schema on prod and cleanup db
- [x] set all master accs to true

## Requests
- [ ] add photo to requests
- [ ] on guard ui remove minutes and second on view

## Errors
- [ ] rethink entities.errors


`alter table users alter column parent_id set default null;`