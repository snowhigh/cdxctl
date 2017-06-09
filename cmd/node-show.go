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
		Hostname string `json:"ansible_hostname"`
		Memory int `json:"ansible_memtotal_mb"`
		DefaultIpv4 struct {
			Address string `json:"address"`
		} `json:"ansible_default_ipv4"`
		Distri string `json:"ansible_distribution"`
		DistriVersion string `json:"ansible_distribution_version"`
		Cores int `json:"ansible_processor_cores"`
	} `json:"ansible_facts"`
}

var Verbose bool

func init() {
	nodeShowCommand.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(nodeShowCommand)
}

var nodeShowCommand = &cobra.Command{
        Use:   "node-show",
        Short: "show node information",
        RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one IP can be added at a time")
		}
		nodeIP := args[0]

		nodeShow(nodeIP)
		log.Printf("Done")
		return nil
        },
}

func nodeShow(nodeIP string) {
	os.Chdir("provision/")
        cmd := exec.Command("/bin/bash", "gather-info.sh")
        env := os.Environ()
        env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
        cmd.Env = env
	if Verbose {
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
	b, _ := json.Marshal(facts)
	fmt.Println(string(b))
}
