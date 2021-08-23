package listtabletservers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client/client"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/spf13/cobra"
)

// Command is the command declaration.
var Command = &cobra.Command{
	Use:   "list-tablet-servers",
	Short: "List all the tablet servers in this database",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
	opConfig      = configs.NewOpListTableServersConfig()
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

	connectedClient, err := client.Connect(configs.NewYBClientConfigFromCliConfig(commandConfig), logger.Named("client"))
	if err != nil {
		logger.Error("failed creating a client", "reason", err)
		return 1
	}
	select {
	case err := <-connectedClient.OnConnectError():
		logger.Error("failed connecting a client", "reason", err)
		return 1
	case <-connectedClient.OnConnected():
		logger.Debug("client connected")
	}
	defer connectedClient.Close()

	responsePayload, err := connectedClient.ListTabletServers(opConfig)
	if err != nil {
		logger.Error("failed reading tablet servers list", "reason", err)
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
