// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	models "homework/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// MockOrderService is an autogenerated mock type for the OrderService type
type MockOrderService struct {
	mock.Mock
}

type MockOrderService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOrderService) EXPECT() *MockOrderService_Expecter {
	return &MockOrderService_Expecter{mock: &_m.Mock}
}

// Accept provides a mock function with given fields: order, pkgTypeStr
func (_m *MockOrderService) Accept(order *models.Order, pkgTypeStr string) error {
	ret := _m.Called(order, pkgTypeStr)

	if len(ret) == 0 {
		panic("no return value specified for Accept")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Order, string) error); ok {
		r0 = rf(order, pkgTypeStr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrderService_Accept_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Accept'
type MockOrderService_Accept_Call struct {
	*mock.Call
}

// Accept is a helper method to define mock.On call
//   - order *models.Order
//   - pkgTypeStr string
func (_e *MockOrderService_Expecter) Accept(order interface{}, pkgTypeStr interface{}) *MockOrderService_Accept_Call {
	return &MockOrderService_Accept_Call{Call: _e.mock.On("Accept", order, pkgTypeStr)}
}

func (_c *MockOrderService_Accept_Call) Run(run func(order *models.Order, pkgTypeStr string)) *MockOrderService_Accept_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Order), args[1].(string))
	})
	return _c
}

func (_c *MockOrderService_Accept_Call) Return(_a0 error) *MockOrderService_Accept_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrderService_Accept_Call) RunAndReturn(run func(*models.Order, string) error) *MockOrderService_Accept_Call {
	_c.Call.Return(run)
	return _c
}

// Exists provides a mock function with given fields: userId
func (_m *MockOrderService) Exists(userId string) (models.Order, bool) {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 models.Order
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (models.Order, bool)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(string) models.Order); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Get(0).(models.Order)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockOrderService_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type MockOrderService_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - userId string
func (_e *MockOrderService_Expecter) Exists(userId interface{}) *MockOrderService_Exists_Call {
	return &MockOrderService_Exists_Call{Call: _e.mock.On("Exists", userId)}
}

func (_c *MockOrderService_Exists_Call) Run(run func(userId string)) *MockOrderService_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockOrderService_Exists_Call) Return(_a0 models.Order, _a1 bool) *MockOrderService_Exists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrderService_Exists_Call) RunAndReturn(run func(string) (models.Order, bool)) *MockOrderService_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// Issue provides a mock function with given fields: ordersToIssue
func (_m *MockOrderService) Issue(ordersToIssue *[]models.Order) error {
	ret := _m.Called(ordersToIssue)

	if len(ret) == 0 {
		panic("no return value specified for Issue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*[]models.Order) error); ok {
		r0 = rf(ordersToIssue)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrderService_Issue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Issue'
type MockOrderService_Issue_Call struct {
	*mock.Call
}

// Issue is a helper method to define mock.On call
//   - ordersToIssue *[]models.Order
func (_e *MockOrderService_Expecter) Issue(ordersToIssue interface{}) *MockOrderService_Issue_Call {
	return &MockOrderService_Issue_Call{Call: _e.mock.On("Issue", ordersToIssue)}
}

func (_c *MockOrderService_Issue_Call) Run(run func(ordersToIssue *[]models.Order)) *MockOrderService_Issue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*[]models.Order))
	})
	return _c
}

func (_c *MockOrderService_Issue_Call) Return(_a0 error) *MockOrderService_Issue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrderService_Issue_Call) RunAndReturn(run func(*[]models.Order) error) *MockOrderService_Issue_Call {
	_c.Call.Return(run)
	return _c
}

// ListOrders provides a mock function with given fields: userId, offset, limit
func (_m *MockOrderService) ListOrders(userId string, offset int, limit int) ([]models.Order, error) {
	ret := _m.Called(userId, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for ListOrders")
	}

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, int) ([]models.Order, error)); ok {
		return rf(userId, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(string, int, int) []models.Order); ok {
		r0 = rf(userId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(userId, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrderService_ListOrders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListOrders'
type MockOrderService_ListOrders_Call struct {
	*mock.Call
}

// ListOrders is a helper method to define mock.On call
//   - userId string
//   - offset int
//   - limit int
func (_e *MockOrderService_Expecter) ListOrders(userId interface{}, offset interface{}, limit interface{}) *MockOrderService_ListOrders_Call {
	return &MockOrderService_ListOrders_Call{Call: _e.mock.On("ListOrders", userId, offset, limit)}
}

func (_c *MockOrderService_ListOrders_Call) Run(run func(userId string, offset int, limit int)) *MockOrderService_ListOrders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *MockOrderService_ListOrders_Call) Return(_a0 []models.Order, _a1 error) *MockOrderService_ListOrders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrderService_ListOrders_Call) RunAndReturn(run func(string, int, int) ([]models.Order, error)) *MockOrderService_ListOrders_Call {
	_c.Call.Return(run)
	return _c
}

// ListReturns provides a mock function with given fields: offset, limit
func (_m *MockOrderService) ListReturns(offset int, limit int) ([]models.Order, error) {
	ret := _m.Called(offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for ListReturns")
	}

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]models.Order, error)); ok {
		return rf(offset, limit)
	}
	if rf, ok := ret.Get(0).(func(int, int) []models.Order); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrderService_ListReturns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReturns'
