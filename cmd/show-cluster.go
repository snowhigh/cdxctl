package cmd

import (
	"fmt"
        "log"
	"database/sql"
	"text/tabwriter"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(showClusterCommand)
}

var showClusterCommand = &cobra.Command{
        Use:   "show-cluster CLUSTER",
        Short: "show cluster info",
        RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one cluster can be displayed at a time")
		}
		clusterID := args[0]
		// Open db for read/write scan result
		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		showCluster(db, clusterID)
		return nil
        },
}

func showCluster(db *sql.DB, clusterID string) {
	var ipv4 string
	var hostname string
	var version string
	var count string

	w := tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', 0)

	rows, err := db.Query(fmt.Sprintf("select version from cluster where cluster='%s'", clusterID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Fprintf(w, "CLUSTER\t%s\n", clusterID)
	fmt.Fprintf(w, "VERSION\t%s\n", version)

	rows, err = db.Query(fmt.Sprintf("select count(ipv4) from cluster_member where cluster='%s'", clusterID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Fprintf(w, "NODES\t%s\n", count)


	fmt.Fprintln(w, "HOSTNAME\tIPV4\tCPU\tMEM\tNIC")
	rows, err = db.Query(fmt.Sprintf("select hostname, ipv4 from cluster_member where cluster='%s'", clusterID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&hostname, &ipv4)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", hostname, ipv4, "1", "2048m", "enp3s0")
	}
	w.Flush()
}
