/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type SaleAccessDefinitionType int

const (
	SaleAccessDefinitionTypeNone SaleAccessDefinitionType = iota + 1
	SaleAccessDefinitionTypeWhitelist
	SaleAccessDefinitionTypeBlacklist
)

var saleAccessDefinitionTypeMap = map[SaleAccessDefinitionType]string{
	SaleAccessDefinitionTypeNone:      "none",
	SaleAccessDefinitionTypeWhitelist: "whitelist",
	SaleAccessDefinitionTypeBlacklist: "blacklist",
}

func (s SaleAccessDefinitionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  saleAccessDefinitionTypeMap[s],
		Value: int32(s),
	})
}

//String - converts int enum to string
func (s SaleAccessDefinitionType) String() string {
	return saleAccessDefinitionTypeMap[s]
}
