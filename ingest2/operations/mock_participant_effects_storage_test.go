// Code generated by mockery v1.0.1. DO NOT EDIT.

package operations

import (
	mock "github.com/stretchr/testify/mock"
	history2 "gitlab.com/tokend/horizon/db2/history2"
)

// mockParticipantEffectsStorage is an autogenerated mock type for the participantEffectsStorage type
type mockParticipantEffectsStorage struct {
	mock.Mock
}

// Insert provides a mock function with given fields: participants
func (_m *mockParticipantEffectsStorage) Insert(participants []history2.ParticipantEffect) error {
	ret := _m.Called(participants)

	var r0 error
	if rf, ok := ret.Get(0).(func([]history2.ParticipantEffect) error); ok {
		r0 = rf(participants)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
