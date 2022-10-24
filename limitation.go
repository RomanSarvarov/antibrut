package antibrut

import "github.com/romsar/antibrut/clock"

type LimitationCode string

const (
	LoginLimitation    LimitationCode = "login"
	PasswordLimitation LimitationCode = "password"
	IPLimitation       LimitationCode = "ip"
)

type Limitation struct {
	Code        LimitationCode
	MaxAttempts int
	Interval    clock.Duration
}
