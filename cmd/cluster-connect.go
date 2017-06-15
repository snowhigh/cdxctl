package cmd

import (
	"fmt"
	"os"

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
	runCommand(tmp_cmd, true)

	// kubectl config set-credentials admin --username=admin --password=admin
	runCommand("kubectl config set-credentials admin --username=admin --password=admin", true)

	// kubectl config set-context default --cluster=default --user=admin
	runCommand("kubectl config set-context default --cluster=default --user=admin", true)

	// kubectl config use-context default
	runCommand("kubectl config use-context default", true)

	// Create /etc/ansible/hosts file
	f, _ := os.Create("/etc/ansible/hosts")
	defer f.Close()
	f.WriteString("[all]\n")
	f.WriteString(fmt.Sprintf("%s\n", nodeIP))
	f.WriteString("[all:vars]\n")
	f.WriteString("ansible_ssh_user=root\n")
	f.WriteString("ansible_python_interpreter=\"/usr/bin/python\"\n")
	f.WriteString("ansible_port=2222\n")
	f.WriteString("ansible_ssh_pass=!Q@W3e4r\n")
	f.Sync()
}
