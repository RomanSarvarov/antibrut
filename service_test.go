package antibrut_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mockery "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
	"github.com/romsar/antibrut/mock"
)

type mocks struct {
	r  *mock.Repository
	rl *mock.RateLimiter
}

func newFakeService(t *testing.T, opts ...antibrut.Option) (*antibrut.Service, mocks) {
	r := mock.NewRepository(t)
	rl := mock.NewRateLimiter(t)
	s := antibrut.NewService(r, rl, opts...)

	return s, mocks{
		r:  r,
		rl: rl,
	}
}

func TestService_AddIPToWhiteList(t *testing.T) {
	t.Parallel()

	t.Run("already exists", func(t *testing.T) {
		s, m := newFakeService(t)

		subnet := antibrut.Subnet("127.0.0.1/10")

		m.r.
			On(
				"FindIPRuleBySubnet",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				subnet,
			).
			Return(&antibrut.IPRule{
				Type:   antibrut.WhiteList,
				Subnet: "127.0.0.1/10",
			}, nil).
			Once()

		err := s.AddIPToWhiteList(context.Background(), "127.0.0.1/10")
		require.NoError(t, err)
	})

	t.Run("create", func(t *testing.T) {
		s, m := newFakeService(t)

		subnet := antibrut.Subnet("127.0.0.1/10")

		m.r.
			On(
				"FindIPRuleBySubnet",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				subnet,
			).
			Return(nil, antibrut.ErrNotFound).
			Once()

		m.r.
			On(
				"CreateIPRule",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				&antibrut.IPRule{
					Type:   antibrut.WhiteList,
					Subnet: subnet,
				},
			).
			Return(&antibrut.IPRule{
				Type:   antibrut.WhiteList,
				Subnet: "127.0.0.1/10",
			}, nil).
			Once()

		err := s.AddIPToWhiteList(context.Background(), "127.0.0.1/10")
		require.NoError(t, err)
	})
}

func TestService_DeleteIPFromWhiteList(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		subnet := antibrut.Subnet("127.0.0.1/10")

		s, m := newFakeService(t)

		m.r.
			On(
				"DeleteIPRules",
				ctx,
				antibrut.IPRuleFilter{
					Type:   antibrut.WhiteList,
					Subnet: subnet,
				},
			).
			Return(int64(5), nil).
			Once()

		err := s.DeleteIPFromWhiteList(ctx, subnet)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		subnet := antibrut.Subnet("127.0.0.1/10")

		gotErr := errors.New("foo bar")

		s, m := newFakeService(t)

		m.r.
			On(
				"DeleteIPRules",
				ctx,
				antibrut.IPRuleFilter{
					Type:   antibrut.WhiteList,
					Subnet: subnet,
				},
			).
			Return(int64(0), gotErr).
			Once()

		err := s.DeleteIPFromWhiteList(ctx, subnet)
		require.ErrorIs(t, err, gotErr)
	})
}

func TestService_AddIPToBlackList(t *testing.T) {
	t.Parallel()

	t.Run("already exists", func(t *testing.T) {
		s, m := newFakeService(t)

		subnet := antibrut.Subnet("127.0.0.1/10")

		m.r.
			On(
				"FindIPRuleBySubnet",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				subnet,
			).
			Return(&antibrut.IPRule{
				Type:   antibrut.BlackList,
				Subnet: "127.0.0.1/10",
			}, nil).
			Once()

		err := s.AddIPToBlackList(context.Background(), "127.0.0.1/10")
		require.NoError(t, err)
	})

	t.Run("create", func(t *testing.T) {
		s, m := newFakeService(t)

		subnet := antibrut.Subnet("127.0.0.1/10")

		m.r.
			On(
				"FindIPRuleBySubnet",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				subnet,
			).
			Return(nil, antibrut.ErrNotFound).
			Once()

		m.r.
			On(
				"CreateIPRule",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				&antibrut.IPRule{
					Type:   antibrut.BlackList,
					Subnet: subnet,
				},
			).
			Return(&antibrut.IPRule{
				Type:   antibrut.BlackList,
				Subnet: "127.0.0.1/10",
			}, nil).
			Once()

		err := s.AddIPToBlackList(context.Background(), "127.0.0.1/10")
		require.NoError(t, err)
	})
}

func TestService_DeleteIPFromBlackList(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		subnet := antibrut.Subnet("127.0.0.1/10")

		s, m := newFakeService(t)

		m.r.
			On(
				"DeleteIPRules",
				ctx,
				antibrut.IPRuleFilter{
					Type:   antibrut.BlackList,
					Subnet: subnet,
				},
			).
			Return(int64(5), nil).
			Once()

		err := s.DeleteIPFromBlackList(ctx, subnet)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		subnet := antibrut.Subnet("127.0.0.1/10")

		gotErr := errors.New("foo bar")

		s, m := newFakeService(t)

		m.r.
			On(
				"DeleteIPRules",
				ctx,
				antibrut.IPRuleFilter{
					Type:   antibrut.BlackList,
					Subnet: subnet,
				},
			).
			Return(int64(0), gotErr).
			Once()

		err := s.DeleteIPFromBlackList(ctx, subnet)
		require.ErrorIs(t, err, gotErr)
	})
}

func TestLogin_IsZero(t *testing.T) {
	t.Parallel()

	require.True(t, antibrut.Login("").IsZero())
	require.False(t, antibrut.Login("foo").IsZero())
}

func TestPassword_IsZero(t *testing.T) {
	t.Parallel()

	require.True(t, antibrut.Password("").IsZero())
	require.False(t, antibrut.Password("foo").IsZero())
}

func TestIP_IsZero(t *testing.T) {
	t.Parallel()

	require.True(t, antibrut.IP("").IsZero())
	require.False(t, antibrut.IP("127.0.0.1").IsZero())
}

