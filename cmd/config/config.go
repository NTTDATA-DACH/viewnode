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

	if stp == nil {
		stp = &Setup{}
	}
	if kubeconfigFlag := cmd.Flags().Lookup("kubeconfig"); kubeconfigFlag != nil {
		stp.KubeCfgPath = kubeconfigFlag.Value.String()
	} else {
		stp.KubeCfgPath = ""
	}
	err := stp.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize setup (%w)", err)
	}
	return stp, nil
}

func GetConfig() *Setup {
	mu.RLock()
	defer mu.RUnlock()

	if stp == nil {
		panic("Setup not initialized! Call Initialize() before accessing GetConfig()")
	}
	return stp
}
