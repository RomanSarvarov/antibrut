package inmem

import (
	"context"
	"testing"
	"time"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
	"github.com/stretchr/testify/require"
)

func TestRepository_FindBucket(t *testing.T) {
	t.Parallel()

	t.Run("not found by limit code", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		limit, err := repo.FindBucket(ctx, "foo", "bar")
		require.Nil(t, limit)
		require.ErrorIs(t, err, antibrut.ErrNotFound)
	})

	t.Run("not found by limit code and value", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()
		bucket := &antibrut.Bucket{
			LimitationCode: "foo",
		}

		_, err := repo.CreateBucket(ctx, bucket)
		require.NoError(t, err)

		limit, err := repo.FindBucket(ctx, "foo", "bar")
		require.Nil(t, limit)
		require.ErrorIs(t, err, antibrut.ErrNotFound)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		tm := clock.NewFromTime(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))

		wantBucket := &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "bar",
			CreatedAt:      tm,
		}

		repo := New()
		repo.buckets = map[antibrut.LimitationCode][]*antibrut.Bucket{
			"foo": {wantBucket},
		}

		gotBucket, err := repo.FindBucket(ctx, "foo", "bar")
		require.NoError(t, err)

		require.NotNil(t, gotBucket.ID)
		require.Equal(t, wantBucket.LimitationCode, gotBucket.LimitationCode)
		require.Equal(t, wantBucket.Value, gotBucket.Value)
		require.Equal(t, wantBucket.CreatedAt, gotBucket.CreatedAt)
	})
}

func TestRepository_CreateBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo := New()
		repo.lastBucketID = 4

		now := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		clock.SetTimeNowFunc(func() time.Time {
			return now
		})
		t.Cleanup(func() { clock.ResetTimeNowFunc() })

		wantBucket := &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "bar",
		}

		gotBucket, err := repo.CreateBucket(ctx, wantBucket)
		require.NoError(t, err)

		require.NotNil(t, gotBucket.ID)
		require.Equal(t, wantBucket.LimitationCode, gotBucket.LimitationCode)
		require.Equal(t, wantBucket.Value, gotBucket.Value)
		require.Equal(t, wantBucket.CreatedAt, gotBucket.CreatedAt)

		require.Len(t, repo.buckets["foo"], 1)
		require.Equal(t, antibrut.BucketID(5), repo.lastBucketID)
	})
}

func TestRepository_DeleteBuckets(t *testing.T) {
	t.Parallel()

	t.Run("delete by limitation code", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		_, err := repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "10",
		})
		require.NoError(t, err)

		_, err = repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: "bar",
			Value:          "20",
		})
		require.NoError(t, err)

		n, err := repo.DeleteBuckets(ctx, antibrut.BucketFilter{
			LimitationCode: "foo",
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		require.Len(t, repo.buckets["foo"], 0)
		_, hasKey := repo.buckets["foo"]
		require.False(t, hasKey)

		require.Len(t, repo.buckets["bar"], 1)
	})

	t.Run("delete by value", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		_, err := repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "10",
		})
		require.NoError(t, err)

		_, err = repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "20",
		})
		require.NoError(t, err)

		n, err := repo.DeleteBuckets(ctx, antibrut.BucketFilter{
			Value: "10",
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		require.Len(t, repo.buckets["foo"], 1)
		require.Equal(t, repo.buckets["foo"][0].Value, "20")
	})

	t.Run("delete by created at to", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		bucket, err := repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "foo",
			CreatedAt:      time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		_, err = repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "bar",
			CreatedAt:      time.Date(2022, 5, 1, 13, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		n, err := repo.DeleteBuckets(ctx, antibrut.BucketFilter{
			CreatedAtTo: bucket.CreatedAt,
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		require.Len(t, repo.buckets["foo"], 1)
		require.Equal(t, repo.buckets["foo"][0].Value, "bar")
	})
}

func TestRepository_FindAttempts(t *testing.T) {
	t.Parallel()

	t.Run("find by bucket id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		attempt, err := repo.CreateAttempt(ctx, &antibrut.Attempt{
			BucketID: 5,
		})
		require.NoError(t, err)

		_, err = repo.CreateAttempt(ctx, &antibrut.Attempt{
			BucketID: 10,
		})
		require.NoError(t, err)

		attempts, err := repo.FindAttempts(ctx, antibrut.AttemptFilter{
			BucketID: attempt.BucketID,
		})
		require.NoError(t, err)
		require.Len(t, attempts, 1)

		require.Equal(t, attempt, attempts[0])
	})

	t.Run("find by created at from", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		_, err := repo.CreateAttempt(ctx, &antibrut.Attempt{
			CreatedAt: time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		attempt, err := repo.CreateAttempt(ctx, &antibrut.Attempt{
			CreatedAt: time.Date(2022, 5, 1, 13, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		attempts, err := repo.FindAttempts(ctx, antibrut.AttemptFilter{
			CreatedAtFrom: attempt.CreatedAt,
		})
		require.NoError(t, err)
		require.Len(t, attempts, 1)

		require.Equal(t, attempt, attempts[0])
	})

	t.Run("find by created at to", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		repo := New()

		attempt, err := repo.CreateAttempt(ctx, &antibrut.Attempt{
			CreatedAt: time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		_, err = repo.CreateAttempt(ctx, &antibrut.Attempt{
			CreatedAt: time.Date(2022, 5, 1, 13, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		attempts, err := repo.FindAttempts(ctx, antibrut.AttemptFilter{
			CreatedAtTo: attempt.CreatedAt,
		})
		require.NoError(t, err)
		require.Len(t, attempts, 1)

		require.Equal(t, attempt, attempts[0])
	})
}

func TestRepository_CreateAttempt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo := New()
		repo.lastAttemptID = 4

		now := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
		clock.SetTimeNowFunc(func() time.Time {
			return now
		})
		t.Cleanup(func() { clock.ResetTimeNowFunc() })

		wantAttempt := &antibrut.Attempt{
			BucketID: 10,
		}

		gotAttempt, err := repo.CreateAttempt(ctx, wantAttempt)
		require.NoError(t, err)

		require.NotNil(t, gotAttempt.ID)
		require.Equal(t, wantAttempt.BucketID, gotAttempt.BucketID)
		require.Equal(t, now, gotAttempt.CreatedAt)

		require.Len(t, repo.attempts[10], 1)
		require.Equal(t, antibrut.AttemptID(5), repo.lastAttemptID)
	})
}
