package test

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/tokend/horizon/log"
)

var testLogger *log.Entry

func init() {
	testLogger, _ = log.New()
	testLogger.Entry.Formatter(&logrus.TextFormatter{DisableColors: true})
	testLogger.Entry.Level(log.DebugLevel)
}
