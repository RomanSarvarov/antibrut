package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResetTimeNowFunc(t *testing.T) {
	SetTimeNowFunc(func() time.Time {
		return time.Time{}
	})

	ResetTimeNowFunc()

	require.Equal(t, time.Now().Second(), timeNowFunc().Second())
}

func TestSetTimeNowFunc(t *testing.T) {
	t.Cleanup(func() { ResetTimeNowFunc() })

	SetTimeNowFunc(func() time.Time {
		return time.Time{}
	})

	require.NotEqual(t, time.Now().Second(), timeNowFunc().Second())
}

func TestNow(t *testing.T) {
	require.Equal(t, time.Now().Second(), Now().Second())
}

func TestNewFromTime(t *testing.T) {
	now := time.Now()

	require.Equal(t, now, NewFromTime(now))
}
