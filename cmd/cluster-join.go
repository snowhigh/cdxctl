package cmd

import (
	"fmt"
        "log"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var clusterJoinIp string
var clusterJoinName string
var clusterJoinNetwork string

func init() {
	clusterJoinCommand.Flags().StringVarP(&clusterJoinIp, "ip", "i", "", "Node IP address")
	clusterJoinCommand.Flags().StringVarP(&clusterJoinName, "name", "n", "", "Cluster name")
	clusterJoinCommand.Flags().StringVarP(&clusterJoinNetwork, "net", "w", "", "Cluster network ex. 192.168.32.0/23")
	RootCmd.AddCommand(clusterJoinCommand)
}

var clusterJoinCommand = &cobra.Command{
        Use:   "cluster-join",
        Short: "add device into cluster",
        RunE: func(cmd *cobra.Command, args []string) error {
		if clusterJoinIp == "" || clusterJoinName == "" {
			return cmd.Help()
		}
		clusterID := clusterJoinName
		nodeIP := clusterJoinIp
		network := clusterJoinNetwork
		clusterJoin(clusterID, nodeIP, network)
		log.Printf("Done")
		return nil
        },
}

func clusterJoin(clusterID string, nodeIP string, network string) {
	os.Chdir("/root/provision")
	// HOST_IP_LIST="$HOSTS" bash upload-preinit-scripts.sh
	cmd := exec.Command("/bin/bash", "upload-preinit-scripts.sh")
	env := os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash pull-all-img.sh
	cmd = exec.Command("/bin/bash", "pull-all-img.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash -x k8sup.sh
	cmd = exec.Command("/bin/bash", "k8sup.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	env = append(env, fmt.Sprintf("CLUSTER_ID=%s", clusterID))
	if ( network != "" ) {
		env = append(env, fmt.Sprintf("NETWORK=%s", network))
	}
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
