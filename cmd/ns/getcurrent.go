package ns

import "github.com/spf13/cobra"

var getCurrent = &cobra.Command{
	Use:     "get-current",
	Short:   "Get current Kubernetes namespace",
	Aliases: []string{"getcurrent"},
}
