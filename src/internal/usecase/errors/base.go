package api_errors

import "errors"

var ErrInvalidBody = errors.New("the request body is filled incorrectly")
var ErrInvalidData = errors.New("the request data is filled incorrectly")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidPage = errors.New("invalid page")
var ErrNotEnoughPermissions = errors.New("not enough permissions")
