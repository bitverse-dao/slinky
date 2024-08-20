// Code generated by mockery v2.44.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	pkgtypes "github.com/skip-mev/slinky/pkg/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// OracleKeeper is an autogenerated mock type for the OracleKeeper type
type OracleKeeper struct {
	mock.Mock
}

// HasCurrencyPair provides a mock function with given fields: ctx, cp
func (_m *OracleKeeper) HasCurrencyPair(ctx types.Context, cp pkgtypes.CurrencyPair) bool {
	ret := _m.Called(ctx, cp)

	if len(ret) == 0 {
		panic("no return value specified for HasCurrencyPair")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, pkgtypes.CurrencyPair) bool); ok {
		r0 = rf(ctx, cp)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewOracleKeeper creates a new instance of OracleKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOracleKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *OracleKeeper {
	mock := &OracleKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
