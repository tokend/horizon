package codes

import (
	"encoding/json"
	"testing"

	"gitlab.com/tokend/go/xdr"
)

func TestCodeProviders(t *testing.T) {
	rawTx := "AAAAAD3w0GHAzCqJfiKeMESG6h6zyi4voJv9i5aDBOUjd4osAAAAAAAAAAAAAAAAAAAAAAAAAABcepj9AAAAAAAAAAEAAAAAAAAAEgAAAAAAAAAO8d+oN6I3bsEziNvff3+lIB6iyBdlHoIXD4RdOo44kL8AAAADAAAAAwAAAAhhc2RhIHNkYQAAAAAAAAAAAAAAAnt9AAAAAAAAAAAAAAAAAAAAAAABI3eKLAAAAEDOGMFv7GM/2cWqtFJyaFSqoIHMTzLh7pkpKxc2gLEoN+JlEhYwyo0X8YhDLWRIxooBhOUlatoD46jmpn9QDEUC"
	var meta xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(rawTx, &meta)
	if err != nil {
		panic(err)
	}

	asJSON, err := json.Marshal(&meta)
	if err != nil {
		panic(err)
	}

	println(string(asJSON))
}
