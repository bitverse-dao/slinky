// Code generated by mockery v2.44.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	rpc "github.com/gagliardetto/solana-go/rpc"

	solana "github.com/gagliardetto/solana-go"
)

// SolanaJSONRPCClient is an autogenerated mock type for the SolanaJSONRPCClient type
type SolanaJSONRPCClient struct {
	mock.Mock
}

// GetMultipleAccountsWithOpts provides a mock function with given fields: ctx, accounts, opts
func (_m *SolanaJSONRPCClient) GetMultipleAccountsWithOpts(ctx context.Context, accounts []solana.PublicKey, opts *rpc.GetMultipleAccountsOpts) (*rpc.GetMultipleAccountsResult, error) {
	ret := _m.Called(ctx, accounts, opts)

	if len(ret) == 0 {
		panic("no return value specified for GetMultipleAccountsWithOpts")
	}

	var r0 *rpc.GetMultipleAccountsResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []solana.PublicKey, *rpc.GetMultipleAccountsOpts) (*rpc.GetMultipleAccountsResult, error)); ok {
		return rf(ctx, accounts, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []solana.PublicKey, *rpc.GetMultipleAccountsOpts) *rpc.GetMultipleAccountsResult); ok {
		r0 = rf(ctx, accounts, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpc.GetMultipleAccountsResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []solana.PublicKey, *rpc.GetMultipleAccountsOpts) error); ok {
		r1 = rf(ctx, accounts, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSolanaJSONRPCClient creates a new instance of SolanaJSONRPCClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSolanaJSONRPCClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SolanaJSONRPCClient {
	mock := &SolanaJSONRPCClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
