package config

import (
	"bullioncoin.githost.io/development/go/keypair"
	"errors"
)

type Core struct {
	*Base
	AccountManagerKey string
	DemurrageOperator string
	AssetRateOperator string
}

func (c *Core) DefineConfigStructure() {
	c.bindEnv("account_manager_key")
	c.bindEnv("sequence_provider")
	c.bindEnv("demurrage_operator")
	c.bindEnv("asset_rate_operator")
}

func (c *Core) Init() error {
	var err error
	c.AccountManagerKey, err = c.getNonEmptyString("account_manager_key")
	if err != nil {
		return err
	}
	kp, err := keypair.Parse(c.AccountManagerKey)
	if err != nil {
		return errors.New("Could not parse Account Manager Key")
	}
	_, ok := kp.(*keypair.Full)
	if !ok {
		return errors.New("Account Manager Key must be a Secret Key, not Public Key")
	}

	c.DemurrageOperator = c.getString("demurrage_operator")
	if c.DemurrageOperator != "" {
		kp, err := keypair.Parse(c.DemurrageOperator)
		if err != nil {
			return errors.New("Could not parse Demurrage Operator Key")
		}
		_, ok := kp.(*keypair.Full)
		if !ok {
			return errors.New("Demurrage Operator Key must be a Secret Key, not Public Key")
		}
	}

	return nil
}
