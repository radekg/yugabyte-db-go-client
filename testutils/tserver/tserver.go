package tserver

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/testutils/common"
	"github.com/radekg/yugabyte-db-go-client/testutils/master"

	"github.com/radekg/yugabyte-db-go-client/client"
	"github.com/radekg/yugabyte-db-go-client/configs"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// TestEnvContext represents a test YugabyteDB TServer environment context.
type TestEnvContext interface {
	// Use these addresses when connecting to the RPC system using your
	// own client instance.
	TServerExternalRPCPort() string
	// Use these addresses when connecting to the YSQL system using your
	// own client instance.
	TServerExternalYSQLPort() string
	// Always clean up after your run.
	Cleanup()
	// Other allocated ports for this TServer.
	OtherPorts() common.AllocatedAdditionalPorts
}

// SetupTServer sets up a single YugabyteDB TServer.
func SetupTServer(t *testing.T,
	mastersCtx master.TestEnvContext,
	config *common.TestTServerConfiguration) TestEnvContext {

	if mastersCtx.Pool() == nil {
		t.Fatalf("mastersCtx pool required")
	}

	if mastersCtx.Network() == nil {
		t.Fatalf("mastersCtx network required")
	}

	config = common.ApplyTServerConfigDefaults(t, config)

	// after every successful step, store a cleanup function here:
	closables := []func(){}

	// used in case of a failure during setup:
	closeClosables := func(closables []func()) {
		for _, closable := range closables {
			defer closable()
		}
	}

	// close all reasources in reverse order:
	prependClosable := func(closable func(), closables []func()) []func() {
		return append([]func(){closable}, closables...)
	}

	// RPC port:
	tserverRPCPort, err := common.NewRandomPortSupplier()
	if err != nil {
		closeClosables(closables)
		t.Fatalf("failed creating tserver RPC random port listener: '%v'", err)
	}
	closables = prependClosable(func() {
		t.Log("cleanup: closing tserver RPC random port listener, if not closed yet")
		tserverRPCPort.Cleanup()
	}, closables)
	if err := tserverRPCPort.Discover(); err != nil {
		closeClosables(closables)
		t.Fatalf("failed extracting host and port from tserver RPC random port listener: '%v'", err)
	}

	// YSQL port:
	tserverYSQLPort, err := common.NewRandomPortSupplier()
	if err != nil {
		closeClosables(closables)
		t.Fatalf("failed creating tserver YSQL random port listener: '%v'", err)
	}
	closables = prependClosable(func() {
		t.Log("cleanup: closing tserver YSQL random port listener, if not closed yet")
		tserverYSQLPort.Cleanup()
	}, closables)
	if err := tserverYSQLPort.Discover(); err != nil {
		closeClosables(closables)
		t.Fatalf("failed extracting host and port from tserver YSQL random port listener: '%v'", err)
	}

	fetchedTServerRPCPort, _ := tserverRPCPort.DiscoveredPort()
	fetchedTServerYSQLPort, _ := tserverYSQLPort.DiscoveredPort()

	tserverRPCBindAddress := fmt.Sprintf("0.0.0.0:%d", common.DefaultYugabyteDBTServerRPCPort)
	if config.TServerID != "" {
		tserverRPCBindAddress = fmt.Sprintf("%s:%d", config.TServerID, common.DefaultYugabyteDBTServerRPCPort)
	}

	// allocated additional requested ports:
	allocatedPorts := common.AllocatedAdditionalPorts{}
	for _, additionalPort := range config.AdditionalPorts {
		portSupplier, err := common.NewRandomPortSupplier()
		if err != nil {
			closeClosables(closables)
			t.Fatalf("failed creating random port listener for port '%s': '%v'", additionalPort, err)
		}
		closables = prependClosable(func() {
			t.Logf("cleanup: closing random port listener for port '%s', if not closed yet", additionalPort)
			portSupplier.Cleanup()
		}, closables)
		if err := portSupplier.Discover(); err != nil {
			closeClosables(closables)
			t.Fatalf("failed extracting host and port from random port listener for '%s': '%v'", additionalPort, err)
		}
		discoveredPort, _ := portSupplier.DiscoveredPort()
		allocatedPorts[additionalPort] = common.NewDefaultAllocatedAdditionalPort(additionalPort, discoveredPort, portSupplier)
	}

	// start RF number of containers and wait for them:

	benchStart := time.Now()
	chanTServerOK := make(chan struct{}, 1)
	chanTServerError := make(chan error, 1)

	// create temp data directory
	tserverDataDirectory, tempDirErr := ioutil.TempDir("", "yb-tserver")
	if tempDirErr != nil {
		closeClosables(closables)
		t.Fatalf("expected temp data directory to be created but received an error: '%v'", tempDirErr)
	}
	closables = prependClosable(func() {
		t.Log("cleanup: tserver data directory")
		os.Remove(tserverDataDirectory)
	}, closables)

	t.Log("TServer data directory created...", tserverDataDirectory)

	tserverCmd := []string{fmt.Sprintf("/home/%s/bin/yb-tserver", config.YbDBContainerUser),
		"--callhome_enabled=false",
		fmt.Sprintf("--fs_data_dirs=%s/tserver", config.YbDBFsDataPath),
		fmt.Sprintf("--tserver_master_addrs=%s", strings.Join(mastersCtx.MasterInternalAddresses(), ",")),
		fmt.Sprintf("--rpc_bind_addresses=%s", tserverRPCBindAddress),
		"--logtostderr",
		"--enable_ysql",
		"--ysql_enable_auth",
		"--minloglevel=1",
		"--placement_cloud=dockertest",
		"--placement_region=test1",
		"--placement_zone=test1a",
		"--stop_on_parent_termination",
		"--undefok=stop_on_parent_termination",
	}

	if config.YbDBCmdSupplier != nil {
		tserverCmd = config.YbDBCmdSupplier(mastersCtx.MasterInternalAddresses(), tserverRPCBindAddress)
	}

	portBindings := map[dc.Port][]dc.PortBinding{
		dc.Port(fmt.Sprintf("%d/tcp", common.DefaultYugabyteDBTServerRPCPort)):  {{HostIP: "0.0.0.0", HostPort: fetchedTServerRPCPort}},
		dc.Port(fmt.Sprintf("%d/tcp", common.DefaultYugabyteDBTServerYSQLPort)): {{HostIP: "0.0.0.0", HostPort: fetchedTServerYSQLPort}},
	}
	for _, allocatedPort := range allocatedPorts {
		portBindings[allocatedPort.Requested()] = []dc.PortBinding{
			{HostIP: "0.0.0.0", HostPort: allocatedPort.Allocated()},
		}
	}

	options := &dockertest.RunOptions{
		Name:       config.TServerID,
		Repository: common.GetEnvOrDefault(common.DefaultYugabyteDBEnvVarImageName, config.YbDBDockerImage),
		Tag:        common.GetEnvOrDefault(common.DefaultYugabyteDEnvVarImageVersion, config.YbDBDockerTag),
		Mounts: []string{
			fmt.Sprintf("%s:%s/tserver", tserverDataDirectory, config.YbDBFsDataPath),
		},
		Cmd: tserverCmd,
		ExposedPorts: []string{
			fmt.Sprintf("%d/tcp", common.DefaultYugabyteDBTServerRPCPort),
			fmt.Sprintf("%d/tcp", common.DefaultYugabyteDBTServerYSQLPort)},
		PortBindings: portBindings,
		Env:          config.YbDBEnv,
		Networks:     []*dockertest.Network{mastersCtx.Network()},
	}

	tserverRPCPort.Cleanup()
	tserverYSQLPort.Cleanup()
	for _, v := range allocatedPorts {
		v.Use()
	}

	// start the container:
	tserver, tserverErr := mastersCtx.Pool().RunWithOptions(options, func(config *dc.HostConfig) {
		config.RestartPolicy.Name = "on-failure"
	})
	if tserverErr != nil {
		closeClosables(closables)
		t.Fatalf("expected tserver to start but received: '%v'", tserverErr)
	}
	closables = prependClosable(func() {
		t.Log("cleanup: closing tserver")
		if !config.NoCleanupContainers {
			tserver.Close()
			mastersCtx.Pool().Purge(tserver)
		}
	}, closables)

	t.Logf("TServer started with container ID '%s', waiting for the server to reply...", tserver.Container.ID)

	go func() {
		poolRetryErr := mastersCtx.Pool().Retry(func() error {

			t.Log("Querying TServer status", tserver.Container.ID)

			client, err := client.Connect(&configs.YBSingleNodeClientConfig{
				MasterHostPort: fmt.Sprintf("127.0.0.1:%s", fetchedTServerRPCPort),
				OpTimeout:      uint32(time.Duration(time.Second * 60).Seconds()),
			}, hclog.Default())

			if err != nil {
				if config.LogRegistrationRetryErrors {
					// this cannot be t.Error because that will fail the test!!!
					t.Log("TServer client could not connect directly, reported an error:", err)
				}
				return err
			}

			select {
			case err := <-client.OnConnectError():
				if config.LogRegistrationRetryErrors {
					// this cannot be t.Error because that will fail the test!!!
					t.Log("TServer OnConnectError reported an error:", err)
				}
				return err
			case <-client.OnConnected():
			}

			defer client.Close()

			request := &ybApi.IsTabletServerReadyRequestPB{}
			response := &ybApi.IsTabletServerReadyResponsePB{}
			isTabletServerReadyError := client.Execute(request, response)

			if isTabletServerReadyError != nil {
				if config.LogRegistrationRetryErrors {
					// this cannot be t.Error because that will fail the test!!!
					t.Log("TServer IsTabletServerReady reported an error:", isTabletServerReadyError)
				}
				return isTabletServerReadyError
			}

			isTabletServerReadyError = errors.NewTabletServerError(response.Error)
			if isTabletServerReadyError != nil {
				if config.LogRegistrationRetryErrors {
					// this cannot be t.Error because that will fail the test!!!
					t.Log("TServer IsTabletServerReady reported an error:", isTabletServerReadyError)
				}
				return isTabletServerReadyError
			}

			t.Log("TServer reported its status:", response)

			return nil

		})
		if poolRetryErr == nil {
			close(chanTServerOK)
			return
		}
		chanTServerError <- poolRetryErr
	}()

	select {
	case <-chanTServerOK:
		t.Logf("TServer replied after: %s", time.Now().Sub(benchStart).String())
	case receivedError := <-chanTServerError:
		closeClosables(closables)
		t.Fatalf("TServer wait finished with error: '%v'", receivedError)
	case <-time.After(time.Second * 45):
		closeClosables(closables)
		t.Fatalf("TServer did not start communicating within timeout")
	}

	return &testEnvContext{
		tServerExternalRPCPortValue:  fetchedTServerRPCPort,
		tServerExternalYSQLPortValue: fetchedTServerYSQLPort,
		cleanupFuncValue: func() {
			for _, closable := range closables {
				closable()
			}
		},
		otherPortsValue: allocatedPorts,
	}

}

type testEnvContext struct {
	tServerExternalRPCPortValue  string
	tServerExternalYSQLPortValue string
	cleanupFuncValue             func()
	otherPortsValue              common.AllocatedAdditionalPorts
}

func (ctx *testEnvContext) TServerExternalRPCPort() string {
	return ctx.tServerExternalRPCPortValue
}

func (ctx *testEnvContext) TServerExternalYSQLPort() string {
	return ctx.tServerExternalYSQLPortValue
}

func (ctx *testEnvContext) Cleanup() {
	ctx.cleanupFuncValue()
}

func (ctx *testEnvContext) OtherPorts() common.AllocatedAdditionalPorts {
	return ctx.otherPortsValue
}
