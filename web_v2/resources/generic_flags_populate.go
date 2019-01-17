package resources

import "gitlab.com/tokend/regources/v2"

//go:generate genny -in=$GOFILE -out=signer_type_flags_populate.go gen "flagValueType=xdr.SignerType"
//go:generate genny -in=$GOFILE -out=asset_policy_type_flags_populate.go gen "flagValueType=xdr.AssetPolicy"
//go:generate genny -in=$GOFILE -out=block_reasons_flags_populate.go gen "flagValueType=xdr.BlockReasons"

//FlagFromflagValueType - unwraps mask into set of flags
func FlagFromflagValueType(mask int32, allFlags []flagValueType) []regources.Flag {
	var result []regources.Flag
	for _, flagValue := range allFlags {
		flagValueAsInt := int32(flagValue)
		if (flagValueAsInt & mask) == flagValueAsInt {
			result = append(result, regources.Flag{
				Value: flagValueAsInt,
				Name:  flagValue.ShortString(),
			})
		}
	}

	return result
}
