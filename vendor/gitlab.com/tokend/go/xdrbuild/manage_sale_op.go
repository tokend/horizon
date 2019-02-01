package xdrbuild

import (
	"encoding/json"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
)

type (
	SaleDetails struct {
		Description string `json:"description"`
		Logo        struct {
			Key      string `json:"key"`
			MimeType string `json:"mime_type"`
			Name     string `json:"name"`
		} `json:"logo"`
		Name             string `json:"name"`
		ShortDescription string `json:"short_description"`
		YoutubeVideoID   string `json:"youtube_video_id"`
	}

	UpdateSaleDetails struct {
		SaleID         uint64
		NewSaleDetails SaleDetails
	}
)

func (u UpdateSaleDetails) Validate() error {
	return ValidateStruct(&u,
		Field(&u.SaleID, Required),
		Field(&u.NewSaleDetails, Required),
	)
}

func (u UpdateSaleDetails) XDR() (*xdr.Operation, error) {
	details, err := json.Marshal(u.NewSaleDetails)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal details")
	}

	op := &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageSale,
			ManageSaleOp: &xdr.ManageSaleOp{
				SaleId: xdr.Uint64(u.SaleID),
				Data: xdr.ManageSaleOpData{
					Action: xdr.ManageSaleActionCreateUpdateDetailsRequest,
					UpdateSaleDetailsData: &xdr.UpdateSaleDetailsData{
						RequestId:  0,
						NewDetails: xdr.Longstring(details),
					},
				},
			},
		},
	}

	return op, nil
}