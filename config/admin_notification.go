package config

import (
	"html/template"

	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/horizon/log"
)

type AdminNotification struct {
	*Base
	Template                     *template.Template
	EmissionThreshold            int64
	EmissionNotificationReceiver string
}

func (a *AdminNotification) DefineConfigStructure() {
	a.bindEnv("emission_threshold")
	a.bindEnv("receiver")
}

func (a *AdminNotification) Init() error {
	var err error
	a.Template = a.getTemplate("admin_notification")

	emissionThreshold, err := amount.Parse(a.getString("emission_threshold"))
	if err != nil {
		log.Error("Unable to parse emission threshold")
		return err
	}
	a.EmissionThreshold = int64(emissionThreshold * amount.One)
	a.EmissionNotificationReceiver, err = a.getNonEmptyString("receiver")
	if err != nil {
		log.Error("Unable to parse emission notification receiver")
		return err
	}
	return nil
}
