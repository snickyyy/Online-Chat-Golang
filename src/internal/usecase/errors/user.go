package api_errors

import "errors"

var ErrProfileNotFound = errors.New("profile not found")
var ErrNeedLoginForChangeProfile = errors.New("your need to login for change profile")
