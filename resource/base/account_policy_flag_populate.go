package base

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func FlagFromXdrAccountPolicy(mask int32, allFlags []xdr.AccountPolicies) []regources.Flag {
	result := []regources.Flag{}
	for _, flagValue := range allFlags {
		flagValueAsInt := int32(flagValue)
		if (flagValueAsInt & mask) == flagValueAsInt {
			result = append(result, regources.Flag{
				Value: flagValueAsInt,
				Name:  flagValue.String(),
			})
		}
	}

	return result
}
