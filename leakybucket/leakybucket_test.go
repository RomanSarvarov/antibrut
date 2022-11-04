package leakybucket

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
	mock "github.com/romsar/antibrut/mock/leakybucket" // leaky bucket mock
	mockery "github.com/stretchr/testify/mock"         // mockery lib
	"github.com/stretchr/testify/require"
)

type mocks struct {
	r *mock.Repository
}

func newFakeService(t *testing.T) (*Service, mocks) {
	t.Helper()

	r := mock.NewRepository(t)
	s := New(r)

	return s, mocks{
		r: r,
	}
}

func TestService_Reset(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		dt := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		filter := antibrut.ResetFilter{
			LimitationCode: antibrut.LimitationCode("foo"),
			Value:          "bar",
			CreatedAtTo:    clock.NewFromTime(dt),
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
		t.Parallel()

		ctx := context.Background()

		dt := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		filter := antibrut.ResetFilter{
			LimitationCode: antibrut.LimitationCode("foo"),
			Value:          "bar",
			CreatedAtTo:    clock.NewFromTime(dt),
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

func TestService_Check(t *testing.T) {
	t.Run("no limitation found", func(t *testing.T) {
		ctx := context.Background()

		s, m := newFakeService(t)

		m.r.
			On(
				"FindLimitation",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.LimitationCode("foo"),
			).
			Return(nil, antibrut.ErrNotFound).
			Once()

		err := s.Check(ctx, "foo", "bar")
		require.ErrorIs(t, err, antibrut.ErrNotFound)
	})

	t.Run("no bucket found", func(t *testing.T) {
		ctx := context.Background()

		s, m := newFakeService(t)

		now := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		clock.SetTimeNowFunc(func() time.Time {
			return now
		})
		t.Cleanup(func() { clock.ResetTimeNowFunc() })

		limit := &antibrut.Limitation{
			Code:        "foo",
			MaxAttempts: 1000,
			Interval:    clock.NewDurationFromTimeDuration(1 * time.Minute),
		}

		m.r.
			On(
				"FindLimitation",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				limit.Code,
			).
			Return(limit, nil).
			Once()

		m.r.
			On(
				"FindBucket",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				limit.Code,
				"bar",
			).
			Return(nil, antibrut.ErrNotFound).
			Once()

		bucket := &antibrut.Bucket{
			LimitationCode: limit.Code,
			Value:          "bar",
		}

		m.r.
			On(
				"CreateBucket",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				bucket,
			).
			Return(&antibrut.Bucket{
				ID:             5000,
				LimitationCode: limit.Code,
				Value:          "bar",
				CreatedAt:      clock.Now().Add(-5 * time.Second),
			}, nil).
			Once()

		attempts := []*antibrut.Attempt{
			{
				ID:        1000,
				BucketID:  bucket.ID,
				CreatedAt: clock.Now().Add(-5 * time.Second),
			},
		}

		m.r.
			On(
				"FindAttempts",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				antibrut.AttemptFilter{
					BucketID:      5000,
					CreatedAtFrom: clock.Now().Add(-limit.Interval.ToDuration()),
					CreatedAtTo:   clock.Now(),
				},
			).
			Return(attempts, nil).
			Once()

		m.r.
			On(
				"CreateAttempt",
				mockery.MatchedBy(func(_ context.Context) bool { return true }),
				&antibrut.Attempt{
					BucketID: 5000,
				},
			).
			Return(&antibrut.Attempt{
				ID:        9999,
				BucketID:  5000,
				CreatedAt: clock.Now(),
			}, nil).
			Once()

		err := s.Check(ctx, limit.Code, "bar")
		require.NoError(t, err)
	})
}

func TestService_CheckMaxAttemptsExceeded(t *testing.T) {
	tests := []struct {
		name           string
		gotLimitation  *antibrut.Limitation
		gotAttemptsCnt int
		wantErr        error
	}{
		{
			name: "positive",
			gotLimitation: &antibrut.Limitation{
				Code:        "foo",
				MaxAttempts: 10,
				Interval:    clock.NewDurationFromTimeDuration(1 * time.Minute),
			},
			gotAttemptsCnt: 10,
			wantErr:        antibrut.ErrMaxAttemptsExceeded,
		},
		{
			name: "negative",
			gotLimitation: &antibrut.Limitation{
				Code:        "foo",
				MaxAttempts: 10,
				Interval:    clock.NewDurationFromTimeDuration(5 * time.Second),
			},
			gotAttemptsCnt: 9,
			wantErr:        nil,
		},
		{
			name: "zero attempts cnt",
			gotLimitation: &antibrut.Limitation{
				Code:        "foo",
				MaxAttempts: 10,
				Interval:    clock.NewDurationFromTimeDuration(5 * time.Second),
			},
			gotAttemptsCnt: 0,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			s, m := newFakeService(t)

			now := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
			clock.SetTimeNowFunc(func() time.Time {
				return now
			})
			t.Cleanup(func() { clock.ResetTimeNowFunc() })

			m.r.
				On(
					"FindLimitation",
					mockery.MatchedBy(func(_ context.Context) bool { return true }),
					tt.gotLimitation.Code,
				).
				Return(tt.gotLimitation, nil).
				Once()

			bucket := &antibrut.Bucket{
				ID:             5000,
				LimitationCode: tt.gotLimitation.Code,
				Value:          "bar",
				CreatedAt:      clock.Now().Add(-5 * time.Second),
			}

			m.r.
				On(
					"FindBucket",
					mockery.MatchedBy(func(_ context.Context) bool { return true }),
					tt.gotLimitation.Code,
					"bar",
				).
				Return(bucket, nil).
				Once()

			attempts := make([]*antibrut.Attempt, 0, tt.gotAttemptsCnt)
			for i := 0; i < tt.gotAttemptsCnt; i++ {
				attempts = append(attempts, &antibrut.Attempt{
					ID:       antibrut.AttemptID(i),
					BucketID: bucket.ID,
				})
			}

			m.r.
				On(
					"FindAttempts",
					mockery.MatchedBy(func(_ context.Context) bool { return true }),
					antibrut.AttemptFilter{
						BucketID:      5000,
						CreatedAtFrom: clock.Now().Add(-tt.gotLimitation.Interval.ToDuration()),
						CreatedAtTo:   clock.Now(),
					},
				).
				Return(attempts, nil).
				Once()

			if tt.wantErr == nil {
				m.r.
					On(
						"CreateAttempt",
						mockery.MatchedBy(func(_ context.Context) bool { return true }),
						&antibrut.Attempt{
							BucketID: 5000,
						},
					).
					Return(&antibrut.Attempt{
						ID:        9999,
						BucketID:  5000,
						CreatedAt: clock.Now(),
					}, nil).
					Once()
			}

			err := s.Check(ctx, tt.gotLimitation.Code, "bar")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
