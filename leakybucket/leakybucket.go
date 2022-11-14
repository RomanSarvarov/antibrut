package leakybucket

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/romsar/antibrut"
)

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

// Service предоставляет сервис для работы с Leaky Bucket алгоритмом.
type Service struct {
	repo Repository

	// timeNow содержит функцию, которая возвращает текущее время.
	timeNow func() time.Time

	// logger механизм логирования.
	logger logger
}

// logger это контракт для механизма логирования.
type logger interface {
	// Printf сохраняет отформатированное сообщение в лог.
	Printf(format string, v ...any)
}

// Option возвращает функцию, модифицирующую Service.
type Option func(s *Service)

// WithTimeNow возвращает функцию, устанавливающую
// callback для получения текущего времени.
func WithTimeNow(f func() time.Time) Option {
	return func(s *Service) {
		s.timeNow = f
	}
}

// WithLogger возвращает функцию,
// устанавливающую механизм логирования.
func WithLogger(l logger) Option {
	return func(s *Service) {
		s.logger = l
	}
}

// New создает Service.
func New(repo Repository, opts ...Option) *Service {
	s := &Service{
		repo: repo,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.timeNow == nil {
		s.timeNow = time.Now
	}

	if s.logger == nil {
		s.logger = log.New(io.Discard, "", 0)
	}

	return s
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
		CreatedAtFrom: s.timeNow().Add(-limit.Interval.ToDuration()),
		CreatedAtTo:   s.timeNow(),
	})
	if err != nil {
		return err
	}

	if len(attempts) >= limit.MaxAttempts {
		s.logger.Printf(
			"Request was blocked [Value: %s] [Limitation: %s] [Attempts: %d] [Max-attempts: %d].\n",
			val,
			c,
			len(attempts),
			limit.MaxAttempts,
		)

		return antibrut.ErrMaxAttemptsExceeded
	}

	_, err = s.repo.CreateAttempt(ctx, &antibrut.Attempt{
		BucketID: bucket.ID,
	})
	if err != nil {
		return err
	}

	s.logger.Printf("Request was allowed [Value: %s] [Limitation: %s].\n", val, c)

	return err
}

// Reset сбрасывает бакеты по определенным признакам.
func (s *Service) Reset(ctx context.Context, filter antibrut.ResetFilter) error {
	s.logger.Printf("Deleting buckets [Filter: %v].\n", filter)

	_, err := s.repo.DeleteBuckets(ctx, filter)

	return err
}
