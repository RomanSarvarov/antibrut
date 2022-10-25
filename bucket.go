package antibrut

import (
	"github.com/romsar/antibrut/clock"
)

type BucketID int64

type Bucket struct {
	ID             BucketID
	LimitationCode LimitationCode
	Value          string
	CreatedAt      clock.Time
}

type BucketFilter struct {
	LimitationCode LimitationCode
	Value          string
	DateTo         clock.Time
}

type ResetFilter = BucketFilter
