package cmd

import (
	"fmt"
        "log"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var nodeInitIp string
var nodeInitVerbose bool

func init() {
	nodeInitCommand.Flags().StringVarP(&nodeInitIp, "ip", "i", "", "Node IP address")
	nodeInitCommand.Flags().BoolVarP(&nodeInitVerbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(nodeInitCommand)
}

var nodeInitCommand = &cobra.Command{
        Use:   "node-init",
        Short: "initialize node information",
        RunE: func(cmd *cobra.Command, args []string) error {
		if nodeInitIp == "" {
			return cmd.Help()
		}
		nodeIP := nodeInitIp
		nodeInit(nodeIP)
		log.Printf("Done")
		return nil
        },
}

func nodeInit(nodeIP string) {
	os.Chdir("/root/provision/")

	// HOST_IP_LIST="$HOSTS" bash qts-qes-switcher.sh
	cmd := exec.Command("/bin/bash", "qts-qes-switcher.sh")
	env := os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	if nodeInitVerbose {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" OPTS="--install --cd --reboot" bash dom-modifier.sh
	cmd = exec.Command("/bin/bash", "dom-modifier.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	env = append(env, fmt.Sprintf("OPTS=--install --cd --reboot"))
	cmd.Env = env
	if nodeInitVerbose {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash -x cdxvirt-coreos-install.sh
	cmd = exec.Command("/bin/bash", "cdxvirt-coreos-install.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	if nodeInitVerbose {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash upload-preinit-scripts.sh
	cmd = exec.Command("/bin/bash", "upload-preinit-scripts.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
