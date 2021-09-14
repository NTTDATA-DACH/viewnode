package main

import "viewnode/cmd"

var version string
var commit string

func main() {
	cmd.SetVersion(version)
	cmd.SetCommit(commit)
	cmd.Execute()
}
