package storage

import (
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/lann/squirrel"
	"gitlab.com/tokend/horizon/db2"
)

const (
	ChanSubmitter = "submitter"
	ChanNewLedger = "new_ledger"
)

func notifyListeners(repo *db2.Repo, channel string) {
	_, err := repo.Exec(sq.Expr(fmt.Sprintf("NOTIFY %s", channel)))
	if err != nil {
		panic(errors.Wrap(err, "failed to notify channel", logan.F{
			"channel": channel,
		}))
	}
}
