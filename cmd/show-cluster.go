package cmd

import (
	"fmt"
        "log"
	"database/sql"
	"text/tabwriter"
	"os"
	"os/exec"
	"strings"
	"io/ioutil"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type Facts struct {
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
	showClusterCommand.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
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


	fmt.Fprintln(w, "HOSTNAME\tIPV4\tOS\tVERSION\tCPU\tMEM\tNIC")
	rows, err = db.Query(fmt.Sprintf("select hostname, ipv4 from cluster_member where cluster='%s'", clusterID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var hosts []string
	for rows.Next() {
		err = rows.Scan(&hostname, &ipv4)
		if err != nil {
			log.Fatal(err)
		}
		hosts = append(hosts, ipv4)
	}
	// call ansible gather-info
	gather_info(hosts, w)
	w.Flush()
}

func gather_info(ips []string, w *tabwriter.Writer) {
	var hostname, ipv4, num_nic, distri, distri_version string
	var cores, mem int
	var iplist string
	for _, ip := range ips {
		iplist = iplist + " " + ip
	}

	os.Chdir("provision/")
        cmd := exec.Command("/bin/bash", "gather-info.sh")
        env := os.Environ()
        env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", iplist))
        cmd.Env = env
	if Verbose {
		cmd.Stdout = os.Stdout
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		dat, err := ioutil.ReadFile(fmt.Sprintf("/tmp/facts/%s", ip))
		if err != nil {
			log.Fatal(err)
		}
		facts := Facts{}
		err = json.Unmarshal(dat, &facts)
		if err != nil {
			log.Fatal(err)
		}
		hostname = facts.Fact.Hostname
		ipv4 = facts.Fact.DefaultIpv4.Address
		distri = facts.Fact.Distri
		distri_version = facts.Fact.DistriVersion
		cores = facts.Fact.Cores
		mem = facts.Fact.Memory

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t%s\n", hostname, ipv4, distri,
		distri_version, cores, mem, strings.Replace(num_nic, "\n", ",", -1))
	}
}
