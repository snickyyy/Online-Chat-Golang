package api_errors

import "errors"

var ErrPasswordsDontMatch = errors.New("passwords don't match")
var ErrInvalidToken = errors.New("token is invalid")
var ErrUserAlreadyExists = errors.New("account with username or email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAlreadyLoggedIn = errors.New("already logged in")
var ErrNotLoggedIn = errors.New("not logged in")
var ErrPasswordLight = errors.New("the password is too light")
var ErrUnauthorized = errors.New("unauthorized")
var ErrSamePassword = errors.New("passwords must be different")
