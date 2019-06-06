/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type KYCRecoveryStatus int

const (
	KYCRecoveryStatusNone KYCRecoveryStatus = iota
	KYCRecoveryStatusOngoing
)

var kycRecoveryStatusStr = map[KYCRecoveryStatus]string{
	KYCRecoveryStatusNone:    "none",
	KYCRecoveryStatusOngoing: "ongoing",
}

func (s KYCRecoveryStatus) String() string {
	return kycRecoveryStatusStr[s]
}

func (s KYCRecoveryStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  kycRecoveryStatusStr[s],
		Value: int32(s),
	})
}
