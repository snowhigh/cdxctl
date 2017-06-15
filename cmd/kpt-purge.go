package cmd

import (
	"fmt"
        "log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var kptPurgeName string

func init() {
	kptPurgeCommand.Flags().StringVarP(&kptPurgeName, "name", "n", "", "package name")
	RootCmd.AddCommand(kptPurgeCommand)
}

var kptPurgeCommand = &cobra.Command{
        Use:   "kpt-purge",
        Short: "Kubernetes package tool, purge packages",
        RunE: func(cmd *cobra.Command, args []string) error {
		os.Chdir("/root/fullstack/cdxvirt")
		if kptPurgeName == "" {
			return cmd.Help()
		}
		if _, err := os.Stat(fmt.Sprintf("/root/fullstack/cdxvirt/%s/ansible", kptPurgeName)); err != nil {
			log.Fatal(err)
		}
		
		tmp_cmd := fmt.Sprintf("ansible-playbook /root/fullstack/cdxvirt/%s/ansible/uninstall.yml", kptPurgeName)
		runCommand(tmp_cmd, true)
		return nil
        },
}
