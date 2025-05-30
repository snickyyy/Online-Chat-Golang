// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "libs/src/internal/dto"

	mock "github.com/stretchr/testify/mock"
)

// ISessionService is an autogenerated mock type for the ISessionService type
type ISessionService struct {
	mock.Mock
}

type ISessionService_Expecter struct {
	mock *mock.Mock
}

func (_m *ISessionService) EXPECT() *ISessionService_Expecter {
	return &ISessionService_Expecter{mock: &_m.Mock}
}

// DecryptAndParsePayload provides a mock function with given fields: session, parseTo
func (_m *ISessionService) DecryptAndParsePayload(session dto.SessionDTO, parseTo interface{}) error {
	ret := _m.Called(session, parseTo)

	if len(ret) == 0 {
		panic("no return value specified for DecryptAndParsePayload")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(dto.SessionDTO, interface{}) error); ok {
		r0 = rf(session, parseTo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ISessionService_DecryptAndParsePayload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecryptAndParsePayload'
type ISessionService_DecryptAndParsePayload_Call struct {
	*mock.Call
}

// DecryptAndParsePayload is a helper method to define mock.On call
//   - session dto.SessionDTO
//   - parseTo interface{}
func (_e *ISessionService_Expecter) DecryptAndParsePayload(session interface{}, parseTo interface{}) *ISessionService_DecryptAndParsePayload_Call {
	return &ISessionService_DecryptAndParsePayload_Call{Call: _e.mock.On("DecryptAndParsePayload", session, parseTo)}
}

func (_c *ISessionService_DecryptAndParsePayload_Call) Run(run func(session dto.SessionDTO, parseTo interface{})) *ISessionService_DecryptAndParsePayload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(dto.SessionDTO), args[1].(interface{}))
	})
	return _c
}

func (_c *ISessionService_DecryptAndParsePayload_Call) Return(_a0 error) *ISessionService_DecryptAndParsePayload_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ISessionService_DecryptAndParsePayload_Call) RunAndReturn(run func(dto.SessionDTO, interface{}) error) *ISessionService_DecryptAndParsePayload_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteSession provides a mock function with given fields: ctx, prefix, session
func (_m *ISessionService) DeleteSession(ctx context.Context, prefix string, session string) error {
	ret := _m.Called(ctx, prefix, session)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, prefix, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ISessionService_DeleteSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteSession'
type ISessionService_DeleteSession_Call struct {
	*mock.Call
}

// DeleteSession is a helper method to define mock.On call
//   - ctx context.Context
//   - prefix string
//   - session string
func (_e *ISessionService_Expecter) DeleteSession(ctx interface{}, prefix interface{}, session interface{}) *ISessionService_DeleteSession_Call {
	return &ISessionService_DeleteSession_Call{Call: _e.mock.On("DeleteSession", ctx, prefix, session)}
}

func (_c *ISessionService_DeleteSession_Call) Run(run func(ctx context.Context, prefix string, session string)) *ISessionService_DeleteSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *ISessionService_DeleteSession_Call) Return(_a0 error) *ISessionService_DeleteSession_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ISessionService_DeleteSession_Call) RunAndReturn(run func(context.Context, string, string) error) *ISessionService_DeleteSession_Call {
	_c.Call.Return(run)
	return _c
}

// GetSession provides a mock function with given fields: ctx, prefix, session
func (_m *ISessionService) GetSession(ctx context.Context, prefix string, session string) (dto.SessionDTO, error) {
	ret := _m.Called(ctx, prefix, session)

	if len(ret) == 0 {
		panic("no return value specified for GetSession")
	}

	var r0 dto.SessionDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (dto.SessionDTO, error)); ok {
		return rf(ctx, prefix, session)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) dto.SessionDTO); ok {
		r0 = rf(ctx, prefix, session)
	} else {
		r0 = ret.Get(0).(dto.SessionDTO)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, prefix, session)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISessionService_GetSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSession'
type ISessionService_GetSession_Call struct {
	*mock.Call
}

// GetSession is a helper method to define mock.On call
//   - ctx context.Context
//   - prefix string
//   - session string
func (_e *ISessionService_Expecter) GetSession(ctx interface{}, prefix interface{}, session interface{}) *ISessionService_GetSession_Call {
	return &ISessionService_GetSession_Call{Call: _e.mock.On("GetSession", ctx, prefix, session)}
}

func (_c *ISessionService_GetSession_Call) Run(run func(ctx context.Context, prefix string, session string)) *ISessionService_GetSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *ISessionService_GetSession_Call) Return(_a0 dto.SessionDTO, _a1 error) *ISessionService_GetSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ISessionService_GetSession_Call) RunAndReturn(run func(context.Context, string, string) (dto.SessionDTO, error)) *ISessionService_GetSession_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserByAuthSession provides a mock function with given fields: ctx, session
func (_m *ISessionService) GetUserByAuthSession(ctx context.Context, session string) (dto.UserDTO, error) {
	ret := _m.Called(ctx, session)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByAuthSession")
	}

	var r0 dto.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (dto.UserDTO, error)); ok {
		return rf(ctx, session)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) dto.UserDTO); ok {
		r0 = rf(ctx, session)
	} else {
		r0 = ret.Get(0).(dto.UserDTO)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, session)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISessionService_GetUserByAuthSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByAuthSession'
