package leakybucket

import (
	"context"
	"errors"
	"testing"
	"time"

	mockery "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
	mock "github.com/romsar/antibrut/mock/leakybucket"
)

type mocks struct {
	r *mock.Repository
}

func newFakeService(t *testing.T) (*Service, mocks) {
	r := mock.NewRepository(t)
	s := New(r)

	return s, mocks{
		r: r,
	}
}

func TestService_Reset(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		dt := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		filter := antibrut.ResetFilter{
			LimitationCode: antibrut.LimitationCode("foo"),
			Value:          "bar",
			DateTo:         clock.NewFromTime(dt),
		}

		s, m := newFakeService(t)

		m.r.
			On(
				"DeleteBuckets",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				filter,
			).
			Return(int64(5), nil).
			Once()

		err := s.Reset(ctx, filter)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()

		dt := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		filter := antibrut.ResetFilter{
			LimitationCode: antibrut.LimitationCode("foo"),
			Value:          "bar",
			DateTo:         clock.NewFromTime(dt),
		}
		gotErr := errors.New("some error")

		s, m := newFakeService(t)

		m.r.
			On(
				"DeleteBuckets",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				filter,
			).
			Return(int64(0), gotErr).
			Once()

		err := s.Reset(ctx, filter)
		require.ErrorIs(t, err, gotErr)
	})
}
