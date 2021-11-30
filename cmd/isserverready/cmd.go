package isserverready

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client/client/implementation"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/spf13/cobra"
)

// Command is the command declaration.
var Command = &cobra.Command{
	Use:   "is-server-ready",
	Short: "Check if server is ready to serve IO requests",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
	opConfig      = configs.NewOpIsServerReadyConfig()
)

func initFlags() {
	Command.Flags().AddFlagSet(commandConfig.FlagSet())
	Command.Flags().AddFlagSet(logConfig.FlagSet())
	Command.Flags().AddFlagSet(opConfig.FlagSet())
}

func init() {
	initFlags()
}

func run(cobraCommand *cobra.Command, _ []string) {
	os.Exit(processCommand())
}

func processCommand() int {

	logger := logConfig.NewLogger("is-server-ready")

	for _, validatingConfig := range []configs.ValidatingConfig{commandConfig, opConfig} {
		if err := validatingConfig.Validate(); err != nil {
			logger.Error("configuration is invalid", "reason", err)
			return 1
		}
	}

	cliConfig, cliConfigErr := configs.NewYBClientConfigFromCliConfig(commandConfig.MasterHostPort[0], commandConfig)
	if cliConfigErr != nil {
		logger.Error("failed creating client configuration", "reason", cliConfigErr)
		return 1
	}
	cliConfig.MasterHostPort = fmt.Sprintf("%s:%d", opConfig.Host, opConfig.Port)
	// TODO: REVISIT: what's the is-tserver flag for?
	cliClient, err := implementation.NewYBConnectedClient(cliConfig, logger.Named("client"))
	if err != nil {
		// careful: different than other commands:
		logger.Error("server not reachable", "reason", err)
		return 2
	}
	select {
	case err := <-cliClient.OnConnectError():
		// TODO: LATER: in this case, this may indicate the service unavailability
		logger.Error("failed connecting a client", "reason", err)
		return 1
	case <-cliClient.OnConnected():
		logger.Debug("client connected")
	}
	defer cliClient.Close()

	responsePayload, err := cliClient.IsTabletServerReady()
	if err != nil {
		// TODO: LATER: in this case, this may indicate the service unavailability
		logger.Error("failed reading server ready response", "reason", err)
		return 1
	}

	jsonBytes, err := json.MarshalIndent(responsePayload, "", "  ")
	if err != nil {
		logger.Error("failed marshaling JSON response", "reason", err)
		return 1
	}

	fmt.Println(string(jsonBytes))

	return 0
}
