package ns

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"viewnode/cmd/config"
)

var initializeConfig = config.Initialize
var currentSetup = config.GetConfig
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
		sort.Strings(namespaceNames)

		for _, namespaceName := range namespaceNames {
			marker := " "
			if namespaceName == setup.Namespace {
				marker = "*"
			}
			fmt.Printf("[%s] %s\n", marker, namespaceName)
		}

		return nil
	},
}
