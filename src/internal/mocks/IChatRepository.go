// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	domain "libs/src/internal/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// IChatRepository is an autogenerated mock type for the IChatRepository type
type IChatRepository struct {
	mock.Mock
}

type IChatRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *IChatRepository) EXPECT() *IChatRepository_Expecter {
	return &IChatRepository_Expecter{mock: &_m.Mock}
}

// Count provides a mock function with given fields: filter, args
func (_m *IChatRepository) Count(filter string, args ...interface{}) (int64, error) {
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

// IChatRepository_Count_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Count'
type IChatRepository_Count_Call struct {
	*mock.Call
}

// Count is a helper method to define mock.On call
//   - filter string
//   - args ...interface{}
func (_e *IChatRepository_Expecter) Count(filter interface{}, args ...interface{}) *IChatRepository_Count_Call {
	return &IChatRepository_Count_Call{Call: _e.mock.On("Count",
		append([]interface{}{filter}, args...)...)}
}

func (_c *IChatRepository_Count_Call) Run(run func(filter string, args ...interface{})) *IChatRepository_Count_Call {
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

func (_c *IChatRepository_Count_Call) Return(_a0 int64, _a1 error) *IChatRepository_Count_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IChatRepository_Count_Call) RunAndReturn(run func(string, ...interface{}) (int64, error)) *IChatRepository_Count_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: obj
func (_m *IChatRepository) Create(obj *domain.Chat) error {
	ret := _m.Called(obj)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Chat) error); ok {
		r0 = rf(obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IChatRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type IChatRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - obj *domain.Chat
func (_e *IChatRepository_Expecter) Create(obj interface{}) *IChatRepository_Create_Call {
	return &IChatRepository_Create_Call{Call: _e.mock.On("Create", obj)}
}

func (_c *IChatRepository_Create_Call) Run(run func(obj *domain.Chat)) *IChatRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Chat))
	})
	return _c
}

func (_c *IChatRepository_Create_Call) Return(_a0 error) *IChatRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IChatRepository_Create_Call) RunAndReturn(run func(*domain.Chat) error) *IChatRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteById provides a mock function with given fields: id
func (_m *IChatRepository) DeleteById(id int64) error {
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

// IChatRepository_DeleteById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteById'
type IChatRepository_DeleteById_Call struct {
	*mock.Call
}

// DeleteById is a helper method to define mock.On call
//   - id int64
func (_e *IChatRepository_Expecter) DeleteById(id interface{}) *IChatRepository_DeleteById_Call {
	return &IChatRepository_DeleteById_Call{Call: _e.mock.On("DeleteById", id)}
}

func (_c *IChatRepository_DeleteById_Call) Run(run func(id int64)) *IChatRepository_DeleteById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *IChatRepository_DeleteById_Call) Return(_a0 error) *IChatRepository_DeleteById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IChatRepository_DeleteById_Call) RunAndReturn(run func(int64) error) *IChatRepository_DeleteById_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteQuery provides a mock function with given fields: query, args
func (_m *IChatRepository) ExecuteQuery(query string, args ...interface{}) error {
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

// IChatRepository_ExecuteQuery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteQuery'
type IChatRepository_ExecuteQuery_Call struct {
	*mock.Call
}

// ExecuteQuery is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *IChatRepository_Expecter) ExecuteQuery(query interface{}, args ...interface{}) *IChatRepository_ExecuteQuery_Call {
	return &IChatRepository_ExecuteQuery_Call{Call: _e.mock.On("ExecuteQuery",
		append([]interface{}{query}, args...)...)}
}

func (_c *IChatRepository_ExecuteQuery_Call) Run(run func(query string, args ...interface{})) *IChatRepository_ExecuteQuery_Call {
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

func (_c *IChatRepository_ExecuteQuery_Call) Return(_a0 error) *IChatRepository_ExecuteQuery_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IChatRepository_ExecuteQuery_Call) RunAndReturn(run func(string, ...interface{}) error) *IChatRepository_ExecuteQuery_Call {
	_c.Call.Return(run)
	return _c
}

// Filter provides a mock function with given fields: query, args
func (_m *IChatRepository) Filter(query string, args ...interface{}) ([]domain.Chat, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Filter")
	}

	var r0 []domain.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) ([]domain.Chat, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) []domain.Chat); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IChatRepository_Filter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Filter'
type IChatRepository_Filter_Call struct {
	*mock.Call
}

// Filter is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *IChatRepository_Expecter) Filter(query interface{}, args ...interface{}) *IChatRepository_Filter_Call {
	return &IChatRepository_Filter_Call{Call: _e.mock.On("Filter",
		append([]interface{}{query}, args...)...)}
}

func (_c *IChatRepository_Filter_Call) Run(run func(query string, args ...interface{})) *IChatRepository_Filter_Call {
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

func (_c *IChatRepository_Filter_Call) Return(_a0 []domain.Chat, _a1 error) *IChatRepository_Filter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IChatRepository_Filter_Call) RunAndReturn(run func(string, ...interface{}) ([]domain.Chat, error)) *IChatRepository_Filter_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with no fields
func (_m *IChatRepository) GetAll() ([]domain.Chat, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.Chat, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.Chat); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IChatRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type IChatRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *IChatRepository_Expecter) GetAll() *IChatRepository_GetAll_Call {
	return &IChatRepository_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *IChatRepository_GetAll_Call) Run(run func()) *IChatRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IChatRepository_GetAll_Call) Return(_a0 []domain.Chat, _a1 error) *IChatRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IChatRepository_GetAll_Call) RunAndReturn(run func() ([]domain.Chat, error)) *IChatRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetById provides a mock function with given fields: id
