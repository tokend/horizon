/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type LicenseInfo struct {
	Key
	Attributes LicenseInfoAttributes `json:"attributes"`
}
type LicenseInfoResponse struct {
	Data     LicenseInfo `json:"data"`
	Included Included    `json:"included"`
}

type LicenseInfoListResponse struct {
	Data     []LicenseInfo `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustLicenseInfo - returns LicenseInfo from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLicenseInfo(key Key) *LicenseInfo {
	var licenseInfo LicenseInfo
	if c.tryFindEntry(key, &licenseInfo) {
		return &licenseInfo
	}
	return nil
}
