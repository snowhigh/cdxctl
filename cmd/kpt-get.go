package cmd

import (
	"fmt"
        "log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var kptGetName string
var kptGetInteractive bool
var kptGetSkipUpload bool

func init() {
	kptGetCommand.Flags().StringVarP(&kptGetName, "name", "n", "", "package name")
	kptGetCommand.Flags().BoolVarP(&kptGetInteractive, "interactive", "t", true, "interactive mode")
	kptGetCommand.Flags().BoolVarP(&kptGetSkipUpload, "skip", "s", false, "skip upload.yml")
	RootCmd.AddCommand(kptGetCommand)
}

var kptGetCommand = &cobra.Command{
        Use:   "kpt-get",
        Short: "Kubernetes package tool, install packages",
        RunE: func(cmd *cobra.Command, args []string) error {
		var tmp_cmd string
		os.Chdir("/root/fullstack/cdxvirt")
		if kptGetName == "" {
			return cmd.Help()
		}
		if _, err := os.Stat(fmt.Sprintf("/root/fullstack/cdxvirt/%s/ansible", kptGetName)); err != nil {
			log.Fatal(err)
		}
		
		if kptGetSkipUpload == false {
			tmp_cmd = fmt.Sprintf("ansible-playbook /root/fullstack/cdxvirt/%s/ansible/upload.yml", kptGetName)
			runCommand(tmp_cmd, kptGetInteractive)
		}
		tmp_cmd = fmt.Sprintf("ansible-playbook /root/fullstack/cdxvirt/%s/ansible/install.yml", kptGetName)
		runCommand(tmp_cmd, kptGetInteractive)
		return nil
        },
}
