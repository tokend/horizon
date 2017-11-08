package notificator

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"gitlab.com/tokend/go/hash"
	"gitlab.com/tokend/horizon/config"
	"gitlab.com/tokend/horizon/log"
	"github.com/go-errors/errors"
	"gitlab.com/distributed_lab/notificator"
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

func (c *Connector) SendLowAvailableEmissions(email, asset string) error {
	letter := Letter{
		Header: "BullionCoin Admin Notification",
	}

	letter.Body = fmt.Sprintf(`Asset %s has low emission. Upload more presigned emissions.`, asset)

	var buff bytes.Buffer
	err := c.conf.AdminNotification.Template.Execute(&buff, letter)
	if err != nil {
		log.WithField("err", err.Error()).Error("failed to render template")
		return err
	}

	payload := &notificator.EmailRequestPayload{
		Destination: email,
		Subject:     letter.Header,
		Message:     buff.String(),
	}

	response, err := c.notificator.Send(NotificatorTypeAdminNotification, email, payload)
	if err != nil {
		c.log.WithError(err).Error("Failed to SendLowAvailableEmissions")
		return err
	}

	if !response.IsSuccess() {
		if retryIn := response.RetryIn(); retryIn != nil {
			return nil
		}
		c.log.WithField("response", response).Warn("Low Available Emissions notification not accepted")
		return errors.New("Low Available Emissions notification not accepted")
	}

	return nil
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
