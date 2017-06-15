package cmd

import (
	"fmt"
        "log"
	"os"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(kptListCommand)
}

var kptListCommand = &cobra.Command{
        Use:   "kpt-list",
        Short: "Kubernetes package tool, list all packages",
        RunE: func(cmd *cobra.Command, args []string) error {
		os.Chdir("/root/fullstack/cdxvirt")

		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				fmt.Println(file.Name())
			}
		}
		return nil
        },
}
