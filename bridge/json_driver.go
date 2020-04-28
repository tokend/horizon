package bridge

import (
	"database/sql/driver"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func DriverValue(data interface{}) (driver.Value, error) {
	data, err := pgdb.JSONValue(data)
	return data, err
}

func DriveScan(src, dest interface{}) error {
	return pgdb.JSONScan(src, dest)
}
