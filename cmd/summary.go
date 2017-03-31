package cmd

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(summaryCommand)
}

var summaryCommand = &cobra.Command{
        Use:   "summary",
        Short: "show environment information",
        RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Total network devices:\t10\n")
		fmt.Printf("Total managed devices:\t25\n")
		fmt.Printf("Total ofline devices:\t5\n")
		fmt.Printf("TBD\n")
		return nil
        },
}