func TestSubnet_Contains(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		c, err := antibrut.Subnet("192.168.5.0/26").Contains("192.168.5.15")
		require.NoError(t, err)
		require.True(t, c)
	})

	t.Run("negative", func(t *testing.T) {
		c, err := antibrut.Subnet("192.168.5.0/26").Contains("192.168.6.15")
		require.NoError(t, err)
		require.False(t, c)
	})
}

func TestService_Work(t *testing.T) {
	t.Parallel()

	clock.SetTimeNowFunc(func() time.Time {
		return time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
	})
	t.Cleanup(func() { clock.ResetTimeNowFunc() })

	t.Run("no prune duration", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		s, _ := newFakeService(t)

		err := s.Work(ctx)
		require.NoError(t, err)
	})

	t.Run("with prune duration", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		d := clock.NewDurationFromTimeDuration(1 * time.Second)
		s, m := newFakeService(t, antibrut.WithPruneDuration(d))

		m.rl.
			On(
				"Reset",
				ctx,
				antibrut.ResetFilter{
					DateTo: clock.Now().Add(-d.ToDuration()),
				},
			).
			Return(nil).
			Once()

		err := s.Work(ctx)
		require.NoError(t, err)
	})
}

func TestNewService(t *testing.T) {
	s := antibrut.NewService(
		mock.NewRepository(t),
		mock.NewRateLimiter(t),
	)

	require.IsType(t, &antibrut.Service{}, s)
}

func TestService_Reset(t *testing.T) {
	t.Parallel()

	t.Run("login success", func(t *testing.T) {
		ctx := context.Background()

		login := antibrut.Login("foobar")

		s, m := newFakeService(t)

		m.rl.
			On(
				"Reset",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.ResetFilter{
					LimitationCode: antibrut.LoginLimitation,
					Value:          login.String(),
				},
			).
			Return(nil).
			Once()

		err := s.Reset(ctx, login, "")
		require.NoError(t, err)
	})

	t.Run("login error", func(t *testing.T) {
		ctx := context.Background()

		login := antibrut.Login("foobar")
		gotErr := errors.New("some error")

		s, m := newFakeService(t)

		m.rl.
			On(
				"Reset",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.ResetFilter{
					LimitationCode: antibrut.LoginLimitation,
					Value:          login.String(),
				},
			).
			Return(gotErr).
			Once()

		err := s.Reset(ctx, login, "")
		require.ErrorIs(t, err, gotErr)
	})

	t.Run("ip success", func(t *testing.T) {
		ctx := context.Background()

		ip := antibrut.IP("foobar")

		s, m := newFakeService(t)

		m.rl.
			On(
				"Reset",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.ResetFilter{
					LimitationCode: antibrut.IPLimitation,
					Value:          ip.String(),
				},
			).
			Return(nil).
			Once()

		err := s.Reset(ctx, "", ip)
		require.NoError(t, err)
	})
}

func TestService_Check(t *testing.T) {
	t.Parallel()

	t.Run("login success", func(t *testing.T) {
		ctx := context.Background()

		login := antibrut.Login("foobar")

		s, m := newFakeService(t)

		m.rl.
			On(
				"Check",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.LoginLimitation,
				login.String(),
			).
			Return(nil).
			Once()

		err := s.Check(ctx, login, "", "")
		require.NoError(t, err)
	})

	t.Run("login error", func(t *testing.T) {
		ctx := context.Background()

		login := antibrut.Login("foobar")
		gotErr := errors.New("some error")

		s, m := newFakeService(t)

		m.rl.
			On(
				"Check",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.LoginLimitation,
				login.String(),
			).
			Return(gotErr).
			Once()

		err := s.Check(ctx, login, "", "")
		require.ErrorIs(t, err, gotErr)
	})

	t.Run("password success", func(t *testing.T) {
		ctx := context.Background()

		password := antibrut.Password("foobar")

		s, m := newFakeService(t)

		m.rl.
			On(
				"Check",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.PasswordLimitation,
				password.String(),
			).
			Return(nil).
			Once()

		err := s.Check(ctx, "", password, "")
		require.NoError(t, err)
	})

	t.Run("ip success", func(t *testing.T) {
		ctx := context.Background()

		ip := antibrut.IP("192.168.5.15")

		s, m := newFakeService(t)

		m.r.
			On(
				"FindIPRulesByIP",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				ip,
			).
			Return([]*antibrut.IPRule{}, nil).
			Once()

		m.rl.
			On(
				"Check",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.IPLimitation,
				ip.String(),
			).
			Return(nil).
			Once()

		err := s.Check(ctx, "", "", ip)
		require.NoError(t, err)
	})

	t.Run("ip with white list", func(t *testing.T) {
		ctx := context.Background()

		ip := antibrut.IP("192.168.5.15")

		s, m := newFakeService(t)

		m.r.
			On(
				"FindIPRulesByIP",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				ip,
			).
			Return([]*antibrut.IPRule{{
				Type:   antibrut.WhiteList,
				Subnet: "192.168.5.0/26",
			}}, nil).
			Once()

		err := s.Check(ctx, "", "", ip)
		require.NoError(t, err)
	})

	t.Run("ip with black list", func(t *testing.T) {
		ctx := context.Background()

		ip := antibrut.IP("192.168.5.15")

		s, m := newFakeService(t)

		m.r.
			On(
				"FindIPRulesByIP",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				ip,
			).
			Return([]*antibrut.IPRule{{
				Type:   antibrut.BlackList,
				Subnet: "192.168.5.0/26",
			}}, nil).
			Once()

		err := s.Check(ctx, "", "", ip)
		require.ErrorIs(t, err, antibrut.ErrIPInBlackList)
	})
}
