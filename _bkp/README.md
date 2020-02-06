## Build
`docker-compose build`

## Run
`docker-compose build`
`docker-compose up`

## API 

Login user
```
POST /auth/v1/login

{"phone":"380671234567", "password":"qwerty"}

#success
HTTP/1.1 200 OK

{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTgsIk5hbWUiOiJKb2huIERvZSIsIlJvbGUiOjQsImF1ZCI6ImR5bmFwcCIsImV4cCI6MTU3NzM2MzU2MiwiaWF0IjoxNTc3Mjc3MTYyLCJpc3MiOiJhdXRoLmR5bmFwcCJ9.asGix3XEgR0CwlRYZYyEYyqPcptPp04OjYZojlYBpyI",
    "refresh_token": "395a6fac-3aee-4891-8ed9-ec5546b8777c"
}

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```

Refresh access token
```
POST /auth/v1/refresh

{"token":"395a6fac-3aee-4891-8ed9-ec5546b8777c"}

#success
HTTP/1.1 200 OK

{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTgsIk5hbWUiOiJKb2huIERvZSIsIlJvbGUiOjQsImF1ZCI6ImR5bmFwcCIsImV4cCI6MTU3NzM2MzU2MiwiaWF0IjoxNTc3Mjc3MTYyLCJpc3MiOiJhdXRoLmR5bmFwcCJ9.asGix3XEgR0CwlRYZYyEYyqPcptPp04OjYZojlYBpyI",
    "refresh_token": "395a6fac-3aee-4891-8ed9-ec5546b8777c"
}

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```

Create User
```
POST /users/v1/register

{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"phone":"380671234567","password":"qwerty", "building_id": 2, "code": "as12da"}

#success
HTTP/1.1 200 OK

{
    "id": 51,
    "email": "test2@test.com"
}

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```

Get User info (for given access token)
```
Header: Authorization Bearer <access-token>
GET /users/v1/user
#success
HTTP/1.1 200 OK

{
    "id": 18,
    "first_name": "John",
    "last_name": "Doe",
    "phone": "380671234567",
    "email": "test@test.com"
}

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```

Create request
```
Header: Authorization Bearer <access-token>
POST /requests/v1/request
{"type":"taxi", "time":1580218937, "description":"blabla comment"}

#success
HTTP/1.1 200 OK

{
    "id": "1",
}

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```

Update request
```
Header: Authorization Bearer <access-token>
PUT /requests/v1/request/<id>
{"type":"taxi", "time":1580218937, "description":"blabla comment", "status":"new"}

#success
HTTP/1.1 200 OK

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```

Get requests lists for user
```
Header: Authorization Bearer <access-token>
GET /requests/v1/my?offset=<uint>&limit=<uint>

#success
HTTP/1.1 200 OK

[
    {
        "id": 2,
        "type": "taxi",
        "user_id": 1,
        "time": 1580218937,
        "description": "blabla comment",
        "status": "new"
    },
    {
        "id": 3,
        "type": "guest",
        "user_id": 1,
        "time": 1580218937,
        "description": "blabla comment",
        "status": "done"
    },
]

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```
Get single request
```
Header: Authorization Bearer <access-token>
GET /requests/v1/request/<id>

#success
HTTP/1.1 200 OK

{
    "id": 2,
    "type": "taxi",
    "user_id": 1,
    "time": 1580218937,
    "description": "blabla comment",
    "status": "new"
}

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```
Delete request
```
Header: Authorization Bearer <access-token>
DELETE /requests/v1/request/<id>

#success
HTTP/1.1 200 OK

# Any internal error
HTTP/1.1 500 Internal Server Error
{"error":"error text"}
```