package cmd

import (
	"fmt"
        "log"
	"database/sql"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(joinClusterCommand)
}

var joinClusterCommand = &cobra.Command{
        Use:   "join-cluster CLUSTER IP",
        Short: "add device into cluster",
        RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Help()
		}
		if len(args) > 2 {
			return fmt.Errorf("Only one IP can be added at a time")
		}
		clusterID := args[0]
		nodeIP := args[1]
		// Open db for read/write scan result
		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		joinCluster(db, clusterID, nodeIP)
		log.Printf("Done")
		return nil
        },
}

func joinCluster(db *sql.DB, clusterID string, nodeIP string) {
	os.Chdir("provision")
	cmd := exec.Command("/bin/bash", "qts-qes-switcher.sh")
	env := os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash -x qts-ipxe-install.sh
	cmd = exec.Command("/bin/bash", "qts-ipxe-install.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash grub-update-reboot.sh CDRamfs
	cmd = exec.Command("/bin/bash", "grub-update-reboot.sh", "CDRamfs")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash -x cdxvirt-coreos-install.sh
	cmd = exec.Command("/bin/bash", "cdxvirt-coreos-install.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash upload-preinit-scripts.sh
	cmd = exec.Command("/bin/bash", "upload-preinit-scripts.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// HOST_IP_LIST="$HOSTS" bash -x k8sup.sh
	cmd = exec.Command("/bin/bash", "k8sup.sh")
	env = os.Environ()
	env = append(env, fmt.Sprintf("HOST_IP_LIST=%s", nodeIP))
	env = append(env, fmt.Sprintf("OPTS=--network=192.168.32.0/23 --cluster=%s", clusterID))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
