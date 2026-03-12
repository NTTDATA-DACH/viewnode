package ns

import "github.com/spf13/cobra"

var NsCmd = &cobra.Command{
	Use:   "ns",
	Short: "Manage Kubernetes namespaces",
	Long:  `Manage Kubernetes namespaces and namespace-related operations.`,
}

func init() {
	NsCmd.AddCommand(listCmd)
}
