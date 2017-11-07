package test

import (
	"gitlab.com/distributed_lab/tokend/horizon/log"
	"github.com/sirupsen/logrus"
)

var testLogger *log.Entry

func init() {
	testLogger, _ = log.New()
	testLogger.Entry.Logger.Formatter.(*logrus.TextFormatter).DisableColors = true
	testLogger.Entry.Logger.Level = logrus.DebugLevel
}
