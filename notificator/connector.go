package notificator

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"gitlab.com/distributed_lab/notificator"
	"gitlab.com/swarmfund/go/hash"
	"gitlab.com/swarmfund/horizon/config"
	"gitlab.com/swarmfund/horizon/log"
)

const (
	NotificatorTypeVerificationEmail = 3
	NotificatorTypeApprovalEmail     = 4
	NotificatorTypeTFA               = 5
	NotificatorTypeRecoveryRequest   = 6
	NotificatorTypeAdminNotification = 7
)

type Connector struct {
	notificator *notificator.Connector
	conf        *config.Notificator

	log *log.Entry
}

func NewConnector(conf *config.Notificator) *Connector {
	return &Connector{
		notificator: notificator.NewConnector(
			notificator.Pair{Secret: conf.Secret, Public: conf.Public},
			conf.GetEndpoint(),
		),
		conf: conf,

		log: log.WithField("service", "notificator"),
	}
}

func (c *Connector) SendTFA(walletID string, phone string, otp string) (*time.Duration, error) {
	payload := &notificator.SMSRequestPayload{
		Destination: phone,
		Message:     fmt.Sprintf("%s is your BullionCoin authentication code.", otp),
	}

	walletIDHash := hash.Hash([]byte(walletID))
	userToken := base64.StdEncoding.EncodeToString(walletIDHash[:])
	response, err := c.notificator.Send(NotificatorTypeTFA, userToken, payload)
	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		if response.IsPermanent() {
			return nil, fmt.Errorf("failed to send sms")
		}

		return response.RetryIn(), nil
	}

	return nil, nil
}

func (c *Connector) send(requestType int, token string, payload notificator.Payload) error {
	response, err := c.notificator.Send(requestType, token, payload)
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return errors.New("notification request not accepted")
	}
	return nil
}
