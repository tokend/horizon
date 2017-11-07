package horizon

import (
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/resource"
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
