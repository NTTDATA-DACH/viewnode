package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"sync"
)

var (
	stp *Setup
	mu  sync.RWMutex
)

func Initialize(cmd *cobra.Command) (*Setup, error) {
	mu.Lock()
	defer mu.Unlock()

	setup := &Setup{}
	if kubeconfigFlag := cmd.Flags().Lookup("kubeconfig"); kubeconfigFlag != nil {
		setup.KubeCfgPath = kubeconfigFlag.Value.String()
	} else {
		setup.KubeCfgPath = ""
	}
	err := setup.Initialize()
	if err != nil {
		stp = nil
		return nil, fmt.Errorf("failed to initialize setup (%w)", err)
	}
	stp = setup
	return setup, nil
}

func GetConfig() *Setup {
	mu.RLock()
	defer mu.RUnlock()

	if stp == nil {
		panic("Setup not initialized! Call Initialize() before accessing GetConfig()")
	}
	return stp
}
