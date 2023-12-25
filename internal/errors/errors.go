package errors

import "errors"

var ErrUserNotExists = errors.New("user not exists")
var ErrIncorrectPassword = errors.New("password is incorrect")

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrPhoneAlreadyExists = errors.New("phone already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")

var ErrInvalidToken = errors.New("token is invalid")

// ErrTokenNotExists используется, если в базе не найден указанный рефреш токен
var ErrTokenNotExists = errors.New("token was not found in the database")

// ErrInvalidUsername используется при проверке access и refresh токена.
// Возвращается в случае, если usernames в токенах не совпадают
var ErrInvalidUsername = errors.New("usernames doesn't match")

var ErrNotAuthorized = errors.New("user is not authorized")
