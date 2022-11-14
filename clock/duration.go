package clock

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Duration это промежуток времени.
// По сути является оберткой над time.Duration.
type Duration struct {
	time.Duration
}

// NewDurationFromTimeDuration создает Duration на основе time.Duration.
func NewDurationFromTimeDuration(d time.Duration) Duration {
	return Duration{
		Duration: d,
	}
}

// ToDuration переводит Duration в time.Duration.
func (d Duration) ToDuration() time.Duration {
	return d.Duration
}

// Value возвращает данные в формате для сохранения в БД.
func (d Duration) Value() (driver.Value, error) {
	return driver.Value(int64(d.Seconds())), nil
}

// Scan определяет способ сканирования данных из БД.
func (d *Duration) Scan(raw any) error {
	switch v := raw.(type) {
	case int64:
		*d = Duration{
			Duration: time.Duration(v) * time.Second,
		}
	case nil:
		*d = Duration{
			Duration: time.Duration(0),
		}
	default:
		return errors.New("cannot sql.Scan() antibrut.Duration")
	}
	return nil
}

// UnmarshalText содержит логику преобразования из текста.
func (d *Duration) UnmarshalText(text []byte) error {
	dd, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration{Duration: dd}
	return nil
}
