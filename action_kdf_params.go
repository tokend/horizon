package horizon

import (
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/resource"
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
