package leakybucket

import (
	"context"

	"github.com/pkg/errors"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
)

// Service предоставляет сервис для работы с Leaky Bucket алгоритмом.
type Service struct {
	repo Repository
}

// Repository декларирует нужные методы для работы с БД.
type Repository interface {
	// FindLimitation находит antibrut.Limitation.
	// Если совпадений нет, вернет antibrut.ErrNotFound.
	FindLimitation(ctx context.Context, c antibrut.LimitationCode) (*antibrut.Limitation, error)

	// FindBucket находит antibrut.Bucket.
	// Если совпадений нет, вернет antibrut.ErrNotFound.
	FindBucket(ctx context.Context, c antibrut.LimitationCode, val string) (*antibrut.Bucket, error)

	// CreateBucket создает antibrut.Bucket.
	CreateBucket(ctx context.Context, bucket *antibrut.Bucket) (*antibrut.Bucket, error)

	// DeleteBuckets удаляет нужные antibrut.Bucket.
	DeleteBuckets(ctx context.Context, filter antibrut.BucketFilter) (int64, error)

	// FindAttempts находит совпадающие antibrut.Attempt.
	FindAttempts(ctx context.Context, filter antibrut.AttemptFilter) ([]*antibrut.Attempt, error)

	// CreateAttempt создает antibrut.Attempt.
	CreateAttempt(ctx context.Context, attempt *antibrut.Attempt) (*antibrut.Attempt, error)
}

// New создает Service.
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Check проверяет "хороший" ли запрос, или его следует отклонить.
// Используется алгоритм Leaky Bucket: https://en.wikipedia.org/wiki/Leaky_bucket.
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

// Reset сбрасывает бакеты по определенным признакам.
func (s *Service) Reset(ctx context.Context, filter antibrut.ResetFilter) error {
	_, err := s.repo.DeleteBuckets(ctx, filter)
	return err
}