type MockOrderService_ListReturns_Call struct {
	*mock.Call
}

// ListReturns is a helper method to define mock.On call
//   - offset int
//   - limit int
func (_e *MockOrderService_Expecter) ListReturns(offset interface{}, limit interface{}) *MockOrderService_ListReturns_Call {
	return &MockOrderService_ListReturns_Call{Call: _e.mock.On("ListReturns", offset, limit)}
}

func (_c *MockOrderService_ListReturns_Call) Run(run func(offset int, limit int)) *MockOrderService_ListReturns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(int))
	})
	return _c
}

func (_c *MockOrderService_ListReturns_Call) Return(_a0 []models.Order, _a1 error) *MockOrderService_ListReturns_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrderService_ListReturns_Call) RunAndReturn(run func(int, int) ([]models.Order, error)) *MockOrderService_ListReturns_Call {
	_c.Call.Return(run)
	return _c
}

// PrintList provides a mock function with given fields: orders
func (_m *MockOrderService) PrintList(orders []models.Order) {
	_m.Called(orders)
}

// MockOrderService_PrintList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrintList'
type MockOrderService_PrintList_Call struct {
	*mock.Call
}

// PrintList is a helper method to define mock.On call
//   - orders []models.Order
func (_e *MockOrderService_Expecter) PrintList(orders interface{}) *MockOrderService_PrintList_Call {
	return &MockOrderService_PrintList_Call{Call: _e.mock.On("PrintList", orders)}
}

func (_c *MockOrderService_PrintList_Call) Run(run func(orders []models.Order)) *MockOrderService_PrintList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]models.Order))
	})
	return _c
}

func (_c *MockOrderService_PrintList_Call) Return() *MockOrderService_PrintList_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockOrderService_PrintList_Call) RunAndReturn(run func([]models.Order)) *MockOrderService_PrintList_Call {
	_c.Call.Return(run)
	return _c
}

// Return provides a mock function with given fields: orders
func (_m *MockOrderService) Return(orders *models.Order) error {
	ret := _m.Called(orders)

	if len(ret) == 0 {
		panic("no return value specified for Return")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Order) error); ok {
		r0 = rf(orders)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrderService_Return_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Return'
type MockOrderService_Return_Call struct {
	*mock.Call
}

// Return is a helper method to define mock.On call
//   - orders *models.Order
func (_e *MockOrderService_Expecter) Return(orders interface{}) *MockOrderService_Return_Call {
	return &MockOrderService_Return_Call{Call: _e.mock.On("Return", orders)}
}

func (_c *MockOrderService_Return_Call) Run(run func(orders *models.Order)) *MockOrderService_Return_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Order))
	})
	return _c
}

func (_c *MockOrderService_Return_Call) Return(_a0 error) *MockOrderService_Return_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrderService_Return_Call) RunAndReturn(run func(*models.Order) error) *MockOrderService_Return_Call {
	_c.Call.Return(run)
	return _c
}

// ReturnToCourier provides a mock function with given fields: id
func (_m *MockOrderService) ReturnToCourier(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for ReturnToCourier")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrderService_ReturnToCourier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReturnToCourier'
type MockOrderService_ReturnToCourier_Call struct {
	*mock.Call
}

// ReturnToCourier is a helper method to define mock.On call
//   - id string
func (_e *MockOrderService_Expecter) ReturnToCourier(id interface{}) *MockOrderService_ReturnToCourier_Call {
	return &MockOrderService_ReturnToCourier_Call{Call: _e.mock.On("ReturnToCourier", id)}
}

func (_c *MockOrderService_ReturnToCourier_Call) Run(run func(id string)) *MockOrderService_ReturnToCourier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockOrderService_ReturnToCourier_Call) Return(_a0 error) *MockOrderService_ReturnToCourier_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrderService_ReturnToCourier_Call) RunAndReturn(run func(string) error) *MockOrderService_ReturnToCourier_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOrderService creates a new instance of MockOrderService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrderService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrderService {
	mock := &MockOrderService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
