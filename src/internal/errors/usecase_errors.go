package errors

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
	BaseError
}

func (e *NotFoundError) NotFountError() string {
	return e.Msg
}
func (e *NotFoundError) Error() string {
	return e.Msg
}

type PermissionError struct {
	BaseError
}

func (e *PermissionError) PermissionError() string {
	return e.Msg
}
func (e *PermissionError) Error() string {
	return e.Msg
}

type AlreadyExistsError struct {
	BaseError
}

func (e *AlreadyExistsError) AlreadyExistsError() string {
	return e.Msg
}
func (e *AlreadyExistsError) Error() string {
	return e.Msg
}

type UnauthorizedError struct {
	BaseError
}

func (e *UnauthorizedError) UnauthorizedError() string {
	return e.Msg
}
func (e *UnauthorizedError) Error() string {
	return e.Msg
}

type BadRequestError struct {
	BaseError
}

func (e *BadRequestError) BadRequestError() string {
	return e.Msg
}
func (e *BadRequestError) Error() string {
	return e.Msg
}
