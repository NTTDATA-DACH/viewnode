package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetVersionAndCommit(t *testing.T) {
	SetVersion("1.2.3")
	SetCommit("abc123")

	require.Equal(t, "1.2.3", version)
	require.Equal(t, "abc123", commit)
}

func TestCurrentBuildYear(t *testing.T) {
	originalBuildTime := buildTime
	t.Cleanup(func() {
		buildTime = originalBuildTime
	})

	buildTime = "2026-03-10T12:00:00Z"

	require.Equal(t, 2026, currentBuildYear())
}

func TestCurrentBuildYearFallsBackToCurrentYear(t *testing.T) {
	originalBuildTime := buildTime
	t.Cleanup(func() {
		buildTime = originalBuildTime
	})

	buildTime = "invalid"

	require.Equal(t, time.Now().Year(), currentBuildYear())
}

func TestVersionCommandPrintsVersionAndCommit(t *testing.T) {
	SetVersion("1.2.3")
	SetCommit("abc123")
	originalBuildTime := buildTime
	t.Cleanup(func() {
		buildTime = originalBuildTime
	})
	buildTime = "2026-03-10T12:00:00Z"

	output := captureStdout(t, func() {
		versionCmd.Run(versionCmd, nil)
	})

	require.Equal(t, "viewnode 1.2.3 (abc123) © 2026 NTT DATA Deutschland SE, Adam Boczek | source: https://github.com/NTTDATA-DACH/viewnode\n", output)
}
