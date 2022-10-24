package antibrut

import "github.com/pkg/errors"

// ErrNotFound это ошибка, которая говорит о том,
// что искомый объект не найден.
var ErrNotFound = errors.New("not found")

// ErrMaxAttemptsExceeded это ошибка, которая говорит о том,
// что количество попыток больше максимально-допустимого значения.
var ErrMaxAttemptsExceeded = errors.New("max attempts exceeded")

// ErrIPInBlackList это ошибка, которая говорит о том,
// что данный IP адрес в черном списке.
var ErrIPInBlackList = errors.New("ip address in black list")
