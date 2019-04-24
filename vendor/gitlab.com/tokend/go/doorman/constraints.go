package doorman

import (
	"net/http"

	"gitlab.com/tokend/go/signcontrol"
)

func SignerOf(address string) SignerConstraint {
	return func(r *http.Request, doorman Doorman) error {
		signer, err := signcontrol.CheckSignature(r)
		if err != nil {
			return err
		}

		signers, err := doorman.AccountSigners(address)
		if err != nil {
			return err
		}

		for _, accountSigner := range signers {
			if accountSigner.AccountID == signer && accountSigner.Weight > 0 {
				return nil
			}
		}
		return ErrNotAllowed
	}
}

func SignatureOf(address string) SignerConstraint {
	return func(r *http.Request, doorman Doorman) error {
		signer, err := signcontrol.CheckSignature(r)
		if err != nil {
			return err
		}

		if signer == address {
			return nil
		}

		return ErrNotAllowed
	}
}
