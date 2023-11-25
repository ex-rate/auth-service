package errors

import "errors"

var ErrUserNotExists = errors.New("user not exists")
var ErrPasswordIncorrect = errors.New("password is incorrect")

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrPhoneAlreadyExists = errors.New("phone already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")
