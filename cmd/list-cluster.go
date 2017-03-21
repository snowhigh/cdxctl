package cmd

import (
	"fmt"
	"os"
        "log"
	"database/sql"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listClusterCommand)
}

var listClusterCommand = &cobra.Command{
        Use:   "list-cluster",
        Short: "list cluster",
        RunE: func(cmd *cobra.Command, args []string) error {
		// Open db for read/write scan result
		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		sqlStmt := `create table if not exists cluster(
			vlan_id integer default 0,
			ipv4 varchar(50),
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
		// Print output
		w := tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', 0)
		fmt.Fprintln(w, "IP\tHOSTNAME\tCLUSTER")

		rows, err := db.Query("select ipv4, hostname, cluster from cluster order by cluster asc")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var ip string
			var hostname string
			var cluster string
			err = rows.Scan(&ip, &hostname, &cluster)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", ip, hostname, cluster)
		}
		w.Flush()
		return nil
        },
}
