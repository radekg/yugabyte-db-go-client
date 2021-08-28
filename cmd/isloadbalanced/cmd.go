package isloadbalanced

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client/client/cli"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/spf13/cobra"
)

// Command is the command declaration.
var Command = &cobra.Command{
	Use:   "is-load-balanced",
	Short: "Check if master leader thinks that the load is balanced across tservers",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
	opConfig      = configs.NewOpIsLoadBalancedConfig()
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

	logger := logConfig.NewLogger("get_master_registration")

	for _, validatingConfig := range []configs.ValidatingConfig{commandConfig} {
		if err := validatingConfig.Validate(); err != nil {
			logger.Error("configuration is invalid", "reason", err)
			return 1
		}
	}

	cliClient, err := cli.NewYBConnectedClient(configs.NewYBClientConfigFromCliConfig(commandConfig), logger.Named("client"))
	if err != nil {
		logger.Error("failed creating a client", "reason", err)
		return 1
	}
	select {
	case err := <-cliClient.OnConnectError():
		logger.Error("failed connecting a client", "reason", err)
		return 1
	case <-cliClient.OnConnected():
		logger.Debug("client connected")
	}
	defer cliClient.Close()

	responsePayload, err := cliClient.IsLoadBalanced(opConfig)
	if err != nil {
		logger.Error("failed reading load balanced state", "reason", err)
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
