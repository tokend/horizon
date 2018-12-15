package operations

import (
	"testing"
	"gitlab.com/tokend/go/xdr"
)

func TestAllOperationsHandled(t *testing.T) {
	opsHandler := NewOperationsHandler(&mockOperationIDProvider{}, &mockParticipantEffectIDProvider{},
		&mockPublicKeyProvider{}, &mockBalanceProvider{})
	for _, opType := range xdr.OperationTypeAll {
		_, ok := opsHandler.allHandlers[opType]
		if !ok {
			t.Fatalf("All operations must be handled. Operation type: %s is not handled", opType)
		}
	}

}
