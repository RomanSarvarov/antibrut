package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/romsar/antibrut/clock"
)

// envFilePath путь до .env файла.
const envFilePath = ".env"

// Config это настройки приложения.
type Config struct {
	// GRPC это настройки GRPC.
	GRPC GRPCConfig

	// SQLiteConfig это настройки для БД SQLite.
	SQLite SQLiteConfig

	// RateLimiterStorageDriver драйвер для лимитинга.
	// Возможные значения: sqlite, inmem.
	RateLimiterStorageDriver string `env:"ANTIBRUT_RATE_LIMITER_DRIVER,notEmpty" envDefault:"inmem"`

	// PruneDuration количество времени, после
	// которого бакеты считаются устаревшими и удаляются.
	// Если указать "0", то удаление неактуальных,
	// бакетов производиться не будет.
	PruneDuration clock.Duration `env:"ANTIBRUT_PRUNE_DURATION" envDefault:"1h"`
}

// GRPCConfig это настройки GRPC.
type GRPCConfig struct {
	// Address адрес сервера.
	Address string `env:"ANTIBRUT_GRPC_ADDRESS,notEmpty" envDefault:":9090"`
}

// SQLiteConfig это настройки для БД SQLite.
type SQLiteConfig struct {
	// DSN настройки подключения к БД (путь к файлу).
	// Либо :memory:, тогда будет использовать оперативная память.
	DSN string `env:"ANTIBRUT_SQLITE_DSN,notEmpty" envDefault:"./data/db.sqlite?_foreign_keys=on"`
}

// New создает Config.
func New() *Config {
	return new(Config)
}

// Load возвращает параметры приложения в виде структуры Config.
func Load() (*Config, error) {
	cfg := New()

	// Загрузим переменные среды из .env файла.
	err := godotenv.Load(envFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "load `%s` file error", envFilePath)
	}

	// Заполним структуру env переменными.
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "env parse error")
	}

	return cfg, nil
}
