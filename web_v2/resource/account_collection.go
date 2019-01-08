package resource

import (
	"gitlab.com/tokend/horizon/db2/core"
)

type AccountCollection struct {
	Base `json:"-"`

	Filters struct {
		AccountType string
	}

	resources []Account
	records   []core.Account
}

func (c *AccountCollection) Fetch(pp PagingParams) error {
	return c.CoreQ.Accounts().Select(&c.records)
}

func (c *AccountCollection) IsAllowed() (bool, error) {
	return c.isSignedByMaster(), nil
}

func (c *AccountCollection) Populate() error {
	for _, r := range c.resources {
		err := r.Populate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *AccountCollection) Response() ([]interface{}, error) {
	response := make([]interface{}, len(c.resources))

	for i := range c.resources {
		response[i] = c.resources[i]
	}

	return response, nil
}
