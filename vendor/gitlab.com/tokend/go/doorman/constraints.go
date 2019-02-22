package doorman

import (
	"net/http"

	"gitlab.com/tokend/go/signcontrol"
)

func SignerOf(address string) SignerConstraint {
	return func(r *http.Request, doorman Doorman) error {
		signer, err := signcontrol.CheckSignature(r)
		if err != nil {
			fmt.Printf("DEBUG: signer of %s: CheckSignature: %v\n", address, err)
			return err
		}

		if signer == address {
			return nil
		}
		signers, err := doorman.AccountSigners(address)
		if err != nil {
			fmt.Printf("DEBUG: signer of %s: AccountSigners: %v\n", address, err)
			return err
		}

		for _, accountSigner := range signers {
			if accountSigner.AccountID == signer && accountSigner.Weight > 0 {
				return nil
			}
		}
		fmt.Printf("DEBUG: signer of %s: NotASigner\n", address)
		return signcontrol.ErrNotAllowed
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

		return signcontrol.ErrNotAllowed
	}
}
