package ns

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"viewnode/cmd/config"
)

var namespaceExists = func(ctx context.Context, setup *config.Setup, namespace string) (bool, error) {
	_, err := setup.Clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

var persistRawConfig = func(setup *config.Setup, rawConfig clientcmdapi.Config) error {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if setup.KubeCfgPath != "" {
		loadingRules.ExplicitPath = setup.KubeCfgPath
	}

	return clientcmd.ModifyConfig(loadingRules, rawConfig, false)
}

var setCurrent = &cobra.Command{
	Use:     "set-current NAMESPACE",
	Short:   "Set current Kubernetes namespace",
	Aliases: []string{"setcurrent"},
	Args:    cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		configCmd := c.Root()
		if configCmd == nil {
			configCmd = c
		}

		if _, err := initializeConfig(configCmd); err != nil {
			return err
		}
		setup := currentSetup()

		namespaceName := args[0]
		exists, err := namespaceExists(context.Background(), setup, namespaceName)
		if err != nil {
			return fmt.Errorf("validating kubernetes namespace %q failed (%w)", namespaceName, err)
		}
		if !exists {
			return fmt.Errorf("kubernetes namespace %q not found", namespaceName)
		}

		rawConfig, err := currentRawConfig(setup)
		if err != nil {
			return fmt.Errorf("getting kubernetes raw config failed (%w)", err)
		}

		currentContext := rawConfig.CurrentContext
		if currentContext == "" || rawConfig.Contexts[currentContext] == nil {
			return fmt.Errorf("current kubernetes context %q not found", currentContext)
		}

		rawConfig.Contexts[currentContext].Namespace = namespaceName
		if err := persistRawConfig(setup, rawConfig); err != nil {
			return fmt.Errorf("setting kubernetes current namespace failed (%w)", err)
		}

		fmt.Printf("current namespace set to %s\n", namespaceName)
		return nil
	},
}
