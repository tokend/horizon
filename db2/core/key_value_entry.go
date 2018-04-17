package core

import "gitlab.com/swarmfund/go/xdr"

type KeyValue struct {
	Key 				string 					`db:"key_value_key"`
	KeyValueType 		xdr.KeyValueEntryType 	`db:"key_value_type"`
	KeyValueDetails 	*KeyValueDetails 		`db:"key_value_details"`
}