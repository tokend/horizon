package core

// AccountKYC is a row of data from the `account_KYC` table
type AccountKYC struct {
	AccountID string `db:"account_id"`
	KYCData   string `db:"KYC_data"`
}
