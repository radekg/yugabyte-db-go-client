package master

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/client/implementation"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/testutils/common"
)

func TestMasterIntegration(t *testing.T) {

	testCtx := SetupMasters(t, &common.TestMasterConfiguration{
		ReplicationFactor: 1,
		MasterPrefix:      "master-it",
	})
	defer testCtx.Cleanup()

	client, err := implementation.MasterLeaderConnectedClient(&configs.CliConfig{
		MasterHostPort: testCtx.MasterExternalAddresses(),
		OpTimeout:      time.Duration(time.Second * 5),
	}, hclog.Default())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	common.Eventually(t, 15, func() error {
		listMastersPb, err := client.ListMasters()
		if err != nil {
			return err
		}
		t.Log(" ==> Received master list", listMastersPb)
		return nil
	})

}
