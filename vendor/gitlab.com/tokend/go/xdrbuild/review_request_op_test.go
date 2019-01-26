package xdrbuild

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/xdr"
)

func TestReviewRequestOp_XDR(t *testing.T) {
	hash := [32]byte{
		0xde, 0xad, 0xbe, 0xaf, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	t.Run("approve", func(t *testing.T) {
		op := ReviewRequestOp{
			ID:     1,
			Hash:   fmt.Sprintf("%x", hash[:]),
			Action: xdr.ReviewRequestOpActionApprove,
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
	})

	t.Run("approve with tasks", func(t *testing.T) {
		reviewDetails := ReviewDetails{
			TasksToAdd:      0,
			TasksToRemove:   8,
			ExternalDetails: "Request can be approved",
		}
		op := ReviewRequestOp{
			ID:            1,
			Hash:          fmt.Sprintf("%x", hash[:]),
			Action:        xdr.ReviewRequestOpActionApprove,
			ReviewDetails: &reviewDetails,
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
		assert.EqualValues(t, op.ReviewDetails.TasksToAdd, body.ReviewDetails.TasksToAdd)
		assert.EqualValues(t, op.ReviewDetails.TasksToRemove, body.ReviewDetails.TasksToRemove)
		assert.EqualValues(t, op.ReviewDetails.ExternalDetails, body.ReviewDetails.ExternalDetails)
	})

	t.Run("approve atomic swap with tasks", func(t *testing.T) {
		reviewDetails := ReviewDetails{
			TasksToAdd:      0,
			TasksToRemove:   8,
			ExternalDetails: "Request can be approved",
		}
		op := ReviewRequestOp{
			ID:            1,
			Hash:          fmt.Sprintf("%x", hash[:]),
			Details:       AtomicSwapDetails{},
			Action:        xdr.ReviewRequestOpActionApprove,
			ReviewDetails: &reviewDetails,
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, xdr.ReviewableRequestTypeAtomicSwap, body.RequestDetails.RequestType)
		assert.EqualValues(t, op.Action, body.Action)
		assert.EqualValues(t, op.ReviewDetails.TasksToAdd, body.ReviewDetails.TasksToAdd)
		assert.EqualValues(t, op.ReviewDetails.TasksToRemove, body.ReviewDetails.TasksToRemove)
		assert.EqualValues(t, op.ReviewDetails.ExternalDetails, body.ReviewDetails.ExternalDetails)
	})

	t.Run("reject", func(t *testing.T) {
		op := ReviewRequestOp{
			ID:     1,
			Hash:   fmt.Sprintf("%x", hash[:]),
			Action: xdr.ReviewRequestOpActionReject,
			Reason: "yoba",
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
		assert.EqualValues(t, op.Reason, body.Reason)
	})

	t.Run("approve withdraw", func(t *testing.T) {
		op := ReviewRequestOp{
			ID:     1,
			Hash:   fmt.Sprintf("%x", hash[:]),
			Action: xdr.ReviewRequestOpActionApprove,
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
	})

	t.Run("reject withdraw", func(t *testing.T) {
		op := ReviewRequestOp{
			ID:     1,
			Hash:   fmt.Sprintf("%x", hash[:]),
			Action: xdr.ReviewRequestOpActionPermanentReject,
			Reason: "yoba",
			Details: WithdrawalDetails{
				ExternalDetails: "foobar",
			},
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
		assert.EqualValues(t, xdr.ReviewableRequestTypeWithdraw, body.RequestDetails.RequestType)
		assert.NotNil(t, body.RequestDetails.Withdrawal)
		assert.EqualValues(t, "foobar", body.RequestDetails.Withdrawal.ExternalDetails)
	})

	t.Run("approve update kyc", func(t *testing.T) {
		op := ReviewRequestOp{
			ID:     1,
			Hash:   fmt.Sprintf("%x", hash[:]),
			Action: xdr.ReviewRequestOpActionApprove,
			Details: UpdateKYCDetails{
				ExternalDetails: "All right",
			},
			ReviewDetails: &ReviewDetails{
				TasksToAdd:    0,
				TasksToRemove: 3,
			},
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
		assert.EqualValues(t, xdr.ReviewableRequestTypeUpdateKyc, body.RequestDetails.RequestType)
		assert.NotNil(t, body.RequestDetails.UpdateKyc)
		assert.EqualValues(t, "All right", body.RequestDetails.UpdateKyc.ExternalDetails)
	})

	t.Run("reject update kyc", func(t *testing.T) {
		op := ReviewRequestOp{
			ID:     1,
			Hash:   fmt.Sprintf("%x", hash[:]),
			Action: xdr.ReviewRequestOpActionReject,
			Reason: "yoba",
			Details: UpdateKYCDetails{
				ExternalDetails: "Invalid identity",
			},
			ReviewDetails: &ReviewDetails{
				TasksToAdd:    2,
				TasksToRemove: 0,
			},
		}
		got, err := op.XDR()
		if err != nil {
			t.Fatal(err)
		}
		body := got.Body.ReviewRequestOp
		assert.EqualValues(t, op.ID, body.RequestId)
		assert.EqualValues(t, hash, body.RequestHash)
		assert.EqualValues(t, op.Action, body.Action)
		assert.EqualValues(t, op.Reason, body.Reason)
		assert.EqualValues(t, xdr.ReviewableRequestTypeUpdateKyc, body.RequestDetails.RequestType)
		assert.NotNil(t, body.RequestDetails.UpdateKyc)
		assert.EqualValues(t, "Invalid identity", body.RequestDetails.UpdateKyc.ExternalDetails)
	})
}
