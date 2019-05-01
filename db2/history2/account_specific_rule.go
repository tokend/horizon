package history2

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/go/xdr"
)

// AccountSpecificRule is a row of data from the `account_specific_rules` table
type AccountSpecificRule struct {
	ID        uint64        `db:"id"`
	Address   *string       `db:"address"`
	Forbids   bool          `db:"forbids"`
	EntryType int32         `db:"entry_type"`
	Key       xdr.LedgerKey `db:"key"`
}

func NewAccountSpecificRule(entry xdr.AccountSpecificRuleEntry) AccountSpecificRule {
	address := entry.AccountId.Address()
	var addressPtr *string
	if address != "" {
		addressPtr = &address
	}

	return AccountSpecificRule{
		ID:        uint64(entry.Id),
		Key:       entry.LedgerKey,
		Address:   addressPtr,
		EntryType: int32(entry.LedgerKey.Type),
	}
}

type LedgerKey xdr.LedgerKey

//Value - implements db driver method for auto marshal
func (r LedgerKey) Value() (driver.Value, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal poll data")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *LedgerKey) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for jsonb")
	}

	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ledger key")
	}

	return nil
}
