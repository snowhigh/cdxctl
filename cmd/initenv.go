package cmd

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initEnvCommand)
}

var initEnvCommand = &cobra.Command{
        Use:   "initenv",
        Short: "Initial environment, dhcp, ipxe, ... servers",
        RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("TBD\n")
		return nil
        },
}
