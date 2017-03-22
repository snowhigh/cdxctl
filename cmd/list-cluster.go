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
		fmt.Fprintln(w, "CLUSTER\tNODES")

		rows, err := db.Query("select distinct cluster from cluster")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var cluster string
			var count string
			err = rows.Scan(&cluster)
			if err != nil {
				log.Fatal(err)
			}
			crows, err := db.Query(fmt.Sprintf("select count(ipv4) from cluster where cluster='%s'", cluster))
			if err != nil {
				log.Fatal(err)
			}
			defer crows.Close()
			if crows.Next() {
				err = crows.Scan(&count)
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Fprintf(w, "%s\t%s\n", cluster, count)
		}
		w.Flush()
		return nil
        },
}