func (_m *IChatRepository) GetById(id int64) (domain.Chat, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 domain.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (domain.Chat, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) domain.Chat); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Chat)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IChatRepository_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type IChatRepository_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id int64
func (_e *IChatRepository_Expecter) GetById(id interface{}) *IChatRepository_GetById_Call {
	return &IChatRepository_GetById_Call{Call: _e.mock.On("GetById", id)}
}

func (_c *IChatRepository_GetById_Call) Run(run func(id int64)) *IChatRepository_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *IChatRepository_GetById_Call) Return(_a0 domain.Chat, _a1 error) *IChatRepository_GetById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IChatRepository_GetById_Call) RunAndReturn(run func(int64) (domain.Chat, error)) *IChatRepository_GetById_Call {
	_c.Call.Return(run)
	return _c
}

// GetListForUser provides a mock function with given fields: userId, limit, offset
func (_m *IChatRepository) GetListForUser(userId int64, limit int, offset int) ([]domain.Chat, error) {
	ret := _m.Called(userId, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetListForUser")
	}

	var r0 []domain.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int, int) ([]domain.Chat, error)); ok {
		return rf(userId, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int64, int, int) []domain.Chat); ok {
		r0 = rf(userId, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int, int) error); ok {
		r1 = rf(userId, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IChatRepository_GetListForUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetListForUser'
type IChatRepository_GetListForUser_Call struct {
	*mock.Call
}

// GetListForUser is a helper method to define mock.On call
//   - userId int64
//   - limit int
//   - offset int
func (_e *IChatRepository_Expecter) GetListForUser(userId interface{}, limit interface{}, offset interface{}) *IChatRepository_GetListForUser_Call {
	return &IChatRepository_GetListForUser_Call{Call: _e.mock.On("GetListForUser", userId, limit, offset)}
}

func (_c *IChatRepository_GetListForUser_Call) Run(run func(userId int64, limit int, offset int)) *IChatRepository_GetListForUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *IChatRepository_GetListForUser_Call) Return(_a0 []domain.Chat, _a1 error) *IChatRepository_GetListForUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IChatRepository_GetListForUser_Call) RunAndReturn(run func(int64, int, int) ([]domain.Chat, error)) *IChatRepository_GetListForUser_Call {
	_c.Call.Return(run)
	return _c
}

// SearchForUser provides a mock function with given fields: userId, name, limit, offset
func (_m *IChatRepository) SearchForUser(userId int64, name string, limit int, offset int) ([]domain.Chat, error) {
	ret := _m.Called(userId, name, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for SearchForUser")
	}

	var r0 []domain.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, string, int, int) ([]domain.Chat, error)); ok {
		return rf(userId, name, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int64, string, int, int) []domain.Chat); ok {
		r0 = rf(userId, name, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, string, int, int) error); ok {
		r1 = rf(userId, name, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IChatRepository_SearchForUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchForUser'
type IChatRepository_SearchForUser_Call struct {
	*mock.Call
}

// SearchForUser is a helper method to define mock.On call
//   - userId int64
//   - name string
//   - limit int
//   - offset int
func (_e *IChatRepository_Expecter) SearchForUser(userId interface{}, name interface{}, limit interface{}, offset interface{}) *IChatRepository_SearchForUser_Call {
	return &IChatRepository_SearchForUser_Call{Call: _e.mock.On("SearchForUser", userId, name, limit, offset)}
}

func (_c *IChatRepository_SearchForUser_Call) Run(run func(userId int64, name string, limit int, offset int)) *IChatRepository_SearchForUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(string), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *IChatRepository_SearchForUser_Call) Return(_a0 []domain.Chat, _a1 error) *IChatRepository_SearchForUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IChatRepository_SearchForUser_Call) RunAndReturn(run func(int64, string, int, int) ([]domain.Chat, error)) *IChatRepository_SearchForUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateById provides a mock function with given fields: id, updateFields
func (_m *IChatRepository) UpdateById(id int64, updateFields map[string]interface{}) error {
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

// IChatRepository_UpdateById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateById'
type IChatRepository_UpdateById_Call struct {
	*mock.Call
}

// UpdateById is a helper method to define mock.On call
//   - id int64
//   - updateFields map[string]interface{}
func (_e *IChatRepository_Expecter) UpdateById(id interface{}, updateFields interface{}) *IChatRepository_UpdateById_Call {
	return &IChatRepository_UpdateById_Call{Call: _e.mock.On("UpdateById", id, updateFields)}
}

func (_c *IChatRepository_UpdateById_Call) Run(run func(id int64, updateFields map[string]interface{})) *IChatRepository_UpdateById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *IChatRepository_UpdateById_Call) Return(_a0 error) *IChatRepository_UpdateById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IChatRepository_UpdateById_Call) RunAndReturn(run func(int64, map[string]interface{}) error) *IChatRepository_UpdateById_Call {
	_c.Call.Return(run)
	return _c
}

// NewIChatRepository creates a new instance of IChatRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIChatRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IChatRepository {
	mock := &IChatRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
