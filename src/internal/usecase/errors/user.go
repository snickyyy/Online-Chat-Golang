package api_errors

import "errors"

var ErrProfileNotFound = errors.New("profile not found")
var ErrNeedLoginForChangeProfile = errors.New("your need to login for change profile")
var ErrUserNotFound = errors.New("user not found")
var ErrChangePasswordTooOften = errors.New("you can't change passwords too often")
var ErrInvalidCode = errors.New("invalid code")
