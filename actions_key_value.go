package horizon

import (
	"gitlab.com/swarmfund/horizon/resource/keyvalue"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
)

type KeyValueShowAction struct {
	Action
	Key string
	KeyValue keyvalue.KeyValue
}

func (action *KeyValueShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			hal.Render(action.W, action.KeyValue)
		},
	)
}

func (action *KeyValueShowAction) loadParams(){
	action.Key = action.GetString("key");
}

func (action *KeyValueShowAction) loadRecord() {
	keyValueRecord,err := action.CoreQ().KeyValue().ByKey(action.Key)
	if err!=nil {
		action.Log.WithError(err).Error("Failed to get key_value from core DB")
		action.Err = &problem.ServerError
		return
	}

	if keyValueRecord == nil {
		action.Err = &problem.NotFound
		return
	}

	action.KeyValue.Populate(keyValueRecord)
}
