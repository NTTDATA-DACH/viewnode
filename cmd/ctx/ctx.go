package ctx

import (
	"github.com/spf13/cobra"
)

var CtxCmd = &cobra.Command{
	Use:   "ctx",
	Short: "Manage Kubernetes context",
	Long:  `Manage Kubernetes context, like listing, getting the current and setting new current context.`,
}

func init() {
	CtxCmd.AddCommand(listCmd, getCurrent, setCurrent)
}
