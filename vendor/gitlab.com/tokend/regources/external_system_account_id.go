package regources

import "gitlab.com/tokend/regources/valueflag"

type ExternalSystemAccountID struct {
	Type      valueflag.Flag `json:"type"`
	Data      string         `json:"data"`
	AssetCode string         `json:"asset_code,omitempty"`
	ExpiresAt *string        `json:"expires_at,omitempty"`
}
