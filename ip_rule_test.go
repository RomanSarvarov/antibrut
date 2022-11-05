package antibrut

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPRule_IsWhiteList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    IPRuleType
		want bool
	}{
		{
			name: "positive",
			t:    WhiteList,
			want: true,
		},
		{
			name: "negative",
			t:    IPRuleType(0),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rule := IPRule{
				Type: tt.t,
			}
			require.Equal(t, tt.want, rule.IsWhiteList())
		})
	}
}

func TestIPRule_IsBlackList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    IPRuleType
		want bool
	}{
		{
			name: "positive",
			t:    BlackList,
			want: true,
		},
		{
			name: "negative",
			t:    IPRuleType(0),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rule := IPRule{
				Type: tt.t,
			}
			require.Equal(t, tt.want, rule.IsBlackList())
		})
	}
}
