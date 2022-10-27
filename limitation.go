package antibrut

import "github.com/romsar/antibrut/clock"

// LimitationCode идентификатор ограничения.
type LimitationCode string

const (
	// LoginLimitation ограничение по логину пользователя.
	LoginLimitation LimitationCode = "login"

	// PasswordLimitation ограничение по паролю пользователя.
	PasswordLimitation LimitationCode = "password"

	// IPLimitation ограничение по IP пользователя.
	IPLimitation LimitationCode = "ip"
)

// Limitation это настройка ограничения по определенному признаку.
type Limitation struct {
	// Code это уникальный текстовый идентификатор.
	Code LimitationCode

	// MaxAttempts задает максимальное количество попыток,
	// при превышении которых будет отдаваться ошибка.
	MaxAttempts int

	// Interval промежуток времени, в течении которого считать MaxAttempts.
	Interval clock.Duration
}
