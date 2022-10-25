package antibrut

import "github.com/pkg/errors"

// ErrNotFound это ошибка, которая говорит о том,
// что искомый объект не найден.
var ErrNotFound = errors.New("not found")
