// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	domain "libs/src/internal/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// IUserRepository is an autogenerated mock type for the IUserRepository type
type IUserRepository struct {
	mock.Mock
}

type IUserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *IUserRepository) EXPECT() *IUserRepository_Expecter {
	return &IUserRepository_Expecter{mock: &_m.Mock}
}

// Count provides a mock function with given fields: filter, args
func (_m *IUserRepository) Count(filter string, args ...interface{}) (int64, error) {
	var _ca []interface{}
	_ca = append(_ca, filter)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Count")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (int64, error)); ok {
		return rf(filter, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) int64); ok {
		r0 = rf(filter, args...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(filter, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IUserRepository_Count_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Count'
type IUserRepository_Count_Call struct {
	*mock.Call
}

// Count is a helper method to define mock.On call
//   - filter string
//   - args ...interface{}
func (_e *IUserRepository_Expecter) Count(filter interface{}, args ...interface{}) *IUserRepository_Count_Call {
	return &IUserRepository_Count_Call{Call: _e.mock.On("Count",
		append([]interface{}{filter}, args...)...)}
}

func (_c *IUserRepository_Count_Call) Run(run func(filter string, args ...interface{})) *IUserRepository_Count_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *IUserRepository_Count_Call) Return(_a0 int64, _a1 error) *IUserRepository_Count_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IUserRepository_Count_Call) RunAndReturn(run func(string, ...interface{}) (int64, error)) *IUserRepository_Count_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: obj
func (_m *IUserRepository) Create(obj *domain.User) error {
	ret := _m.Called(obj)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IUserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type IUserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - obj *domain.User
func (_e *IUserRepository_Expecter) Create(obj interface{}) *IUserRepository_Create_Call {
	return &IUserRepository_Create_Call{Call: _e.mock.On("Create", obj)}
}

func (_c *IUserRepository_Create_Call) Run(run func(obj *domain.User)) *IUserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.User))
	})
	return _c
}

func (_c *IUserRepository_Create_Call) Return(_a0 error) *IUserRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IUserRepository_Create_Call) RunAndReturn(run func(*domain.User) error) *IUserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteById provides a mock function with given fields: id
func (_m *IUserRepository) DeleteById(id int64) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteById")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IUserRepository_DeleteById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteById'
type IUserRepository_DeleteById_Call struct {
	*mock.Call
}

// DeleteById is a helper method to define mock.On call
//   - id int64
func (_e *IUserRepository_Expecter) DeleteById(id interface{}) *IUserRepository_DeleteById_Call {
	return &IUserRepository_DeleteById_Call{Call: _e.mock.On("DeleteById", id)}
}

func (_c *IUserRepository_DeleteById_Call) Run(run func(id int64)) *IUserRepository_DeleteById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *IUserRepository_DeleteById_Call) Return(_a0 error) *IUserRepository_DeleteById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IUserRepository_DeleteById_Call) RunAndReturn(run func(int64) error) *IUserRepository_DeleteById_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteQuery provides a mock function with given fields: query, args
func (_m *IUserRepository) ExecuteQuery(query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteQuery")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) error); ok {
		r0 = rf(query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IUserRepository_ExecuteQuery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteQuery'
type IUserRepository_ExecuteQuery_Call struct {
	*mock.Call
}

// ExecuteQuery is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *IUserRepository_Expecter) ExecuteQuery(query interface{}, args ...interface{}) *IUserRepository_ExecuteQuery_Call {
	return &IUserRepository_ExecuteQuery_Call{Call: _e.mock.On("ExecuteQuery",
		append([]interface{}{query}, args...)...)}
}

func (_c *IUserRepository_ExecuteQuery_Call) Run(run func(query string, args ...interface{})) *IUserRepository_ExecuteQuery_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *IUserRepository_ExecuteQuery_Call) Return(_a0 error) *IUserRepository_ExecuteQuery_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IUserRepository_ExecuteQuery_Call) RunAndReturn(run func(string, ...interface{}) error) *IUserRepository_ExecuteQuery_Call {
	_c.Call.Return(run)
	return _c
}

// Filter provides a mock function with given fields: query, args
func (_m *IUserRepository) Filter(query string, args ...interface{}) ([]domain.User, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Filter")
	}

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) ([]domain.User, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) []domain.User); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IUserRepository_Filter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Filter'
type IUserRepository_Filter_Call struct {
	*mock.Call
}