type ISessionService_GetUserByAuthSession_Call struct {
	*mock.Call
}

// GetUserByAuthSession is a helper method to define mock.On call
//   - ctx context.Context
//   - session string
func (_e *ISessionService_Expecter) GetUserByAuthSession(ctx interface{}, session interface{}) *ISessionService_GetUserByAuthSession_Call {
	return &ISessionService_GetUserByAuthSession_Call{Call: _e.mock.On("GetUserByAuthSession", ctx, session)}
}

func (_c *ISessionService_GetUserByAuthSession_Call) Run(run func(ctx context.Context, session string)) *ISessionService_GetUserByAuthSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ISessionService_GetUserByAuthSession_Call) Return(_a0 dto.UserDTO, _a1 error) *ISessionService_GetUserByAuthSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ISessionService_GetUserByAuthSession_Call) RunAndReturn(run func(context.Context, string) (dto.UserDTO, error)) *ISessionService_GetUserByAuthSession_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserByEmailSession provides a mock function with given fields: ctx, session
func (_m *ISessionService) GetUserByEmailSession(ctx context.Context, session string) (dto.UserDTO, error) {
	ret := _m.Called(ctx, session)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmailSession")
	}

	var r0 dto.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (dto.UserDTO, error)); ok {
		return rf(ctx, session)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) dto.UserDTO); ok {
		r0 = rf(ctx, session)
	} else {
		r0 = ret.Get(0).(dto.UserDTO)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, session)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISessionService_GetUserByEmailSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByEmailSession'
type ISessionService_GetUserByEmailSession_Call struct {
	*mock.Call
}

// GetUserByEmailSession is a helper method to define mock.On call
//   - ctx context.Context
//   - session string
func (_e *ISessionService_Expecter) GetUserByEmailSession(ctx interface{}, session interface{}) *ISessionService_GetUserByEmailSession_Call {
	return &ISessionService_GetUserByEmailSession_Call{Call: _e.mock.On("GetUserByEmailSession", ctx, session)}
}

func (_c *ISessionService_GetUserByEmailSession_Call) Run(run func(ctx context.Context, session string)) *ISessionService_GetUserByEmailSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ISessionService_GetUserByEmailSession_Call) Return(_a0 dto.UserDTO, _a1 error) *ISessionService_GetUserByEmailSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ISessionService_GetUserByEmailSession_Call) RunAndReturn(run func(context.Context, string) (dto.UserDTO, error)) *ISessionService_GetUserByEmailSession_Call {
	_c.Call.Return(run)
	return _c
}

// IsExist provides a mock function with given fields: ctx, prefix, session
func (_m *ISessionService) IsExist(ctx context.Context, prefix string, session string) bool {
	ret := _m.Called(ctx, prefix, session)

	if len(ret) == 0 {
		panic("no return value specified for IsExist")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, prefix, session)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ISessionService_IsExist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsExist'
type ISessionService_IsExist_Call struct {
	*mock.Call
}

// IsExist is a helper method to define mock.On call
//   - ctx context.Context
//   - prefix string
//   - session string
func (_e *ISessionService_Expecter) IsExist(ctx interface{}, prefix interface{}, session interface{}) *ISessionService_IsExist_Call {
	return &ISessionService_IsExist_Call{Call: _e.mock.On("IsExist", ctx, prefix, session)}
}

func (_c *ISessionService_IsExist_Call) Run(run func(ctx context.Context, prefix string, session string)) *ISessionService_IsExist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *ISessionService_IsExist_Call) Return(_a0 bool) *ISessionService_IsExist_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ISessionService_IsExist_Call) RunAndReturn(run func(context.Context, string, string) bool) *ISessionService_IsExist_Call {
	_c.Call.Return(run)
	return _c
}

// SetSession provides a mock function with given fields: ctx, session
func (_m *ISessionService) SetSession(ctx context.Context, session dto.SessionDTO) (string, error) {
	ret := _m.Called(ctx, session)

	if len(ret) == 0 {
		panic("no return value specified for SetSession")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.SessionDTO) (string, error)); ok {
		return rf(ctx, session)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.SessionDTO) string); ok {
		r0 = rf(ctx, session)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.SessionDTO) error); ok {
		r1 = rf(ctx, session)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISessionService_SetSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetSession'
type ISessionService_SetSession_Call struct {
	*mock.Call
}

// SetSession is a helper method to define mock.On call
//   - ctx context.Context
//   - session dto.SessionDTO
func (_e *ISessionService_Expecter) SetSession(ctx interface{}, session interface{}) *ISessionService_SetSession_Call {
	return &ISessionService_SetSession_Call{Call: _e.mock.On("SetSession", ctx, session)}
}

func (_c *ISessionService_SetSession_Call) Run(run func(ctx context.Context, session dto.SessionDTO)) *ISessionService_SetSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.SessionDTO))
	})
	return _c
}

func (_c *ISessionService_SetSession_Call) Return(_a0 string, _a1 error) *ISessionService_SetSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ISessionService_SetSession_Call) RunAndReturn(run func(context.Context, dto.SessionDTO) (string, error)) *ISessionService_SetSession_Call {
	_c.Call.Return(run)
	return _c
}

// NewISessionService creates a new instance of ISessionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISessionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISessionService {
	mock := &ISessionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
