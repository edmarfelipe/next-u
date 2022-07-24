# Identity

Authentication, Authorization and User Registration

## Endpoints

### Sign Up
```json
POST /identity/v1/sign-up
{
	"name": "Jon Snow",
	"username": "jon",
	"password": "1234",
	"email": "jonsnow@mail.com"
}
```

### Sign In
```json
POST /identity/v1/sign-in
{
	"username": "jon",
	"password": "1234",
}
```

### Recovery Password
```json
POST /identity/v1/recover-password
{
	"email": "jonsnow@mail.com"
}
```

### Change Password
```json
POST /identity/v1/change-password
{
	"userName": "jon",
	"oldPassword": "oldPass",
	"newPassword": "newPass"
}
```

### Change Password with a token
```json
POST /identity/v1/change-password
{
	"token": "oldPass",
	"newPassword": "newPass"
}
```


### Enable User
```json
POST /identity/v1/enable/{username}
```

### Disable User
```json
POST /identity/v1/disable/{username}
```