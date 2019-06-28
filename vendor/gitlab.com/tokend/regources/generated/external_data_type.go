/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type ExternalDataType int

const (
	ExternalDataTypeAddress ExternalDataType = iota + 1
	ExternalDataTypeAddressWithPayload
)

var externalDataTypeMap = map[ExternalDataType]string{
	ExternalDataTypeAddress:            "address",
	ExternalDataTypeAddressWithPayload: "address_with_payload",
}

func (s ExternalDataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  externalDataTypeMap[s],
		Value: int32(s),
	})
}

//String - converts int enum to string
func (s ExternalDataType) String() string {
	return externalDataTypeMap[s]
}

func (s *ExternalDataType) UnmarshalJSON(b []byte) error {
	var res Flag
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	*s = ExternalDataType(res.Value)
	return nil
}
