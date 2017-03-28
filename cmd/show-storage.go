package cmd

import (
	"fmt"
	"os"
        "log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"database/sql"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(showStorageCommand)
}

type StorageNode struct {
	Nodes []struct {
		Hwslot []struct {
			DevName string `json:"dev_name"`
			DiskFwrev string `json:"disk_fwrev"`
			DiskModel string `json:"disk_model"`
			DiskSerial string `json:"disk_serial"`
			DiskType string `json:"disk_type"`
			OsdID string `json:"osd_id"`
			Slot string `json:"slot"`
		} `json:"hwslot"`
		NodeName string `json:"node_name"`
		OsdName string `json:"osd_name"`
	} `json:"nodes"`
}

var showStorageCommand = &cobra.Command{
        Use:   "show-storage CLUSTERID",
        Short: "show storage",
        RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		clusterID := args[0]

		// Open db for read/write scan result
		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		w := tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', 0)
		fmt.Fprintln(w, "IP\tOSD\tDISK")
		// get cluster node ip
		rows, err := db.Query(fmt.Sprintf("select cluster, ipv4 from cluster_member where cluster='%s'", clusterID))

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

                var cluster string
                var ip string

		// loop to get storage OSD info use api
		for rows.Next() {
			err = rows.Scan(&cluster, &ip)
			if err != nil {
				log.Fatal(err)
			}

			client := &http.Client{}
			resp, err := client.Get("http://" + ip + ":30001/api/v1/storage/node/" + ip)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			data :=  StorageNode{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Fatal(err)
			}
			var OSD int
			var DISK int
			for k := range data.Nodes[0].Hwslot {
				if data.Nodes[0].Hwslot[k].OsdID != "DOM" {
					if data.Nodes[0].Hwslot[k].OsdID != "" {
						OSD++
					}
					DISK++
				}
			}
			fmt.Fprintf(w, "%s\t%d\t%d\n", ip, OSD, DISK )
			w.Flush()
		}

		return nil
        },
}
