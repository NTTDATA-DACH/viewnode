package ns

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"viewnode/cmd/config"
	"viewnode/cmd/internal/listing"
)

var initializeConfig = config.Initialize
var currentSetup = config.GetConfig
var namespaceFilter string
var listNamespaces = func(ctx context.Context, setup *config.Setup) ([]string, error) {
	namespaces, err := setup.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaceNames := make([]string, 0, len(namespaces.Items))
	for _, namespace := range namespaces.Items {
		namespaceNames = append(namespaceNames, namespace.Name)
	}

	return namespaceNames, nil
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all Kubernetes namespaces",
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

		namespaceNames, err := listNamespaces(context.Background(), setup)
		if err != nil {
			return fmt.Errorf("getting kubernetes namespaces failed (%w)", err)
		}

		rawConfig, err := currentRawConfig(setup)
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		activeNamespace := "default"
		if context := rawConfig.Contexts[rawConfig.CurrentContext]; context != nil && context.Namespace != "" {
			activeNamespace = context.Namespace
		}

		for _, entry := range listing.PrepareNamespaceEntries(namespaceNames, activeNamespace) {
			marker := " "
			if entry.IsActive {
				marker = "*"
			}
			fmt.Printf("[%s] %s\n", marker, entry.Name)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&namespaceFilter, "filter", "f", "", "show only namespaces according to filter")
}
