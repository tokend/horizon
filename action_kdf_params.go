package horizon

import (
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
)

type KdfParamsAction struct {
	Action
}

func (action *KdfParamsAction) JSON() {
	action.ValidateBodyType()
	action.Do(func() {
		var response resource.KdfParams
		response.Populate()
		hal.Render(action.W, response)
	})
}
