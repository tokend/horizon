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
	action.Key = action.GetString("key")
}

func (action *KeyValueShowAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.KeyValueRecord, err = q.ByKey(action.Key)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get key_value from core DB")
		action.Err = &problem.ServerError
		return
	}

	if action.KeyValueRecord == nil {
		action.Err = &problem.NotFound
		return
	}
}

type KeyValueShowAllAction struct {
	Action
	CoreRecord []core.KeyValue
}

func (action *KeyValueShowAllAction) JSON() {
	action.Do(
		action.loadRecord,
		func() {
			res, err := action.getPopulatedKeyValues()
			if err != nil {
				action.Log.WithError(err).Error("Failed to populate all key_values")
				action.Err = &problem.ServerError
				return
			}

			hal.Render(action.W, res)
		},
	)
}

func (action *KeyValueShowAllAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.CoreRecord, err = q.All()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get all key_values from core DB")
		action.Err = &problem.ServerError
		return
	}

	if action.CoreRecord == nil {
		action.Err = &problem.NotFound
		return
	}
}

func (action *KeyValueShowAllAction) getPopulatedKeyValues() ([]keyvalue.KeyValue, error){
	var res []keyvalue.KeyValue
	for i, keyValue := range action.CoreRecord {
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