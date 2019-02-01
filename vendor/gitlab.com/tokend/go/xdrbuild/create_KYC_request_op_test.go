package xdrbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/keypair"
)

func TestCreateKYCRequestOp_XDR(t *testing.T) {
	kp, _ := keypair.Random()
	var allTasks uint32 = 3
	t.Run("valid", func(t *testing.T) {
		op := CreateChangeRoleRequestOp{
			RequestID:          0,
			DestiantionAccount: kp.Address(),
			AccountRoleToSet:   1,
			KYCData:            "Some KYC data",
			AllTasks:           &allTasks,
		}
		assert.NoError(t, op.Validate())
		got, err := op.XDR()
		assert.NoError(t, err)
		body := got.Body.CreateChangeRoleRequestOp
		assert.EqualValues(t, op.RequestID, body.RequestId)
		assert.EqualValues(t, op.AccountRoleToSet, body.AccountRoleToSet)
		assert.EqualValues(t, op.KYCData, body.KycData)
		assert.EqualValues(t, op.DestiantionAccount, body.DestinationAccount.Address())
		assert.EqualValues(t, op.AllTasks, body.AllTasks)
	})

	t.Run("missing account type", func(t *testing.T) {
		op := CreateChangeRoleRequestOp{
			RequestID:          0,
			KYCData:            "Some KYC data",
			DestiantionAccount: kp.Address(),
			AllTasks:           nil,
		}
		assert.Error(t, op.Validate())
	})

	t.Run("missing KYC data", func(t *testing.T) {
		op := CreateChangeRoleRequestOp{
			RequestID:          0,
			AccountRoleToSet:   1,
			DestiantionAccount: kp.Address(),
			AllTasks:           nil,
		}
		assert.Error(t, op.Validate())
	})

	t.Run("missing updated account", func(t *testing.T) {
		op := CreateChangeRoleRequestOp{
			RequestID:        0,
			AccountRoleToSet: 1,
			KYCData:          "Some KYC data",
			AllTasks:         nil,
		}
		assert.Error(t, op.Validate())
	})
}
