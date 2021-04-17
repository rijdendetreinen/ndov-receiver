package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// Version contains the version information
var Version VersionInformation

// VersionInformation is a simple struct containing the version information
type VersionInformation struct {
	Version string
	Commit  string
	Date    string
}

// VersionStringLong returns a version string
func (v VersionInformation) VersionStringLong() string {
	return fmt.Sprintf("%v (%v; built %v)", v.Version, v.Commit, v.Date)
}

// VersionStringShort returns a shortened version string
func (v VersionInformation) VersionStringShort() string {
	return fmt.Sprintf("%v (%v)", v.Version, v.Commit)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ndov-receiver",
	Short: "Receive NDOV messages",
	Long:  "Receive NDOV messages and queue them",
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")

		viper.AddConfigPath("./config")
		viper.AddConfigPath("/etc/ndov-receiver/")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Info("Using config file:")
	}

	log.Debug("Configuration loaded")
}
