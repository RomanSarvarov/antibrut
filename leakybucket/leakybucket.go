package leakybucket

import (
	"context"

	"github.com/pkg/errors"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
)

type Service struct {
	repo repository
}

type repository interface {
	FindLimitation(ctx context.Context, c antibrut.LimitationCode) (*antibrut.Limitation, error)

	FindBucket(ctx context.Context, c antibrut.LimitationCode, val string) (*antibrut.Bucket, error)
	CreateBucket(ctx context.Context, bucket *antibrut.Bucket) (*antibrut.Bucket, error)
	DeleteBuckets(ctx context.Context, filter antibrut.BucketFilter) (int64, error)

	FindAttempts(ctx context.Context, filter antibrut.AttemptFilter) ([]*antibrut.Attempt, error)
	CreateAttempt(ctx context.Context, attempt *antibrut.Attempt) (*antibrut.Attempt, error)
}

func New(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Check(ctx context.Context, c antibrut.LimitationCode, val string) error {
	limit, err := s.repo.FindLimitation(ctx, c)
	if err != nil {
		return err
	}

	bucket, err := s.repo.FindBucket(ctx, limit.Code, val)
	if err != nil {
		if !errors.Is(err, antibrut.ErrNotFound) {
			return err
		}

		bucket, err = s.repo.CreateBucket(ctx, &antibrut.Bucket{
			LimitationCode: c,
			Value:          val,
		})
		if err != nil {
			return err
		}
	}

	attempts, err := s.repo.FindAttempts(ctx, antibrut.AttemptFilter{
		BucketID:      bucket.ID,
		CreatedAtFrom: clock.Now().Add(-limit.Interval.ToDuration()),
		CreatedAtTo:   clock.Now(),
	})
	if err != nil {
		return err
	}

	if len(attempts) >= limit.MaxAttempts {
		return antibrut.ErrMaxAttemptsExceeded
	}

	_, err = s.repo.CreateAttempt(ctx, &antibrut.Attempt{
		BucketID: bucket.ID,
	})

	return err
}

func (s *Service) Reset(ctx context.Context, filter antibrut.ResetFilter) error {
	_, err := s.repo.DeleteBuckets(ctx, filter)
	return err
}
