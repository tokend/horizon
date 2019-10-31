package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// LicenseQ is a helper struct to aid in configuring queries that loads
// license structs.
type LicenseQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewLicenseQ - creates new instance of LicenseQ
func NewLicenseQ(repo *db2.Repo) LicenseQ {
	return LicenseQ{
		repo: repo,
		selector: sq.Select("license.id",
			"license.hash",
			"license.prev_hash",
			"license.ledger_hash",
			"license.admin_count",
			"license.due_date",
		).From("license license"),
	}
}

func (q LicenseQ) GetLatest() (*License, error) {
	q.selector = q.selector.OrderBy("license.id desc").Limit(1)
	return q.Get()
}

// Get - loads a row from `license`
// returns nil, nil - if license does not exists
// returns error if more than one License found
func (q LicenseQ) Get() (*License, error) {
	var result License
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load license")
	}

	return &result, nil
}

//Select - selects slice from the db, if no license found - returns nil, nil
func (q LicenseQ) Select() ([]License, error) {
	var result []License
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load license license")
	}

	return result, nil
}
