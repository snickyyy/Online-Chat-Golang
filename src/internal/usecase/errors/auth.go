package api_errors

import "errors"

var ErrPasswordsDontMatch = errors.New("passwords don't match")
var ErrInvalidToken = errors.New("token is invalid")
var ErrUserAlreadyExists = errors.New("account with this fields already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
