package antibrut

import (
	"time"
)

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
	CreatedAt time.Time
}

// AttemptFilter это структура для фильтрации
// выборки попыток запроса.
type AttemptFilter struct {
	// BucketID это идентификатор бакета,
	// к которому принадлежит попытка запроса.
	BucketID BucketID

	// CreatedAtFrom это минимальная дата создания.
	CreatedAtFrom time.Time

	// CreatedAtTo это максимальная дата создания.
	CreatedAtTo time.Time
}
