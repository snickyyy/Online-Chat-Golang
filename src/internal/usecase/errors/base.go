package api_errors

import "errors"

var ErrInvalidBody = errors.New("the request body is filled incorrectly")
