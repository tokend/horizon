// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package base

import "gitlab.com/tokend/go/xdr"

func FlagFromXdrBlockReasons(mask int32, allFlags []xdr.BlockReasons) []Flag {
	result := []Flag{}
	for _, flagValue := range allFlags {
		flagValueAsInt := int32(flagValue)
		if (flagValueAsInt & mask) == flagValueAsInt {
			result = append(result, Flag{
				Value: flagValueAsInt,
				Name:  flagValue.String(),
			})
		}
	}

	return result
}
