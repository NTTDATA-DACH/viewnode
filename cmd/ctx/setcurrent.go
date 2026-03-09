package ctx

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"viewnode/cmd/config"
)

var setCurrent = &cobra.Command{
	Use:     "set-current CONTEXT",
	Short:   "Set current Kubernetes context",
	Aliases: []string{"sc", "setcurrent"},
	Args:    cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		setup := config.GetConfig()
		rawConfig, err := setup.ClientConfig.RawConfig()
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		contextName := args[0]
		if _, exists := rawConfig.Contexts[contextName]; !exists {
			return fmt.Errorf("kubernetes context %q not found", contextName)
		}

		rawConfig.CurrentContext = contextName
		if err := clientcmd.ModifyConfig(clientcmd.NewDefaultClientConfigLoadingRules(), rawConfig, false); err != nil {
			return fmt.Errorf("setting kubernetes current context failed (%w)", err)
		}

		fmt.Printf("current context set to %s\n", contextName)
		return nil
	},
}
