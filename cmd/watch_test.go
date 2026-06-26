package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRootCmdWatchFlagAbsentLeavesWatchDisabled(t *testing.T) {
	resetRootCommandState()

	watchLoopCalled := false
	runOnceFunc = func(_ context.Context) error {
		return nil
	}
	runWatchFunc = func(ctx context.Context, interval time.Duration, runOnce func(context.Context) error, sleep func(context.Context, time.Duration) error) error {
		watchLoopCalled = true
		return nil
	}

	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	RootCmd.SetArgs([]string{})

	err := RootCmd.Execute()

	require.NoError(t, err)
	require.False(t, watchEnabled)
	require.Zero(t, watchInterval)
	require.False(t, watchLoopCalled)
}

func TestRootCmdWatchFlagDefaultsToOneSecondWhenPresentWithoutValue(t *testing.T) {
	resetRootCommandState()

	var capturedInterval time.Duration
	runOnceFunc = func(_ context.Context) error {
		return nil
	}
	runWatchFunc = func(ctx context.Context, interval time.Duration, runOnce func(context.Context) error, sleep func(context.Context, time.Duration) error) error {
		capturedInterval = interval
		return nil
	}

	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	RootCmd.SetArgs([]string{"--watch"})

	err := RootCmd.Execute()

	require.NoError(t, err)
	require.True(t, watchEnabled)
	require.Equal(t, 1, watchInterval)
	require.Equal(t, time.Second, capturedInterval)
}

func TestRootCmdWatchFlagAcceptsExplicitIntervals(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{name: "long flag", args: []string{"--watch", "5"}},
		{name: "short flag", args: []string{"-w", "5"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetRootCommandState()

			var capturedInterval time.Duration
			runOnceFunc = func(_ context.Context) error {
				return nil
			}
			runWatchFunc = func(ctx context.Context, interval time.Duration, runOnce func(context.Context) error, sleep func(context.Context, time.Duration) error) error {
				capturedInterval = interval
				return nil
			}

			t.Cleanup(func() {
				resetRootCommandState()
				RootCmd.SetArgs(nil)
				RootCmd.SetOut(nil)
				RootCmd.SetErr(nil)
			})

			RootCmd.SetArgs(tc.args)

			err := RootCmd.Execute()

			require.NoError(t, err)
			require.True(t, watchEnabled)
			require.Equal(t, 5, watchInterval)
			require.Equal(t, 5*time.Second, capturedInterval)
		})
	}
}

func TestRootCmdWatchFlagRejectsZeroInterval(t *testing.T) {
	resetRootCommandState()

	var output bytes.Buffer
	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	RootCmd.SetErr(&output)
	RootCmd.SetArgs([]string{"--watch", "0"})

	err := RootCmd.Execute()

	require.ErrorContains(t, err, "must be >= 1 second")
}

func TestRootCmdWatchFlagRejectsNegativeInterval(t *testing.T) {
	resetRootCommandState()

	var output bytes.Buffer
	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	RootCmd.SetErr(&output)
	RootCmd.SetArgs([]string{"--watch=-1"})

	err := RootCmd.Execute()

	require.ErrorContains(t, err, "must be >= 1 second")
}

func TestRootCmdWatchFlagRejectsNonIntegerInterval(t *testing.T) {
	resetRootCommandState()

	var output bytes.Buffer
	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	RootCmd.SetErr(&output)
	RootCmd.SetArgs([]string{"--watch", "abc"})

	err := RootCmd.Execute()

	require.ErrorContains(t, err, "invalid argument")
	require.Contains(t, output.String(), "invalid argument")
}

func TestRunWatchSleepsAfterEachCompletedRefresh(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := 10 * time.Millisecond
	var (
		events     []string
		runEnds    []time.Time
		sleepCalls []struct {
			at       time.Time
			interval time.Duration
		}
	)

	runCount := 0
	err := runWatch(ctx, interval, func(context.Context) error {
		events = append(events, "run")
		time.Sleep(50 * time.Millisecond)
		runEnds = append(runEnds, time.Now())
		runCount++
		if runCount == 2 {
			cancel()
		}
		return nil
	}, func(ctx context.Context, duration time.Duration) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		events = append(events, "sleep")
		sleepCalls = append(sleepCalls, struct {
			at       time.Time
			interval time.Duration
		}{
			at:       time.Now(),
			interval: duration,
		})
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, []string{"sleep", "run", "sleep", "run"}, events)
	require.Len(t, sleepCalls, 2)
	require.Len(t, runEnds, 2)
	require.Equal(t, interval, sleepCalls[0].interval)
	require.Equal(t, interval, sleepCalls[1].interval)
	require.False(t, sleepCalls[1].at.Before(runEnds[0]))
}
