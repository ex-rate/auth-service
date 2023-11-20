package errors

import "errors"

var ErrUserNotExists = errors.New("user not exists")
var ErrPasswordIncorrect = errors.New("password is incorrect")
