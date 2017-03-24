package cmd

import (
	"fmt"
        "log"
	"database/sql"
	"os"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var clusterVersion string

func init() {
        createClusterCommand.Flags().StringVarP(&clusterVersion, "version", "v", "v1.10.1", "Define Cluster Version.")
	RootCmd.AddCommand(createClusterCommand)
}

var createClusterCommand = &cobra.Command{
        Use:   "create-cluster",
        Short: "create kubernetes cluster",
        RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one cluster can be created at a time")
		}
		clusterID := args[0]

		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		sqlStmt := `create table if not exists cluster(
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

		createCluster(db, clusterID)
		return nil
        },
}

func createCluster(db *sql.DB, clusterID string) {
	rows, err := db.Query(fmt.Sprintf("select cluster from cluster where cluster='%s'", clusterID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		log.Printf("Cluster %s already exist\n", clusterID)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert or replace into cluster(cluster, version) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(clusterID, clusterVersion)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	// Print output
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CLUSTER\tVERSION\tNODES")
	fmt.Fprintf(w, "%s\t%s\t0\n", clusterID, clusterVersion)
	w.Flush()
}

