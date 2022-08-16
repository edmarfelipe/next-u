# Identity

Authentication, Authorization and User Registration

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