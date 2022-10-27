package antibrut

import (
	"github.com/romsar/antibrut/clock"
)

// BucketID это идентификатор бакета.
type BucketID int64

// Bucket это бакет попыток запроса.
// Бакет имеет определенную пропускную способность,
// а также имеет максимальный объем.
// Если количество запросов превышает пропускную,
// то в моменте превышения максимального объема,
// запрос будет отклонён.
type Bucket struct {
	// ID это идентификатор бакета.
	ID BucketID

	// LimitationCode это идентификатор лимита.
	LimitationCode LimitationCode

	// Value это значение бакета.
	Value string

	// CreatedAt это дата создания.
	CreatedAt clock.Time
}

// BucketFilter предоставляет структуру для фильтрации бакетов.
type BucketFilter struct {
	// LimitationCode это идентификатор лимита.
	LimitationCode LimitationCode

	// Value это значение бакета.
	Value string

	// DateTo это максимальная дата создания.
	DateTo clock.Time
}

// ResetFilter предоставляет структуру для сброса бакетов.
type ResetFilter = BucketFilter
