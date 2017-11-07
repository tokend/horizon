package horizon

import (
	"encoding/json"
	"testing"

	"gitlab.com/distributed_lab/tokend/horizon/resource"
	"gitlab.com/distributed_lab/tokend/horizon/test"
)

func TestRootAction(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	server := test.NewStaticMockServer(`{
			"info": {
				"network": "test",
				"build": "test-core"
			}
		}`)
	defer server.Close()

	ht.App.horizonVersion = "test-horizon"
	ht.App.config.StellarCoreURL = server.URL

	w := ht.Get("/")
	if ht.Assert.Equal(200, w.Code) {
		var actual resource.Root
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)
		ht.Assert.Equal("test-horizon", actual.HorizonVersion)
		ht.Assert.Equal("test-core", actual.StellarCoreVersion)
	}
}
