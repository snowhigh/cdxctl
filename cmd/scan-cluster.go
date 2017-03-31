package cmd

import (
	"fmt"
        "log"
	"time"
	"database/sql"
	"strings"
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/grandcat/zeroconf"
)

func init() {
	RootCmd.AddCommand(scanClusterCommand)
}

var scanClusterCommand = &cobra.Command{
        Use:   "scan-cluster",
        Short: "scan for k8sup clusters",
        RunE: func(cmd *cobra.Command, args []string) error {
		// Open db for read/write scan result
		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		sqlStmt := `create table if not exists cluster_member(
			vlan_id integer default 0,
			ipv4 varchar(40),
			hostname varchar(100),
			cluster varchar(100),
			last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			unique(vlan_id, ipv4) on conflict replace
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return err
		}

		sqlStmt = `create table if not exists cluster(
			cluster varchar(100),
			version varchar(40),
			last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			unique(cluster) on conflict replace
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return err
		}
		scanCluster(db)
		log.Printf("Done")
		return nil
        },
}

func scanCluster(db *sql.DB) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			var hostname string
			var ipv4 string
			var clusterID string
			clusterID = "Unknown"
			for _, b := range entry.Text {
				if strings.Contains(b, "clusterID") {
					clusterID = strings.Split(b, "=")[1]
				}
			}
			hostname = strings.Split(entry.HostName, ".")[0]
			ipv4 = fmt.Sprintf("%s", entry.AddrIPv4[0])
			// fmt.Printf("%s %s:%d %s %s\n", e.HostName, e.AddrIPv4, e.Port, e.Text, e.ServiceInstanceName())
			// cluster member
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := tx.Prepare("insert or replace into cluster_member(vlan_id, ipv4, hostname, cluster) values(?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(0, ipv4, hostname, clusterID)
			if err != nil {
				log.Fatal(err)
			}
			tx.Commit()
			// cluster
			tx, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			stmt, err = tx.Prepare("insert or replace into cluster(cluster, version) values(?,?)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(clusterID, "Unknown")
			if err != nil {
				log.Fatal(err)
			}
			tx.Commit()
		}
	}(entries)
	// Send the "stop browsing" signal after the desired timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = resolver.Browse(ctx, "_etcd._tcp", "local", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(5 * time.Second)
}
