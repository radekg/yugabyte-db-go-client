package master

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/radekg/yugabyte-db-go-client/client"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/testutils/common"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"github.com/stretchr/testify/assert"
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
		t.Log("Received master list", response)
		return nil
	})

}

func TestMasterReconnect(t *testing.T) {

	request := &ybApi.ListMastersRequestPB{}

	testCtx := SetupMasters(t, &common.TestMasterConfiguration{
		ReplicationFactor: 1,
		MasterPrefix:      "master-it",
	})
	defer testCtx.Cleanup()

	client := client.NewYBClient(&configs.YBClientConfig{
		MasterHostPort:         testCtx.MasterExternalAddresses(),
		OpTimeout:              time.Duration(time.Second * 5),
		MaxReconnectAttempts:   1,
		ReconnectRetryInterval: time.Duration(time.Millisecond * 100),
	}, hclog.Default())

	errNotConnected := client.Execute(request, &ybApi.ListMastersResponsePB{})
	assert.NotNil(t, errNotConnected, "expected an error")

	common.Eventually(t, 15, func() error {
		if err := client.Connect(); err != nil {
			return err
		}
		return nil
	})

	defer client.Close()

	common.Eventually(t, 15, func() error {

		response := &ybApi.ListMastersResponsePB{}
		err := client.Execute(request, response)
		if err != nil {
			return err
		}
		t.Log("Received master list", response)
		return nil
	})

	testCtx.Cleanup()

	response := &ybApi.ListMastersResponsePB{}
	err := client.Execute(request, response)
	assert.NotNil(t, err)

	wasReconnectFailedError := false
	if tMultiError, ok := err.(*multierror.Error); ok {
		for _, me := range tMultiError.Errors {
			if me.Error() == errors.ErrorMessageReconnectFailed {
				wasReconnectFailedError = true
				break
			}
		}
	}
	assert.True(t, wasReconnectFailedError, "expected reconnect failed error")

}
