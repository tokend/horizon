package config

import (
	"gitlab.com/tokend/horizon/log"
	"net/url"
)

type Notificator struct {
	*Base
	Endpoint          string
	AdminNotification AdminNotification
	Secret            string
	Public            string
}

func (n *Notificator) DefineConfigStructure() {
	n.bindEnv("endpoint")
	n.bindEnv("secret")
	n.bindEnv("public")

	n.AdminNotification.Base = NewBase(n.Base, "admin_notification")
	n.AdminNotification.DefineConfigStructure()
}

func (n *Notificator) Init() error {
	var err error
	n.Endpoint, err = n.getURLAsString("endpoint")
	if err != nil {
		return err
	}

	n.Secret, err = n.getNonEmptyString("secret")
	if err != nil {
		return err
	}

	n.Public, err = n.getNonEmptyString("public")
	if err != nil {
		return err
	}

	err = n.AdminNotification.Init()
	if err != nil {
		return err
	}

	return nil
}

func (n *Notificator) GetEndpoint() url.URL {
	result, err := url.Parse(n.Endpoint)
	if err != nil {
		log.WithField("service", "conf_notificator").WithError(err).Fatal("Failed to get endpoint")
	}

	return *result
}
