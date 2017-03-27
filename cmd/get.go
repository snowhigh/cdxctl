// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resource (e.g. cdxctl get cluster)",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		client := &http.Client{
			Timeout: time.Duration(3e9),
		}

		if len(args) < 1 {
			log.Fatal(fmt.Errorf("%s\n", "No resource specified!"))
			os.Exit(1)
		}

		ResourceName := args[0]
		switch ResourceName {
		case "cluster":
			resp, err := client.Get("http://" + ep + "/api/v1/cluster")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			var nodeData NodeData
			err = json.Unmarshal(body, &nodeData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Node amount: ", nodeData.NodeNum)
			for _, no := range nodeData.Nodes {
				fmt.Println("Node name:", no.Name, "\t", "cpu:", no.CPU, "\t", "memory:", no.Memory)
			}

		case "setting":
			resp, err := client.Get("http://" + ep + "/api/v1/storage/setting")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			var settingData SettingData
			err = json.Unmarshal(body, &settingData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Max OSD: " + settingData.SetMaxOSD)
		default:
			log.Fatal(fmt.Errorf("%s\n", "Resource name error!"))
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
