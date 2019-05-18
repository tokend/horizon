package history2

import (
	"encoding/json"
	"gitlab.com/tokend/regources/generated"
)

type SaleAccessDefinitionType int

const (
	SaleAccessDefinitionTypeNone SaleAccessDefinitionType = iota + 1
	SaleAccessDefinitionTypeWhitelist
	SaleAccessDefinitionTypeBlacklist
)

var saleAccessDefinitionTypeMap = map[SaleAccessDefinitionType]string{
	SaleAccessDefinitionTypeNone:      "open",
	SaleAccessDefinitionTypeWhitelist: "closed",
	SaleAccessDefinitionTypeBlacklist: "canceled",
}

func (s SaleAccessDefinitionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(regources.Flag{
		Name:  saleAccessDefinitionTypeMap[s],
		Value: int32(s),
	})
}

//String - converts int enum to string
func (s SaleAccessDefinitionType) String() string {
	return saleAccessDefinitionTypeMap[s]
}
