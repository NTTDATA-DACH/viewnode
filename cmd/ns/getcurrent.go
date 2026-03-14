package ns

import (
	"fmt"

	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"viewnode/cmd/config"
)

var currentRawConfig = func(setup *config.Setup) (clientcmdapi.Config, error) {
	return setup.ClientConfig.RawConfig()
}

var getCurrent = &cobra.Command{
	Use:     "get-current",
	Short:   "Get current Kubernetes namespace",
	Aliases: []string{"getcurrent"},
	RunE: func(c *cobra.Command, args []string) error {
		configCmd := c.Root()
		if configCmd == nil {
			configCmd = c
		}

		if _, err := initializeConfig(configCmd); err != nil {
			return err
		}
		setup := currentSetup()

		rawConfig, err := currentRawConfig(setup)
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		namespace := "default"
		if context := rawConfig.Contexts[rawConfig.CurrentContext]; context != nil && context.Namespace != "" {
			namespace = context.Namespace
		}

		fmt.Println(namespace)
		return nil
	},
}
