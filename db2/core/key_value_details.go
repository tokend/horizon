package core

import (
	_"time"

	"database/sql/driver"
	"gitlab.com/distributed_lab/logan/v3/errors"
	_"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
)

type KeyValueDetails struct {
	KYCSettings *KycSettings `json:"kyc_settings"`
}

func (r KeyValueDetails) Value() (driver.Value, error) {
	result, err := db2.DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	return result, nil
}

func (r *KeyValueDetails) Scan(src interface{}) error {
	err := db2.DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details")
	}

	return nil
}


type KycSettings struct {

}