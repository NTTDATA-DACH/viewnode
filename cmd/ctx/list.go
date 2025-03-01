package ctx

import (
	"fmt"
	"github.com/spf13/cobra"
	"sort"
	"viewnode/cmd/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kubernetes contexts",
	RunE: func(c *cobra.Command, args []string) error {
		setup := config.GetConfig()
		rawConfig, err := setup.ClientConfig.RawConfig()
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		ctxs := make([]string, 0, len(rawConfig.Contexts))
		for k := range rawConfig.Contexts {
			ctxs = append(ctxs, k)
		}
		sort.Strings(ctxs)

		for _, ctx := range ctxs {
			marker := " "
			if ctx == rawConfig.CurrentContext {
				marker = "*"
			}
			fmt.Printf("[%s] %s\n", marker, ctx)
		}
		return nil
	},
}
