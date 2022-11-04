package inmem

import (
	"context"
	"sync"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
)

// Repository предоставляет API для работы с хранилищем.
type Repository struct {
	// lastBucketID последний идентификатор бакетов.
	lastBucketID antibrut.BucketID

	// bucketsMu мьютекс по бакетам.
	bucketsMu sync.Mutex

	// buckets содержит в себе бакеты.
	buckets map[antibrut.LimitationCode][]*antibrut.Bucket

	// lastAttemptID последний идентификатор попыток.
	lastAttemptID antibrut.AttemptID

	// attemptsMu мьютекс по попыткам.
	attemptsMu sync.Mutex

	// attempts содержит в себе попытки.
	attempts map[antibrut.BucketID][]*antibrut.Attempt
}

// New создает репозиторий.
func New() *Repository {
	return &Repository{
		buckets:  make(map[antibrut.LimitationCode][]*antibrut.Bucket, 0),
		attempts: make(map[antibrut.BucketID][]*antibrut.Attempt, 0),
	}
}

// FindBucket находит antibrut.Bucket.
// Если совпадений нет, вернет antibrut.ErrNotFound.
func (r *Repository) FindBucket(
	_ context.Context,
	c antibrut.LimitationCode,
	val string,
) (*antibrut.Bucket, error) {
	r.bucketsMu.Lock()
	defer r.bucketsMu.Unlock()

	buckets, ok := r.buckets[c]
	if !ok {
		return nil, antibrut.ErrNotFound
	}

	for _, b := range buckets {
		if b.Value == val {
			return b, nil
		}
	}

	return nil, antibrut.ErrNotFound
}

// CreateBucket создает antibrut.Bucket.
func (r *Repository) CreateBucket(_ context.Context, bucket *antibrut.Bucket) (*antibrut.Bucket, error) {
	r.bucketsMu.Lock()
	defer r.bucketsMu.Unlock()

	bucket.ID = r.lastBucketID + 1
	bucket.CreatedAt = clock.Now()

	r.buckets[bucket.LimitationCode] = append(r.buckets[bucket.LimitationCode], bucket)

	r.lastBucketID = bucket.ID

	return bucket, nil
}

// DeleteBuckets удаляет нужные antibrut.Bucket.
func (r *Repository) DeleteBuckets(_ context.Context, filter antibrut.BucketFilter) (n int64, err error) {
	r.bucketsMu.Lock()
	defer r.bucketsMu.Unlock()

	var counter int64

	for limitationCode, buckets := range r.buckets {
		if filter.LimitationCode != "" && filter.LimitationCode != limitationCode {
			continue
		}

		for i, bucket := range buckets {
			if filter.Value != "" && filter.Value != bucket.Value {
				continue
			}

			if !filter.CreatedAtTo.IsZero() && bucket.CreatedAt.After(filter.CreatedAtTo) {
				continue
			}

			counter++

			r.deleteAttemptsByBucketID(bucket.ID)

			if len(r.buckets[limitationCode]) == 1 {
				delete(r.buckets, limitationCode)
				break
			}

			r.buckets[limitationCode] = append(r.buckets[limitationCode][:i], r.buckets[limitationCode][i+1:]...)
		}
	}

	return counter, nil
}

// FindAttempts находит совпадающие antibrut.Attempt.
func (r *Repository) FindAttempts(_ context.Context, filter antibrut.AttemptFilter) ([]*antibrut.Attempt, error) {
	r.attemptsMu.Lock()
	defer r.attemptsMu.Unlock()

	result := make([]*antibrut.Attempt, 0)

	for bucketID, attempts := range r.attempts {
		if filter.BucketID > 0 && filter.BucketID != bucketID {
			continue
		}

		for _, attempt := range attempts {
			if !filter.CreatedAtFrom.IsZero() && attempt.CreatedAt.Before(filter.CreatedAtFrom) {
				continue
			}

			if !filter.CreatedAtTo.IsZero() && attempt.CreatedAt.After(filter.CreatedAtTo) {
				continue
			}

			result = append(result, attempt)
		}
	}

	return result, nil
}

// CreateAttempt создает antibrut.Attempt.
func (r *Repository) CreateAttempt(_ context.Context, attempt *antibrut.Attempt) (*antibrut.Attempt, error) {
	r.attemptsMu.Lock()
	defer r.attemptsMu.Unlock()

	attempt.ID = r.lastAttemptID + 1
	attempt.CreatedAt = clock.Now()

	r.attempts[attempt.BucketID] = append(r.attempts[attempt.BucketID], attempt)

	r.lastAttemptID = attempt.ID

	return attempt, nil
}

// deleteAttemptsByBucketID удаляет antibrut.Attempt на основе ID бакета.
func (r *Repository) deleteAttemptsByBucketID(bucketID antibrut.BucketID) {
	r.attemptsMu.Lock()
	defer r.attemptsMu.Unlock()

	delete(r.attempts, bucketID)
}
