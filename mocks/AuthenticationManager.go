// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	DatabaseAbstraction "EntitlementServer/DatabaseAbstraction"

	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// AuthenticationManager is an autogenerated mock type for the AuthenticationManager type
type AuthenticationManager struct {
	mock.Mock
}

type AuthenticationManager_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthenticationManager) EXPECT() *AuthenticationManager_Expecter {
	return &AuthenticationManager_Expecter{mock: &_m.Mock}
}

// AuthenticateUser provides a mock function with given fields: username, password
func (_m *AuthenticationManager) AuthenticateUser(username string, password string) (bool, error) {
	ret := _m.Called(username, password)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationManager_AuthenticateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AuthenticateUser'
type AuthenticationManager_AuthenticateUser_Call struct {
	*mock.Call
}

// AuthenticateUser is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *AuthenticationManager_Expecter) AuthenticateUser(username interface{}, password interface{}) *AuthenticationManager_AuthenticateUser_Call {
	return &AuthenticationManager_AuthenticateUser_Call{Call: _e.mock.On("AuthenticateUser", username, password)}
}

func (_c *AuthenticationManager_AuthenticateUser_Call) Run(run func(username string, password string)) *AuthenticationManager_AuthenticateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationManager_AuthenticateUser_Call) Return(_a0 bool, _a1 error) *AuthenticationManager_AuthenticateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// AuthenticationMiddleware provides a mock function with given fields: c
func (_m *AuthenticationManager) AuthenticationMiddleware(c *gin.Context) {
	_m.Called(c)
}

// AuthenticationManager_AuthenticationMiddleware_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AuthenticationMiddleware'
type AuthenticationManager_AuthenticationMiddleware_Call struct {
	*mock.Call
}

// AuthenticationMiddleware is a helper method to define mock.On call
//   - c *gin.Context
func (_e *AuthenticationManager_Expecter) AuthenticationMiddleware(c interface{}) *AuthenticationManager_AuthenticationMiddleware_Call {
	return &AuthenticationManager_AuthenticationMiddleware_Call{Call: _e.mock.On("AuthenticationMiddleware", c)}
}

func (_c *AuthenticationManager_AuthenticationMiddleware_Call) Run(run func(c *gin.Context)) *AuthenticationManager_AuthenticationMiddleware_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *AuthenticationManager_AuthenticationMiddleware_Call) Return() *AuthenticationManager_AuthenticationMiddleware_Call {
	_c.Call.Return()
	return _c
}

// ComparePasswords provides a mock function with given fields: hashedPassword, password
func (_m *AuthenticationManager) ComparePasswords(hashedPassword string, password string) (bool, error) {
	ret := _m.Called(hashedPassword, password)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(hashedPassword, password)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(hashedPassword, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationManager_ComparePasswords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ComparePasswords'
type AuthenticationManager_ComparePasswords_Call struct {
	*mock.Call
}

// ComparePasswords is a helper method to define mock.On call
//   - hashedPassword string
//   - password string
func (_e *AuthenticationManager_Expecter) ComparePasswords(hashedPassword interface{}, password interface{}) *AuthenticationManager_ComparePasswords_Call {
	return &AuthenticationManager_ComparePasswords_Call{Call: _e.mock.On("ComparePasswords", hashedPassword, password)}
}

func (_c *AuthenticationManager_ComparePasswords_Call) Run(run func(hashedPassword string, password string)) *AuthenticationManager_ComparePasswords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationManager_ComparePasswords_Call) Return(_a0 bool, _a1 error) *AuthenticationManager_ComparePasswords_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// CreateToken provides a mock function with given fields: userid
func (_m *AuthenticationManager) CreateToken(userid int) (string, error) {
	ret := _m.Called(userid)

	var r0 string
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(userid)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationManager_CreateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateToken'
type AuthenticationManager_CreateToken_Call struct {
	*mock.Call
}

// CreateToken is a helper method to define mock.On call
//   - userid int
func (_e *AuthenticationManager_Expecter) CreateToken(userid interface{}) *AuthenticationManager_CreateToken_Call {
	return &AuthenticationManager_CreateToken_Call{Call: _e.mock.On("CreateToken", userid)}
}

func (_c *AuthenticationManager_CreateToken_Call) Run(run func(userid int)) *AuthenticationManager_CreateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *AuthenticationManager_CreateToken_Call) Return(_a0 string, _a1 error) *AuthenticationManager_CreateToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// CreateUser provides a mock function with given fields: username, password
func (_m *AuthenticationManager) CreateUser(username string, password string) error {
	ret := _m.Called(username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticationManager_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type AuthenticationManager_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *AuthenticationManager_Expecter) CreateUser(username interface{}, password interface{}) *AuthenticationManager_CreateUser_Call {
	return &AuthenticationManager_CreateUser_Call{Call: _e.mock.On("CreateUser", username, password)}
}

func (_c *AuthenticationManager_CreateUser_Call) Run(run func(username string, password string)) *AuthenticationManager_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationManager_CreateUser_Call) Return(_a0 error) *AuthenticationManager_CreateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

// HashPassword provides a mock function with given fields: password
func (_m *AuthenticationManager) HashPassword(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationManager_HashPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HashPassword'
type AuthenticationManager_HashPassword_Call struct {
	*mock.Call
}

// HashPassword is a helper method to define mock.On call
//   - password string
func (_e *AuthenticationManager_Expecter) HashPassword(password interface{}) *AuthenticationManager_HashPassword_Call {
	return &AuthenticationManager_HashPassword_Call{Call: _e.mock.On("HashPassword", password)}
}

func (_c *AuthenticationManager_HashPassword_Call) Run(run func(password string)) *AuthenticationManager_HashPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AuthenticationManager_HashPassword_Call) Return(_a0 string, _a1 error) *AuthenticationManager_HashPassword_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// ValidateToken provides a mock function with given fields: token
func (_m *AuthenticationManager) ValidateToken(token string) (bool, DatabaseAbstraction.User, error) {
	ret := _m.Called(token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 DatabaseAbstraction.User
	if rf, ok := ret.Get(1).(func(string) DatabaseAbstraction.User); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Get(1).(DatabaseAbstraction.User)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(token)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// AuthenticationManager_ValidateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateToken'
type AuthenticationManager_ValidateToken_Call struct {
	*mock.Call
}

// ValidateToken is a helper method to define mock.On call
//   - token string
func (_e *AuthenticationManager_Expecter) ValidateToken(token interface{}) *AuthenticationManager_ValidateToken_Call {
	return &AuthenticationManager_ValidateToken_Call{Call: _e.mock.On("ValidateToken", token)}
}

func (_c *AuthenticationManager_ValidateToken_Call) Run(run func(token string)) *AuthenticationManager_ValidateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AuthenticationManager_ValidateToken_Call) Return(_a0 bool, _a1 DatabaseAbstraction.User, _a2 error) *AuthenticationManager_ValidateToken_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

type mockConstructorTestingTNewAuthenticationManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthenticationManager creates a new instance of AuthenticationManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthenticationManager(t mockConstructorTestingTNewAuthenticationManager) *AuthenticationManager {
	mock := &AuthenticationManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
