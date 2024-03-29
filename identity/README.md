# Identity

Authentication, Authorization and User Registration

## Overview

```mermaid
C4Container
  title Authentication, Authorization and User Registration

  Person(customer, Customer, "A student of the university")
  System_Ext(email_system, "E-Mail System", "SendGrid API")

  Container_Boundary("Identity") {
      Container(backend_api, "API Application", "Java, Docker Container", "Provides user registration and authentication via API", "Golang")
      ContainerDb(database, "Database", "MongoDB", "Stores user registration information.",)
  }

  Rel(customer, backend_api, "Uses", "sync, JSON/HTTPS")
  Rel_Back(database, backend_api, "Reads from and writes to", "sync, TCP/IP")
  UpdateRelStyle(database, backend_api, $offsetX="-50", $offsetY="20")

  Rel(email_system, customer, "Sends e-mails to")
  UpdateRelStyle(email_system, customer, $offsetX="-40", $offsetY="20")

  Rel(backend_api, email_system, "Sends e-mails using", "sync, JSON/HTTPS")
```

## REST API

### Create User

```js
POST /identity/v1/signup
```

Supported attributes:

| Attribute                | Type     | Required  | Description     |
|:-------------------------|:---------|:----------|:----------------|
| `name`                   | string   | Yes       |                 |
| `email`                  | string   | Yes       | Must be unique. |
| `password`               | string   | Yes       |                 |

Example:

```js
curl --request POST \
  --url http://127.0.0.1:3000/identity/v1/signup \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Jon Snow",
	"password": "1234",
	"email": "jon.snow@mail.com"
}'
```

Example response:
```json
{
	"name": "Jon Snow",
	"password": "1234",
	"email": "jon.snow3@mail.com"
}
```

### Authorize

Use this endpoint to authenticate a user with email and password.

```js
POST /identity/v1/authorize
```

Supported attributes:

| Attribute                | Type     | Required  | Description     |
|:-------------------------|:---------|:----------|:----------------|
| `email`                  | string   | Yes       |                 |
| `password`               | string   | Yes       |                 |

Example:

```js
curl --request POST \
  --url http://127.0.0.1:3000/identity/v1/authorize \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "jon.snow@mail.com",
	"password": "1234"
}'
```

Example response:
```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjYxOTA4MTQ0LCJuYW1lIjoiSm9uIFNub3cifQ.T0p2bUJaz840-A7yygRH8tjlvb9r7jmbWRVgASQwBWw"
}
```

### Change Password

```js
POST /identity/v1/password/change
```

Supported attributes:

| Attribute                | Type     | Required  | Description           |
|:-------------------------|:---------|:----------|:----------------------|
| `email`                  | string   | Yes       |                       |
| `oldPassword`            | string   | Yes       |                       |
| `newPassword`            | string   | Yes       |                       |

Example:

```js
curl --request POST \
  --url http://127.0.0.1:3000/identity/v1/password/change \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "jon.snow@mail.com",
	"oldPassword": "oldPassword",
	"newPassword": "newpass"
}'
```


### Recovery Password

```js
POST /identity/v1/password/reset
```

Supported attributes:

| Attribute                | Type     | Required  | Description  |
|:-------------------------|:---------|:----------|:-------------|
| `email`                  | string   | Yes       |              |

Example:

```js
curl --request POST \
  --url http://127.0.0.1:3000/identity/v1/password/reset \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "jon.snow@mail.com"
}'
```

### Change Password with a token
```js
POST /identity/v1/password/change/{token}
```

Example:

```js
curl --request POST \
  --url http://127.0.0.1:3000/identity/v1/password/change/972abbe5-b8e1-4679-add5-4d8217fdc054 \
  --data '{
	"newPassword": "newpass"
  }'
```

### Enable User

```js
PATCH /identity/v1/enable/{id}
```

Example:

```js
curl --request PATCH \
  --url http://127.0.0.1:3000/identity/v1/enable/62fae8ed7463e5388a23bf21 \
  --header 'Authorization: Bearer ...'
```


### Disable User

```js
PATCH /identity/v1/disable/{id}
```

Example:

```js
curl --request PATCH \
  --url http://127.0.0.1:3000/identity/v1/disable/62fae8ed7463e5388a23bf21 \
  --header 'Authorization: Bearer ...'
```

### List users

Example:

```js
curl --request GET \
  --url http://127.0.0.1:3000/identity/v1/ \
  --header 'Authorization: Bearer ...'
```

Example response:

```json
[
  {
    "id": "62fae8ed7463e5388a23bf21",
    "name": "Jon Snow",
    "email": "jon.snow@mail.com",
    "active": true
  }
]
```

### Change user role

Supported attributes:

| Attribute                | Type     | Required  | Description     |
|:-------------------------|:---------|:----------|:----------------|
| `role`                   | string   | Yes       |                 |


Example:

```js
curl --request GET \
  --url http://127.0.0.1:3000/identity/v1/change-role/6313f8b2602913651f9bc8a8 \
  --header 'Authorization: Bearer ...'
  --header 'Content-Type: application/json' \
  --data '{
	"role": "admin"
}'
```

Example response:

```json
{
	"id": "6313f8b2602913651f9bc8a8",
	"name": "Jon Snow",
	"email": "jon.snow4@mail.com",
	"role": "admin",
	"active": true
}
```