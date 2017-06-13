package cmd

import (
	"fmt"
        "log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var clusterConnectIp string

func init() {
	clusterConnectCommand.Flags().StringVarP(&clusterConnectIp, "ip", "i", "", "Node IP address")
	RootCmd.AddCommand(clusterConnectCommand)
}

var clusterConnectCommand = &cobra.Command{
        Use:   "cluster-connect",
        Short: "Setup kubectl config to selected cluster",
        RunE: func(cmd *cobra.Command, args []string) error {
		if clusterConnectIp == "" {
			return cmd.Help()
		}
		nodeIP := clusterConnectIp
		clusterConnect(nodeIP)
		fmt.Println("Done. Try \"kubectl get nodes\"")
		return nil
        },
}

func clusterConnect(nodeIP string) {
	os.Chdir("/root/provision/")
	os.Link("/root/provision/playbooks/files/cluster/kubectl", "/usr/local/bin/kubectl")

	tmp_cmd := fmt.Sprintf("kubectl config set-cluster default --server=https://%s:6443 --insecure-skip-tls-verify=true", nodeIP)
	runCommand(tmp_cmd)

	// kubectl config set-credentials admin --username=admin --password=admin
	runCommand("kubectl config set-credentials admin --username=admin --password=admin")

	// kubectl config set-context default --cluster=default --user=admin
	runCommand("kubectl config set-context default --cluster=default --user=admin")

	// kubectl config use-context default
	runCommand("kubectl config use-context default")
}

func runCommand(tmp_cmd string) {
	split_cmd := strings.Split(tmp_cmd, " ")
        cmd := exec.Command(split_cmd[0], split_cmd[1:]...)
        cmd.Stderr = os.Stderr
        cmd.Stdout = os.Stdout
        err := cmd.Run()
        if err != nil {
                log.Fatal(err)
        }
}
