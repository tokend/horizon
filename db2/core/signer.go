package core

import (
	sq "github.com/lann/squirrel"
)

// Signer is a row of data from the `signers` table from stellar-core
type Signer struct {
	Accountid  string
	Publickey  string
	Weight     int32
	SignerType int32  `db:"signer_type"`
	Identity   int32  `db:"identity_id"`
	Name       string `db:"signer_name"`
}

func (s Signer) GetPublicKey() string {
	return s.Publickey
}

func (s Signer) GetIdentity() int32 {
	return s.Identity
}

func (s Signer) GetName() string {
	return s.Name
}

// SignersByAddress loads all signer rows for `addy`
func (q *Q) SignersByAddress(dest interface{}, addy string) error {
	sql := selectSigner.Where("accountid = ?", addy)
	return q.Select(dest, sql)
}

var selectSigner = sq.Select(
	"si.accountid",
	"si.publickey",
	"si.weight",
	"si.signer_type",
	"si.identity_id",
	"si.signer_name",
).From("signers si").OrderBy("identity_id DESC")
