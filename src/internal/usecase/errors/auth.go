package api_errors

import "errors"

var ErrPasswordsDontMatch = errors.New("passwords don't match")
var ErrInvalidToken = errors.New("token is invalid")

