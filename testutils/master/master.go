package master

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
	"github.com/radekg/yugabyte-db-go-client/testutils/common"

	"github.com/radekg/yugabyte-db-go-client/client/implementation"
	"github.com/radekg/yugabyte-db-go-client/configs"
)

// TestEnvContext represents a test YugabyteDB master environment context.
type TestEnvContext interface {
	// Use these addresses when connecting to the RPC system using your
	// own client instance.
	MasterExternalAddresses() []string
	// use these addresses when adding TServers to your test cluster.
	MasterInternalAddresses() []string
	// Always clean up after your run.
	Cleanup()

	// Gives access to the underlying docker pool used by this context.
	Pool() *dockertest.Pool
	// Network contains the docker network in which masters are running.
	// Your TServers must be started in the same network.
	Network() *dockertest.Network
}

// SetupMasters sets up RF number of YugabyteDB masters.
func SetupMasters(t *testing.T, config *common.TestMasterConfiguration) TestEnvContext {

	config = common.ApplyMasterConfigDefaults(t, config)

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

	masterRPCPorts := []common.RandomPortSupplier{}
	fetchedRPCPorts := []string{}

	runID := time.Now().Unix()

	masterExternalAddresses := []string{}
	masterInternalAddresses := []string{}
	containerNames := []string{}

	// get RF number of RPC ports:
	for i := 0; i < config.ReplicationFactor; i = i + 1 {

		containerName := fmt.Sprintf("%s-%d", config.MasterPrefix, i)

		// master RPC port:
		masterRPCPort, err := common.NewRandomPortSupplier()
		if err != nil {
			closeClosables(closables)
			t.Fatalf("failed creating master RPC random port listener: '%v'", err)
		}
		closables = prependClosable(func() {
			t.Log("cleanup: closing master RPC random port listener, if not closed yet")
			masterRPCPort.Cleanup()
		}, closables)
		if err := masterRPCPort.Discover(); err != nil {
			closeClosables(closables)
			t.Fatalf("failed extracting host and port from master RPC random port listener: '%v'", err)
		}
		fetchedMasterRPCPort, _ := masterRPCPort.DiscoveredPort()
		masterRPCPorts = append(masterRPCPorts, masterRPCPort)
		fetchedRPCPorts = append(fetchedRPCPorts, fetchedMasterRPCPort)

		masterExternalAddresses = append(masterExternalAddresses, fmt.Sprintf("127.0.0.1:%s", fetchedMasterRPCPort))
		masterInternalAddresses = append(masterInternalAddresses, fmt.Sprintf("%s:%d", containerName, common.DefaultYugabyteDBMasterRPCPort))

		containerNames = append(containerNames, containerName)
	}

	// create new pool using the default Docker endpoint:
	pool, poolErr := dockertest.NewPool("")
	if poolErr != nil {
		closeClosables(closables)
		t.Fatalf("expected docker pool to come up but received: '%v'", poolErr)
	}

	networkName := fmt.Sprintf("yb-net-%d", runID)

	dockerNetwork, dockerNetworkErr := pool.CreateNetwork(networkName)
	if dockerNetworkErr != nil {
		closeClosables(closables)
		t.Fatalf("expected docker network to be created but received: '%v'", dockerNetworkErr)
	}
	closables = prependClosable(func() {
		t.Log("cleanup: removing docker network")
		pool.RemoveNetwork(dockerNetwork)
	}, closables)

	// start RF number of containers and wait for them:

	benchStart := time.Now()
	chanMastersOK := make(chan string, config.ReplicationFactor)
	chanMastersError := make(chan error, config.ReplicationFactor)

	for i := 0; i < config.ReplicationFactor; i = i + 1 {

		// create temp data directory
		masterDataDirectory, tempDirErr := ioutil.TempDir("", fmt.Sprintf("yb-master-%d", i))
		if tempDirErr != nil {
			closeClosables(closables)
			t.Fatalf("expected temp data directory to be created but received an error: '%v'", tempDirErr)
		}
		closables = prependClosable(func() {
			t.Log("cleanup: master data directory")
			os.Remove(masterDataDirectory)
		}, closables)

		t.Log("Master data directory created...", masterDataDirectory)

		options := &dockertest.RunOptions{
			Name:       containerNames[i],
			Repository: config.YbDBDockerImage,
			Tag:        config.YbDBDockerTag,
			Mounts: []string{
				fmt.Sprintf("%s:%s/master", masterDataDirectory, config.YbDBFsDataPath),
			},
			Cmd: []string{fmt.Sprintf("/home/%s/bin/yb-master", config.YbDBContainerUser),
				"--callhome_enabled=false",
				fmt.Sprintf("--fs_data_dirs=%s/master", config.YbDBFsDataPath),
				fmt.Sprintf("--master_addresses=%s", strings.Join(masterInternalAddresses, ",")),
				fmt.Sprintf("--rpc_bind_addresses=%s", masterInternalAddresses[i]),
				"--logtostderr",
				"--minloglevel=1",
				"--placement_cloud=dockertest",
				"--stop_on_parent_termination",
				"--undefok=stop_on_parent_termination",
				"--replication_factor=1",
			},
			ExposedPorts: []string{
				fmt.Sprintf("%d/tcp", common.DefaultYugabyteDBMasterRPCPort)},
			PortBindings: map[dc.Port][]dc.PortBinding{
				dc.Port(fmt.Sprintf("%d/tcp", common.DefaultYugabyteDBMasterRPCPort)): {{HostIP: "0.0.0.0", HostPort: fetchedRPCPorts[i]}},
			},
			Env:      config.YbDBEnv,
			Networks: []*dockertest.Network{dockerNetwork},
		}

		masterRPCPorts[i].Cleanup()

		go func(masterIndex int) {

			// start the container:
			masterServer, masterErr := pool.RunWithOptions(options, func(config *dc.HostConfig) {
				config.RestartPolicy.Name = "on-failure"
			})
			if masterErr != nil {
				closeClosables(closables)
				chanMastersError <- fmt.Errorf("expected master to start but received: '%v'", masterErr)
				return
			}
			closables = prependClosable(func() {
				t.Log("cleanup: closing master")
				if !config.NoCleanupContainers {
					masterServer.Close()
					pool.Purge(masterServer)
				}
			}, closables)

			t.Logf("Master started with container ID '%s', waiting for the server to reply...", masterServer.Container.ID)

			poolRetryErr := pool.Retry(func() error {

				t.Log("Querying master registration:", masterInternalAddresses[masterIndex])

				client, err := implementation.MasterLeaderConnectedClient(&configs.CliConfig{
					MasterHostPort: masterExternalAddresses,
					OpTimeout:      time.Duration(time.Second * 5),
				}, hclog.Default())
				if err != nil {
					if config.LogRegistrationRetryErrors {
						// this cannot be t.Error because that will fail the test!!!
						t.Log("Master", masterInternalAddresses[masterIndex], "connect reported an error:", err)
					}
					return err
				}

				defer client.Close()

				registrationPb, err := client.GetMasterRegistration()
				if err != nil {
					if config.LogRegistrationRetryErrors {
						// this cannot be t.Error because that will fail the test!!!
						t.Log("Master", masterInternalAddresses[masterIndex], "reported an error:", err)
					}
					return err
				}

				t.Log("Master", masterInternalAddresses[masterIndex], "reported its registration:", registrationPb)

				return nil
			})
			if poolRetryErr == nil {
				chanMastersOK <- masterInternalAddresses[masterIndex]
				return
			}
			chanMastersError <- poolRetryErr

		}(i)

	}

	reported := 0

outLoop:
	for {
		select {
		case masterAddress := <-chanMastersOK:
			t.Logf("Master '%s' replied after: %s", masterAddress, time.Now().Sub(benchStart).String())
			reported = reported + 1
			if reported == config.ReplicationFactor {
				break outLoop
			}
		case receivedError := <-chanMastersError:
			closeClosables(closables)
			t.Fatalf("Master wait finished with error: '%v'", receivedError)
			reported = reported + 1
			if reported == config.ReplicationFactor {
				break outLoop
			}
		case <-time.After(time.Second * 45):
			closeClosables(closables)
			t.Fatalf("Masters did not start communicating within timeout")
		}
	}

	return &testEnvContext{
		masterExternalAddressesValue: masterExternalAddresses,
		masterInternalAddressesValue: masterInternalAddresses,
		cleanupFuncValue: func() {
			for _, closable := range closables {
				closable()
			}
		},
		poolValue:    pool,
		networkValue: dockerNetwork,
	}

}

type testEnvContext struct {
	masterExternalAddressesValue []string
	masterInternalAddressesValue []string
	cleanupFuncValue             func()
	poolValue                    *dockertest.Pool
	networkValue                 *dockertest.Network
}

func (ctx *testEnvContext) MasterExternalAddresses() []string {
	return ctx.masterExternalAddressesValue
}

func (ctx *testEnvContext) MasterInternalAddresses() []string {
	return ctx.masterInternalAddressesValue
}

func (ctx *testEnvContext) Cleanup() {
	ctx.cleanupFuncValue()
}

func (ctx *testEnvContext) Pool() *dockertest.Pool {
	return ctx.poolValue
}

func (ctx *testEnvContext) Network() *dockertest.Network {
	return ctx.networkValue
}
