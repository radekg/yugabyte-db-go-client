package main

import (
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client/cmd/getmasterregistration"
	"github.com/radekg/yugabyte-db-go-client/cmd/listmasters"
	"github.com/radekg/yugabyte-db-go-client/cmd/listtabletservers"
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
	rootCmd.AddCommand(getmasterregistration.Command)
	rootCmd.AddCommand(listmasters.Command)
	rootCmd.AddCommand(listtabletservers.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
