package base

import (
	"gitlab.com/tokend/regources/valueflag"
)

//go:generate genny -in=$GOFILE -out=signer_type_flags_populate.go gen "flagValueType=xdr.SignerType"
//go:generate genny -in=$GOFILE -out=asset_policy_type_flags_populate.go gen "flagValueType=xdr.AssetPolicy"
//go:generate genny -in=$GOFILE -out=block_reasons_flags_populate.go gen "flagValueType=xdr.BlockReasons"
//go:generate genny -in=$GOFILE -out=contract_state_flags_populate.go gen "flagValueType=xdr.ContractState"

func FlagFromflagValueType(mask int32, allFlags []flagValueType) []valueflag.Flag {
	result := []valueflag.Flag{}
	for _, flagValue := range allFlags {
		flagValueAsInt := int32(flagValue)
		if (flagValueAsInt & mask) == flagValueAsInt {
			result = append(result, valueflag.Flag{
				Value: flagValueAsInt,
				Name:  flagValue.String(),
			})
		}
	}

	return result
}
