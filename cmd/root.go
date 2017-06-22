package cmd

import (
	"fmt"
	"os"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ep      string
	cfgFile string
)

// NodeData test...
type NodeData struct {
	NodeNum int `json:"node_num"`
	Nodes   []struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
		Name   string `json:"name"`
	} `json:"nodes"`
}

// SettingData test...
type SettingData struct {
	SetMaxOSD string `json:"set_max_osd"`
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetEnvPrefix("cdxctl")
	viper.SetDefault("install_path", "/usr/local/bin")
	viper.BindEnv("install_path")

	RootCmd.PersistentFlags().StringVarP(&ep, "endpoint", "e", "127.0.0.1:5001", "Endpoint")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cdxctl/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		fmt.Println(cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")        // name of config file (without extension)
		viper.AddConfigPath("$HOME/.cdxctl") // adding home directory as first search path
		viper.AutomaticEnv()                 // read in environment variables that match
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		endpoint := viper.Get("endpoint")
		switch endpoint.(type) {
		case string:
			ep = endpoint.(string)
		}
	} else {
		fmt.Println(err)
	}
}

// RootCmd is the root CLI command
var RootCmd = &cobra.Command{
	Use:           "cdxctl",
	Short:         "cdxvirt platfrom CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runCommand(tmp_cmd string, interactive bool) {
	split_cmd := strings.Split(tmp_cmd, " ")
	cmd := exec.Command(split_cmd[0], split_cmd[1:]...)
	if interactive {
		cmd.Stdin = os.Stdin
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
