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

var nodeSetBondIp string
var nodeSetBondName string
var nodeSetBondNics string

func init() {
	nodeSetBondCommand.Flags().StringVarP(&nodeSetBondIp, "ip", "i", "", "Node IP address")
	nodeSetBondCommand.Flags().StringVarP(&nodeSetBondName, "name", "n", "", "Node Bond dev name")
	nodeSetBondCommand.Flags().StringVarP(&nodeSetBondNics, "nics", "f", "", "Node NIC names ex. eth0,eth1")
	RootCmd.AddCommand(nodeSetBondCommand)
}

var nodeSetBondCommand = &cobra.Command{
        Use:   "node-set-bond",
        Short: "set node bond interfaces",
        RunE: func(cmd *cobra.Command, args []string) error {
		if nodeSetBondIp == "" || nodeSetBondName == "" || nodeSetBondNics == "" {
			return cmd.Help()
		}
		nodeIP := nodeSetBondIp
		name := nodeSetBondName
		nics := nodeSetBondNics
		fmt.Sprintf("nics: %s\n", nics)
		nicsa := strings.Split(nics, ",")
		nic1, nic2 := nicsa[0], nicsa[1]
		fmt.Sprintf("nic1: %s\n", nic1)
		fmt.Sprintf("nic2: %s\n", nic2)

		nodeSetBond(nodeIP, name, nic1, nic2)
		log.Printf("Done")
		return nil
        },
}

func nodeSetBond(nodeIP string, name string, nic1 string, nic2 string) {
	os.Chdir("provision/")
	cmd := exec.Command("/bin/bash", "set-bond.sh")
        env := os.Environ()
        env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	env = append(env, fmt.Sprintf("BONDDEV=%s", name))
	env = append(env, fmt.Sprintf("BONDNIC1=%s", nic1))
	env = append(env, fmt.Sprintf("BONDNIC2=%s", nic2))
        cmd.Env = env
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
