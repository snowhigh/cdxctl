package cmd

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var VERSION string = "1.0.0"

func init() {
	RootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
        Use:   "version",
        Short: "show version information",
        RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("version %s\n", VERSION)
		return nil
        },
}
