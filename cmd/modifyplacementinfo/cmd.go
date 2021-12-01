package modifyplacementinfo

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
	Use:   "modify-placement-info",
	Short: "Modifies the placement information (cloud, region, and zone) for a deployment",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
	opConfig      = configs.NewOpModifyPlacementInfoConfig()
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

	logger := logConfig.NewLogger("modify-placement-info")

	for _, validatingConfig := range []configs.ValidatingConfig{commandConfig, opConfig} {
		if err := validatingConfig.Validate(); err != nil {
			logger.Error("configuration is invalid", "reason", err)
			return 1
		}
	}

	cliClient, err := implementation.MasterLeaderConnectedClient(commandConfig, logger.Named("client"))
	if err != nil {
		logger.Error("could not connect to a leader master", "reason", err)
		return 1
	}
	defer cliClient.Close()

	responsePayload, err := cliClient.ModifyPlacementInfo(opConfig)
	if err != nil {
		logger.Error("failed modifying placement info", "reason", err)
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
