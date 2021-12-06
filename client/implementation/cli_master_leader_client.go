package implementation

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/configs"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// MasterLeaderConnectedClient creates a client for every configured master host port and tests each one for
// the leader status. First client for whom the relevant master server returns Raft peer role: Leader
// if returned to the caller.
// If no master leader can be found or all clients experienced a connection error, an error is returned
// to the caller.
func MasterLeaderConnectedClient(commandConfig *configs.CliConfig, logger hclog.Logger) (YBCliClient, error) {

	validConfigs := map[string]*configs.YBClientConfig{}
	for _, hostPort := range commandConfig.MasterHostPort {
		cliConfig, cliConfigErr := configs.NewYBClientConfigFromCliConfig(hostPort, commandConfig)
		if cliConfigErr != nil {
			logger.Error("failed creating client configuration, skipping",
				"reason", cliConfigErr,
				"host-port", hostPort)
			continue
		}
		validConfigs[hostPort] = cliConfig
	}

	chanConnectedClient := make(chan YBCliClient, 1)
	chanErrors := make(chan error)

	for hostPort, cliConfig := range validConfigs {
		go func(thisHostPort string, thisConfig *configs.YBClientConfig) {
			cliClient, err := NewYBConnectedClient(thisConfig, logger.Named("client"))
			if err != nil {
				logger.Error("failed creating a client",
					"reason", err,
					"host-port", thisHostPort)
				chanErrors <- err
				return
			}

			select {
			case err := <-cliClient.OnConnectError():

				logger.Error("connection error",
					"reason", err,
					"host-port", thisHostPort)
				cliClient.Close()
				chanErrors <- err

			case <-cliClient.OnConnected():

				masterRegistration, err := cliClient.GetMasterRegistration()

				if err != nil {
					logger.Error("failed querying master registration",
						"reason", err,
						"host-port", thisHostPort)
					cliClient.Close()
					chanErrors <- err
					return
				}

				if masterRegistration == nil {
					logger.Trace("master did not send with registration info",
						"host-port", thisHostPort)
					cliClient.Close()
					chanErrors <- fmt.Errorf("master %s did not send registration info", thisHostPort)
					return
				}

				if masterRegistration.Role == nil {
					logger.Trace("master did not report its raft peer role",
						"host-port", thisHostPort)
					cliClient.Close()
					chanErrors <- fmt.Errorf("master %s did not report its raft peer role", thisHostPort)
					return
				}

				if *masterRegistration.Role != ybApi.RaftPeerPB_LEADER {
					logger.Trace("master not leader",
						"host-port", thisHostPort)
					cliClient.Close()
					chanErrors <- fmt.Errorf("master %s not leader", thisHostPort)
					return
				}

				logger.Info("found master leader", "host-port", thisHostPort)
				chanConnectedClient <- cliClient
			}

		}(hostPort, cliConfig)
	}

	var done uint64
	max := uint64(len(validConfigs))

	for {
		select {
		case <-chanErrors:
			atomic.AddUint64(&done, 1)
			if atomic.LoadUint64(&done) == max {
				return nil, fmt.Errorf("no reachable leaders")
			}
		case connectedClient := <-chanConnectedClient:
			return connectedClient, nil
		case <-time.After(commandConfig.OpTimeout):
			return nil, fmt.Errorf("failed to connect to a leader master within timeout")
		}
	}

}
