package ctx

import (
	"fmt"
	"github.com/spf13/cobra"
	"viewnode/cmd/config"
	"viewnode/cmd/internal/listing"
)

var initializeConfig = config.Initialize
var currentSetup = config.GetConfig

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all Kubernetes contexts",
	Aliases: []string{"ls"},
	RunE: func(c *cobra.Command, args []string) error {
		configCmd := c.Root()
		if configCmd == nil {
			configCmd = c
		}
		if _, err := initializeConfig(configCmd); err != nil {
			return err
		}
		setup := currentSetup()
		rawConfig, err := setup.ClientConfig.RawConfig()
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		ctxs := make([]string, 0, len(rawConfig.Contexts))
		for k := range rawConfig.Contexts {
			ctxs = append(ctxs, k)
		}

		for _, entry := range listing.PrepareContextEntries(ctxs, rawConfig.CurrentContext) {
			marker := " "
			if entry.IsActive {
				marker = "*"
			}
			fmt.Printf("[%s] %s\n", marker, entry.Name)
		}
		return nil
	},
}
