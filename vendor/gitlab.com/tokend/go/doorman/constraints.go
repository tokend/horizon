package doorman

import (
	"fmt"
	"net/http"

	"gitlab.com/tokend/go/signcontrol"
)

func SignerOf(address string) SignerConstraint {
	return func(r *http.Request, doorman Doorman) error {
		signer, err := signcontrol.CheckSignature(r)
		if err != nil {
			fmt.Printf("%s: signature is not valid: %v\n", r.URL.String(), err)
			return err
		}

		if signer == address {
			return nil
		}
		signers, err := doorman.AccountSigners(address)
		if err != nil {
			fmt.Printf("%s: failed to get account signers: %v\n", r.URL.String(), err)
			return err
		}

		for _, accountSigner := range signers {
			fmt.Printf("%s: %s %s %d \n", r.URL.String(), accountSigner.AccountID, signer, accountSigner.Weight)
			if accountSigner.AccountID == signer && accountSigner.Weight > 0 {
				return nil
			}
		}
		fmt.Printf("%s: you die\n", r.URL.String())
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
