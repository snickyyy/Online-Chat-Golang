package usecase_errors

type INotFoundError interface {
	NotFountError() string
}

type IPermissionError interface {
	PermissionError() string
}

type IAlreadyExistsError interface {
	AlreadyExistsError() string
}

type IUnauthorizedError interface {
	UnauthorizedError() string
}

type IBadRequestError interface {
	BadRequestError() string
}

type NotFoundError struct {
	Msg string
}

func (e NotFoundError) NotFountError() string {
	return e.Msg
}
func (e NotFoundError) Error() string {
	return e.Msg
}

type PermissionError struct {
	Msg string
}

func (e PermissionError) PermissionError() string {
	return e.Msg
}
func (e PermissionError) Error() string {
	return e.Msg
}

type AlreadyExistsError struct {
	Msg string
}

func (e AlreadyExistsError) AlreadyExistsError() string {
	return e.Msg
}
func (e AlreadyExistsError) Error() string {
	return e.Msg
}

type UnauthorizedError struct {
	Msg string
}

func (e UnauthorizedError) UnauthorizedError() string {
	return e.Msg
}
func (e UnauthorizedError) Error() string {
	return e.Msg
}

type BadRequestError struct {
	Msg string
}

func (e BadRequestError) BadRequestError() string {
	return e.Msg
}
func (e BadRequestError) Error() string {
	return e.Msg
}
