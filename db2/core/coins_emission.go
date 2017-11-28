package core

type CoinsEmission struct {
	SerialNumber string `db:"serial_number"`
	Amount       int64  `db:"amount"`
	LastModified int64  `db:"lastmodified"`
}

type AssetAmount struct {
	Asset  string `db:"asset"`
	Amount int64  `db:"amount"`
}

type AssetStat struct {
	Asset     string `db:"asset"`
	Hundreds  int64  `db:"hundreds"`
	Ones      int64  `db:"ones"`
	Remainder int64  `db:"remainder"`
}
