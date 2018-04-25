package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/resource/keyvalue"
)

type KeyValueShowAction struct {
	Action
	Key string
	KeyValueRecord *core.KeyValue
}

func (action *KeyValueShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			var res keyvalue.KeyValue
			err := res.Populate(action.KeyValueRecord)
			if err != nil {
				action.Log.WithError(err).Error("Failed to populate key_value")
				action.Err = &problem.ServerError
				return
			}
			hal.Render(action.W, res)
		},
	)
}

func (action *KeyValueShowAction) loadParams(){
	action.Key = action.GetString("key");
}

func (action *KeyValueShowAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.KeyValueRecord,err = q.ByKey(action.Key)
	if err!=nil {
		action.Log.WithError(err).Error("Failed to get key_value from core DB")
		action.Err = &problem.ServerError
		return
	}

	if action.KeyValueRecord == nil {
		action.Err = &problem.NotFound
		return
	}
}
