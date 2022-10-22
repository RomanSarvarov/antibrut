package leakybucket

type Service struct {
	s storage
}

func New(s storage) *Service {
	return &Service{
		s: s,
	}
}

type storage interface {
	/*FindLimit(c LimitCode) (*Limit, error)
	FindOrCreateBucket(lc LimitCode, val string) (*Bucket, error)*/
}

type LimitCode string

type Limit struct {
	Code      LimitCode `db:"code"`
	Attempts  int       `db:"attempts"`
	PerSecond int       `db:"per_second"`
}

type Bucket struct {
	Max       uint
	LimitCode LimitCode
	Value     string
	Count     int
}
