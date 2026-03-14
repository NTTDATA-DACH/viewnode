package ns

import "github.com/spf13/cobra"

var setCurrent = &cobra.Command{
	Use:     "set-current NAMESPACE",
	Short:   "Set current Kubernetes namespace",
	Aliases: []string{"setcurrent"},
}
