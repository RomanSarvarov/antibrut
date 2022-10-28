package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewDurationFromTimeDuration(t *testing.T) {
	d := time.Duration(555)
	require.Equal(t, Duration{Duration: d}, NewDurationFromTimeDuration(d))
}

func TestDuration_ToDuration(t *testing.T) {
	d := time.Duration(555)
	require.Equal(t, d, Duration{Duration: d}.ToDuration())
}
