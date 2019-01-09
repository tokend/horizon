package resource

import (
	"gitlab.com/tokend/horizon/db2/core"
)

type AccountCollection struct {
	Base

	Filters struct {
		AccountType string
	}

	resources []Account
	records   []core.Account
}

func NewAccountCollection() (*AccountCollection, error) {
	return &AccountCollection{}, nil
}

func (c *AccountCollection) Fetch() error {
	return c.CoreQ.Accounts().Select(&c.records)
}

func (c *AccountCollection) IsAllowed() (bool, error) {
	return c.isSignedByMaster(), nil
}

func (c *AccountCollection) Populate() error {
	for _, record := range c.records {
		resource, err := NewAccount(record.AccountID)
		if err != nil {
			return err
		}

		resource.record = &record

		err = resource.Populate()
		if err != nil {
			return err
		}

		c.resources = append(c.resources, *resource)
	}

	return nil
}

func (c *AccountCollection) Response() ([]Response, error) {
	response := make([]Response, len(c.resources))

	for i := range c.resources {
		r, err := c.resources[i].Response()
		if err != nil {
			return nil, err
		}

		response[i] = r
	}

	return response, nil
}
