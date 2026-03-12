package ns

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"viewnode/cmd/config"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all Kubernetes namespaces",
	Aliases: []string{"ls"},
	RunE: func(c *cobra.Command, args []string) error {
		configCmd := c.Root()
		if configCmd == nil {
			configCmd = c
		}

		if _, err := config.Initialize(configCmd); err != nil {
			return err
		}
		setup := config.GetConfig()

		namespaces, err := setup.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("getting kubernetes namespaces failed (%w)", err)
		}

		namespaceNames := make([]string, 0, len(namespaces.Items))
		for _, namespace := range namespaces.Items {
			namespaceNames = append(namespaceNames, namespace.Name)
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
