package snapshotsdeleteschedule

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
	Use:   "delete-snapshot-schedule",
	Short: "Delete snapshot schedule",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
	opConfig      = configs.NewOpSnapshotDeleteScheduleConfig()
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

	logger := logConfig.NewLogger("delete-snapshot-schedule")

	for _, validatingConfig := range []configs.ValidatingConfig{commandConfig, opConfig} {
		if err := validatingConfig.Validate(); err != nil {
			logger.Error("configuration is invalid", "reason", err)
			return 1
		}
	}

	cliConfig, cliConfigErr := configs.NewYBClientConfigFromCliConfig(commandConfig)
	if cliConfigErr != nil {
		logger.Error("failed creating client configuration", "reason", cliConfigErr)
		return 1
	}
	cliClient, err := cli.NewYBConnectedClient(cliConfig, logger.Named("client"))
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

	responsePayload, err := cliClient.SnapshotsDeleteSchedule(opConfig)
	if err != nil {
		logger.Error("failed deleting snapshot schedule", "reason", err)
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
