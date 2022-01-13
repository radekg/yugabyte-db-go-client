package master

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
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
	})

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

func TestMasterLeaderStepdown(t *testing.T) {

	debugLog := hclog.Default()
	debugLog.SetLevel(hclog.Debug)

	metricsCallback := common.NewTestMetricsCallback(t)

	request := &ybApi.ListMastersRequestPB{}

	testCtx := SetupMasters(t, &common.TestMasterConfiguration{
		ReplicationFactor: 1,
		MasterPrefix:      "master-it",
	})
	defer testCtx.Cleanup()

	client := client.NewYBClient(&configs.YBClientConfig{
		MasterHostPort:       testCtx.MasterExternalAddresses(),
		OpTimeout:            time.Duration(time.Second * 5),
		MaxReconnectAttempts: 20,
	}).WithLogger(debugLog).WithMetricsCallback(metricsCallback)

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

	stepdownRequest := &ybApi.LeaderStepDownRequestPB{
		TabletId: []byte("00000000000000000000000000000000"),
	}
	stepdownResponse := &ybApi.LeaderStepDownResponsePB{}
	stepdownErr := client.Execute(stepdownRequest, stepdownResponse)
	assert.Nil(t, stepdownErr)

	t.Log("Received stepdown response", stepdownResponse)

	tsRequest := &ybApi.ListTabletServersRequestPB{}
	tsResponse := &ybApi.ListTabletServersResponsePB{}
	err := client.Execute(tsRequest, tsResponse)
	assert.Nil(t, err)
	t.Log("Received tablet servers list", tsResponse)

	assert.Greater(t, metricsCallback.InspectClientReconnectAttempt(t), 0, "expected some reconnect attempts")
	assert.Equal(t, 1, metricsCallback.InspectClientReconnectSuccess(t), "expected one successful reconnect")

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
	})

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

	wasReconnectFailedError := strings.HasPrefix(err.Error(), errors.ErrorMessageReconnectFailed)
	assert.True(t, wasReconnectFailedError, "expected reconnect failed error")

}
