package config

import (
	"gitlab.com/swarmfund/go/keypair"
	"errors"
)

type Core struct {
	*Base
	DemurrageOperator string
}

func (c *Core) DefineConfigStructure() {
	c.bindEnv("demurrage_operator")
}

func (c *Core) Init() error {
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
