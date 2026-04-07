package ctx

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"viewnode/cmd/config"
	"viewnode/cmd/internal/listing"
)

var initializeConfig = config.Initialize
var currentSetup = config.GetConfig
var contextFilter string

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all Kubernetes contexts",
	Aliases: []string{"ls"},
	RunE: func(c *cobra.Command, args []string) error {
		configCmd := c.Root()
		if configCmd == nil {
			configCmd = c
		}

		filterValue := ""
		if filterFlag := c.Flags().Lookup("filter"); filterFlag != nil && filterFlag.Changed {
			filterValue = normalizeContextFilter(filterFlag.Value.String())
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
		ctxs = filterContextNames(ctxs, filterValue)
		if filterValue != "" && len(ctxs) == 0 {
			fmt.Printf("no contexts matched filter %q\n", filterValue)
			return nil
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

func filterContextNames(contextNames []string, filter string) []string {
	if filter == "" {
		return contextNames
	}

	normalizedFilter := strings.ToLower(filter)
	filteredNames := make([]string, 0, len(contextNames))
	for _, contextName := range contextNames {
		if strings.Contains(strings.ToLower(contextName), normalizedFilter) {
			filteredNames = append(filteredNames, contextName)
		}
	}

	return filteredNames
}

func normalizeContextFilter(filter string) string {
	return strings.TrimSpace(filter)
}

func init() {
	listCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
}
