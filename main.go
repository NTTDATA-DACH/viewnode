package main

import "kubectl-view-node/cmd"

var version string
var commit string

func main() {
	cmd.SetVersion(version)
	cmd.SetCommit(commit)
	cmd.Execute()
}
