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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set resource (e.g. cdxctl set setting 1)",
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
		case "setting":

			if len(args) < 2 {
				log.Fatal(fmt.Errorf("%s\n", "No Max OSD amount specified!"))
				os.Exit(1)
			}
			maxOSD := args[1]

			url := "http://" + ep + "/api/v1/storage/setting"

			settingData := &SettingData{SetMaxOSD: maxOSD}
			jsonData, err := json.Marshal(settingData)
			if err != nil {
				log.Fatal(err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(body, &settingData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("The max OSD amount:", settingData.SetMaxOSD, "has been set!")

		default:
			log.Fatal(fmt.Errorf("%s\n", "Resource name error!"))
			os.Exit(1)
		}

	},
}

func init() {
	RootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
