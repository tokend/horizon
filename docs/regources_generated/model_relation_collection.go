/*
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package regources

import "encoding/json"

type RelationCollection struct {
	Data []map[string]interface{} `json:"data"`
}

func (r RelationCollection) MarshalJSON() ([]byte, error) {
	if r.Data == nil {
		r.Data = []Key{}
	}

	type temp RelationCollection
	return json.Marshal(temp(r))
}
