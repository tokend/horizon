// Code generated by mockery v1.0.0. DO NOT EDIT.

package exchange

import core "gitlab.com/tokend/horizon/db2/core"
import mock "github.com/stretchr/testify/mock"

// mockAssetProvider is an autogenerated mock type for the assetProvider type
type mockAssetProvider struct {
	mock.Mock
}

// GetAssetPairsForCodes provides a mock function with given fields: baseAssets, quoteAssets
func (_m *mockAssetProvider) GetAssetPairsForCodes(baseAssets []string, quoteAssets []string) ([]core.AssetPair, error) {
	ret := _m.Called(baseAssets, quoteAssets)

	var r0 []core.AssetPair
	if rf, ok := ret.Get(0).(func([]string, []string) []core.AssetPair); ok {
		r0 = rf(baseAssets, quoteAssets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.AssetPair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, []string) error); ok {
		r1 = rf(baseAssets, quoteAssets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAssetsForPolicy provides a mock function with given fields: policy
func (_m *mockAssetProvider) GetAssetsForPolicy(policy uint32) ([]core.Asset, error) {
	ret := _m.Called(policy)

	var r0 []core.Asset
	if rf, ok := ret.Get(0).(func(uint32) []core.Asset); ok {
		r0 = rf(policy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Asset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint32) error); ok {
		r1 = rf(policy)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLoadAssetByCode provides a mock function with given fields: code
func (_m *mockAssetProvider) GetLoadAssetByCode(code string) (*core.Asset, error) {
	ret := _m.Called(code)

	var r0 *core.Asset
	if rf, ok := ret.Get(0).(func(string) *core.Asset); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Asset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
