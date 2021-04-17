package main

import "github.com/rijdendetreinen/ndov-receiver/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = cmd.VersionInformation{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
	cmd.RootCmd.Version = cmd.Version.Version
	cmd.RootCmd.SetVersionTemplate("ndov-receiver " + cmd.Version.VersionStringLong() + "\n")
	cmd.Execute()
}
