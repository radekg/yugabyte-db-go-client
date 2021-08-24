package main

import (
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client/cmd/checkexists"
	"github.com/radekg/yugabyte-db-go-client/cmd/getloadmovecompletion"
	"github.com/radekg/yugabyte-db-go-client/cmd/getmasterregistration"
	"github.com/radekg/yugabyte-db-go-client/cmd/getuniverseconfig"
	"github.com/radekg/yugabyte-db-go-client/cmd/isloadbalanced"
	"github.com/radekg/yugabyte-db-go-client/cmd/isserverready"
	"github.com/radekg/yugabyte-db-go-client/cmd/listmasters"
	"github.com/radekg/yugabyte-db-go-client/cmd/listtabletservers"
	"github.com/radekg/yugabyte-db-go-client/cmd/ping"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ybcli",
	Short: "ybcli",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(checkexists.Command)
	rootCmd.AddCommand(getloadmovecompletion.Command)
	rootCmd.AddCommand(getmasterregistration.Command)
	rootCmd.AddCommand(getuniverseconfig.Command)
	rootCmd.AddCommand(isloadbalanced.Command)
	rootCmd.AddCommand(isserverready.Command)
	rootCmd.AddCommand(listmasters.Command)
	rootCmd.AddCommand(listtabletservers.Command)
	rootCmd.AddCommand(ping.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
