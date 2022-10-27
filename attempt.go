package antibrut

import "github.com/romsar/antibrut/clock"

// AttemptID это идентификатор попыток запросов.
type AttemptID int64

// Attempt это попытка запроса.
type Attempt struct {
	// ID это идентификатор попыток запросов.
	ID AttemptID

	// BucketID это идентификатор бакета,
	// к которому принадлежит попытка запроса.
	BucketID BucketID

	// CreatedAt это дата создания.
	CreatedAt clock.Time
}

// AttemptFilter это структура для фильтрации
// выборки попыток запроса.
type AttemptFilter struct {
	// BucketID это идентификатор бакета,
	// к которому принадлежит попытка запроса.
	BucketID BucketID

	// CreatedAtFrom это минимальная дата создания.
	CreatedAtFrom clock.Time

	// CreatedAtTo это максимальная дата создания.
	CreatedAtTo clock.Time
}
