package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/resource/keyvalue"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type KeyValueShowAction struct {
	Action
	key string
	keyValueRecord *core.KeyValue
}

func (action *KeyValueShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			var res keyvalue.KeyValue
			err := res.Populate(action.keyValueRecord)
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
	action.key = action.GetString("key")
}

func (action *KeyValueShowAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.keyValueRecord, err = q.ByKey(action.key)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get key_value from core DB")
		action.Err = &problem.ServerError
		return
	}

	if action.keyValueRecord == nil {
		action.Err = &problem.NotFound
		return
	}
}

type KeyValueShowAllAction struct {
	Action
	coreRecords 	[]core.KeyValue
	recordsToRender []keyvalue.KeyValue
}

func (action *KeyValueShowAllAction) JSON() {
	action.Do(
		action.loadRecord,
		func() {
			hal.Render(action.W, action.recordsToRender)
		},
	)
}

func (action *KeyValueShowAllAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.coreRecords, err = q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get all key_values from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.recordsToRender, err = action.getPopulatedKeyValues()
	if err != nil {
		action.Log.WithError(err).Error("Failed to populate all key_values")
		action.Err = &problem.ServerError
		return
	}
}

func (action *KeyValueShowAllAction) getPopulatedKeyValues() ([]keyvalue.KeyValue, error){
	var res []keyvalue.KeyValue
	for i, keyValue := range action.coreRecords {
		res = append(res, keyvalue.KeyValue{})
		err := res[i].Populate(&keyValue)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to populate key_value", logan.F{
				"keyValue": keyValue,
			})
		}
	}
	return res, nil
}