// Filter is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *IUserRepository_Expecter) Filter(query interface{}, args ...interface{}) *IUserRepository_Filter_Call {
	return &IUserRepository_Filter_Call{Call: _e.mock.On("Filter",
		append([]interface{}{query}, args...)...)}
}

func (_c *IUserRepository_Filter_Call) Run(run func(query string, args ...interface{})) *IUserRepository_Filter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *IUserRepository_Filter_Call) Return(_a0 []domain.User, _a1 error) *IUserRepository_Filter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IUserRepository_Filter_Call) RunAndReturn(run func(string, ...interface{}) ([]domain.User, error)) *IUserRepository_Filter_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with no fields
func (_m *IUserRepository) GetAll() ([]domain.User, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IUserRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type IUserRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *IUserRepository_Expecter) GetAll() *IUserRepository_GetAll_Call {
	return &IUserRepository_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *IUserRepository_GetAll_Call) Run(run func()) *IUserRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IUserRepository_GetAll_Call) Return(_a0 []domain.User, _a1 error) *IUserRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IUserRepository_GetAll_Call) RunAndReturn(run func() ([]domain.User, error)) *IUserRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetById provides a mock function with given fields: id
func (_m *IUserRepository) GetById(id int64) (domain.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (domain.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) domain.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IUserRepository_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type IUserRepository_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id int64
func (_e *IUserRepository_Expecter) GetById(id interface{}) *IUserRepository_GetById_Call {
	return &IUserRepository_GetById_Call{Call: _e.mock.On("GetById", id)}
}

func (_c *IUserRepository_GetById_Call) Run(run func(id int64)) *IUserRepository_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *IUserRepository_GetById_Call) Return(_a0 domain.User, _a1 error) *IUserRepository_GetById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IUserRepository_GetById_Call) RunAndReturn(run func(int64) (domain.User, error)) *IUserRepository_GetById_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUsername provides a mock function with given fields: username
func (_m *IUserRepository) GetByUsername(username string) (domain.User, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetByUsername")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.User, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) domain.User); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IUserRepository_GetByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUsername'
type IUserRepository_GetByUsername_Call struct {
	*mock.Call
}

// GetByUsername is a helper method to define mock.On call
//   - username string
func (_e *IUserRepository_Expecter) GetByUsername(username interface{}) *IUserRepository_GetByUsername_Call {
	return &IUserRepository_GetByUsername_Call{Call: _e.mock.On("GetByUsername", username)}
}

func (_c *IUserRepository_GetByUsername_Call) Run(run func(username string)) *IUserRepository_GetByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IUserRepository_GetByUsername_Call) Return(_a0 domain.User, _a1 error) *IUserRepository_GetByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IUserRepository_GetByUsername_Call) RunAndReturn(run func(string) (domain.User, error)) *IUserRepository_GetByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateById provides a mock function with given fields: id, updateFields
func (_m *IUserRepository) UpdateById(id int64, updateFields map[string]interface{}) error {
	ret := _m.Called(id, updateFields)

	if len(ret) == 0 {
		panic("no return value specified for UpdateById")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, map[string]interface{}) error); ok {
		r0 = rf(id, updateFields)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IUserRepository_UpdateById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateById'
type IUserRepository_UpdateById_Call struct {
	*mock.Call
}

// UpdateById is a helper method to define mock.On call
//   - id int64
//   - updateFields map[string]interface{}
func (_e *IUserRepository_Expecter) UpdateById(id interface{}, updateFields interface{}) *IUserRepository_UpdateById_Call {
	return &IUserRepository_UpdateById_Call{Call: _e.mock.On("UpdateById", id, updateFields)}
}

func (_c *IUserRepository_UpdateById_Call) Run(run func(id int64, updateFields map[string]interface{})) *IUserRepository_UpdateById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *IUserRepository_UpdateById_Call) Return(_a0 error) *IUserRepository_UpdateById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IUserRepository_UpdateById_Call) RunAndReturn(run func(int64, map[string]interface{}) error) *IUserRepository_UpdateById_Call {
	_c.Call.Return(run)
	return _c
}

// NewIUserRepository creates a new instance of IUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUserRepository {
	mock := &IUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
