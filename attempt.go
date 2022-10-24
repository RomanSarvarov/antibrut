package antibrut

import "github.com/romsar/antibrut/clock"

type AttemptID int64

type Attempt struct {
	ID        AttemptID
	BucketID  BucketID
	CreatedAt clock.Time
}

type AttemptFilter struct {
	BucketID      BucketID
	CreatedAtFrom clock.Time
	CreatedAtTo   clock.Time
}
