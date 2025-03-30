package api_errors

import "errors"

var ErrChatAlreadyExists = errors.New("chat already exists")
var ErrChatNotFound = errors.New("chat not found")
var ErrNotEnoughPermissionsForInviting = errors.New("not enough permissions for inviting")
var ErrUserAlreadyInChat = errors.New("user already in chat")
var ErrInviterNotInChat = errors.New("inviter not in chat")
