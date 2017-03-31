package cmd

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listNetworkCommand)
}

var listNetworkCommand = &cobra.Command{
        Use:   "list-network",
        Short: "list network",
        RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("TBD\n")
		return nil
        },
}
