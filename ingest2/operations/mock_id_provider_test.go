// Code generated by mockery v1.0.0. DO NOT EDIT.

package operations

import mock "github.com/stretchr/testify/mock"
import xdr "gitlab.com/tokend/go/xdr"

// MockIDProvider is an autogenerated mock type for the IDProvider type
type MockIDProvider struct {
	mock.Mock
}

// MustAccountID provides a mock function with given fields: raw
func (_m *MockIDProvider) MustAccountID(raw xdr.AccountId) uint64 {
	ret := _m.Called(raw)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(xdr.AccountId) uint64); ok {
		r0 = rf(raw)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// MustBalanceID provides a mock function with given fields: raw
func (_m *MockIDProvider) MustBalanceID(raw xdr.BalanceId) uint64 {
	ret := _m.Called(raw)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(xdr.BalanceId) uint64); ok {
		r0 = rf(raw)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}
