package cron_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"gitlab-mr-notifier/internal/cron"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	o := cron.New()
	require.NoError(t, o.Start("1d", "9:30", func() {}))
	require.NoError(t, o.Start("1w", "9:30", func() {}))

	require.Error(t, o.Start("1m", "9:30", func() {}))
	require.Error(t, o.Start("1d", "29:30", func() {}))
	require.Error(t, o.Start("1w", "29:30", func() {}))
	require.Error(t, o.Start("wd", "29:30", func() {}))
	require.Error(t, o.Start("fw", "29:30", func() {}))

	require.NoError(t, o.Start("1m", "", func() {}))
	require.NoError(t, o.Start("1s", "", func() {}))

	require.Error(t, o.Start("1d", "", func() {}))
	require.Error(t, o.Start("", "", nil))
}

func TestFunctional(t *testing.T) {
	o := cron.New()
	var val int32
	require.NoError(t, o.Start("5s", "", func() {
		fmt.Println("on job")
		atomic.AddInt32(&val, 1)
	}))

	// wait a second, should be still 0
	time.Sleep(time.Second)
	require.Equal(t, int32(0), atomic.LoadInt32(&val))

	// wait more, should be 1
	time.Sleep(5 * time.Second)
	require.Equal(t, int32(1), atomic.LoadInt32(&val))
}
