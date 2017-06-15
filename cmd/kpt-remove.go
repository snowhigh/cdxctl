package cmd

import (
	"fmt"
        "log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var kptRemoveName string

func init() {
	kptRemoveCommand.Flags().StringVarP(&kptRemoveName, "name", "n", "", "package name")
	RootCmd.AddCommand(kptRemoveCommand)
}

var kptRemoveCommand = &cobra.Command{
        Use:   "kpt-remove",
        Short: "Kubernetes package tool, remove packages",
        RunE: func(cmd *cobra.Command, args []string) error {
		os.Chdir("/root/fullstack/cdxvirt")
		if kptRemoveName == "" {
			return cmd.Help()
		}
		if _, err := os.Stat(fmt.Sprintf("/root/fullstack/cdxvirt/%s/ansible", kptRemoveName)); err != nil {
			log.Fatal(err)
		}
		
		tmp_cmd := fmt.Sprintf("ansible-playbook /root/fullstack/cdxvirt/%s/ansible/uninstall.yml", kptRemoveName)
		runCommand(tmp_cmd, true)
		return nil
        },
}
