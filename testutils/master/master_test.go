package master

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/client"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/testutils/common"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

func TestMasterIntegration(t *testing.T) {

	testCtx := SetupMasters(t, &common.TestMasterConfiguration{
		ReplicationFactor: 1,
		MasterPrefix:      "master-it",
	})
	defer testCtx.Cleanup()

	client := client.NewYBClient(&configs.YBClientConfig{
		MasterHostPort: testCtx.MasterExternalAddresses(),
		OpTimeout:      time.Duration(time.Second * 5),
	}, hclog.Default())

	common.Eventually(t, 15, func() error {
		if err := client.Connect(); err != nil {
			return err
		}
		return nil
	})

	defer client.Close()

	common.Eventually(t, 15, func() error {
		request := &ybApi.ListMastersRequestPB{}
		response := &ybApi.ListMastersResponsePB{}
		err := client.Execute(request, response)
		if err != nil {
			return err
		}
		t.Log(" ==> Received master list", response)
		return nil
	})

}
