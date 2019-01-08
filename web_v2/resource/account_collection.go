package resource

import "gitlab.com/tokend/horizon/db2/core"

type AccountCollection struct {
	Base `json:"-"`

	resources []Account
	records   []core.Account
}

func (c *AccountCollection) Fetch() error {
	return nil
}

func (c *AccountCollection) IsAllowed() (bool, error) {
	return c.isSignedByAdmin(), nil
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

func (c *AccountCollection) Response() (interface{}, error) {
	return c.resources, nil
}

