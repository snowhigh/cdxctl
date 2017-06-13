package cmd

import (
	"fmt"
        "log"
	"os"
	"os/exec"
	"encoding/json"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type NodeFacts struct {
	Fact struct {
		Distri string `json:"ansible_distribution"`
		DistriVersion string `json:"ansible_distribution_version"`
		Hostname string `json:"ansible_hostname"`
		Cores int `json:"ansible_processor_cores"`
		Memory int `json:"ansible_memtotal_mb"`
		DefaultIpv4 struct {
			Address string `json:"address"`
                        Interface string `json:"interface"`
			Gateway string `json:"gateway"`
		} `json:"ansible_default_ipv4"`
		Interfaces []interface{} `json:"ansible_interfaces"`
	} `json:"ansible_facts"`
}

var nodeShowIp string
var nodeShowVerbose bool

func init() {
	nodeShowCommand.Flags().StringVarP(&nodeShowIp, "ip", "i", "", "Node IP address")
	nodeShowCommand.Flags().BoolVarP(&nodeShowVerbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(nodeShowCommand)
}

var nodeShowCommand = &cobra.Command{
        Use:   "node-show",
        Short: "show node information",
        RunE: func(cmd *cobra.Command, args []string) error {
		if nodeShowIp == "" {
			return cmd.Help()
		}
		nodeIP := nodeShowIp
		nodeShow(nodeIP)
		log.Printf("Done")
		return nil
        },
}

func nodeShow(nodeIP string) {
	os.Chdir("/root/provision/")
        cmd := exec.Command("/bin/bash", "gather-info.sh")
        env := os.Environ()
        env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
        cmd.Env = env
	if nodeShowVerbose {
		cmd.Stdout = os.Stdout
        	cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	dat, err := ioutil.ReadFile(fmt.Sprintf("/tmp/facts/%s", nodeIP))
	if err != nil {
		log.Fatal(err)
	}
	facts := NodeFacts{}
	err = json.Unmarshal(dat, &facts)
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(facts, "", "\t")
	fmt.Println(string(b))
}
