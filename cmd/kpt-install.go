package cmd

import (
	"fmt"
        "log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var kptInstallName string
var kptInstallInteractive bool
var kptInstallSkipUpload bool

func init() {
	kptInstallCommand.Flags().StringVarP(&kptInstallName, "name", "n", "", "package name")
	kptInstallCommand.Flags().BoolVarP(&kptInstallInteractive, "interactive", "t", true, "interactive mode")
	kptInstallCommand.Flags().BoolVarP(&kptInstallSkipUpload, "skip", "s", false, "skip upload.yml")
	RootCmd.AddCommand(kptInstallCommand)
}

var kptInstallCommand = &cobra.Command{
        Use:   "kpt-install",
        Short: "Kubernetes package tool, install packages",
        RunE: func(cmd *cobra.Command, args []string) error {
		var tmp_cmd string
		os.Chdir("/root/fullstack/cdxvirt")
		if kptInstallName == "" {
			return cmd.Help()
		}
		if _, err := os.Stat(fmt.Sprintf("/root/fullstack/cdxvirt/%s/ansible", kptInstallName)); err != nil {
			log.Fatal(err)
		}
		
		if kptInstallSkipUpload == false {
			tmp_cmd = fmt.Sprintf("ansible-playbook /root/fullstack/cdxvirt/%s/ansible/upload.yml", kptInstallName)
			runCommand(tmp_cmd, kptInstallInteractive)
		}
		tmp_cmd = fmt.Sprintf("ansible-playbook /root/fullstack/cdxvirt/%s/ansible/install.yml", kptInstallName)
		runCommand(tmp_cmd, kptInstallInteractive)
		return nil
        },
}
