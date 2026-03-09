package ctx

import (
	"fmt"
	"github.com/spf13/cobra"
	"viewnode/cmd/config"
)

var getCurrent = &cobra.Command{
	Use:     "get-current",
	Short:   "Get current Kubernetes context",
	Aliases: []string{"gc", "getcurrent"},
	RunE: func(c *cobra.Command, args []string) error {
		setup := config.GetConfig()
		rawConfig, err := setup.ClientConfig.RawConfig()
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		fmt.Println(rawConfig.CurrentContext)
		return nil
	},
}
