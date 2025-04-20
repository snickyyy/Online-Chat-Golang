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
