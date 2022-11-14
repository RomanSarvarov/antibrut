package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // драйвер для работы пакеты database/sql с sqlite3.
	"github.com/pressly/goose/v3"
	"github.com/romsar/antibrut"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Repository предоставляет API для работы с хранилищем.
type Repository struct {
	db database

	// timeNow содержит функцию, которая возвращает текущее время.
	timeNow func() time.Time
}

// database декларирует методы для работы с БД.
type database interface {
	io.Closer

	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// Option возвращает функцию, модифицирующую Repository.
type Option func(r *Repository)

// WithTimeNow возвращает функцию, устанавливающую
// callback для получения текущего времени.
func WithTimeNow(f func() time.Time) Option {
	return func(r *Repository) {
		r.timeNow = f
	}
}

// New создает и открывает подключение к БД.
func New(dsn string, opts ...Option) (*Repository, error) {
	if dsn == "" {
		return nil, errors.New("dsn required")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot create database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return nil, fmt.Errorf("enable wal error: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, fmt.Errorf("foreign keys pragma error: %w", err)
	}

	r := &Repository{
		db: db,
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.timeNow == nil {
		r.timeNow = time.Now
	}

	return r, nil
}

// Close закрывает подключение к БД.
func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// migrateMu мьютекс для миграций.
// Goose использует глобальные объекты,
// по-этому в тестах возникает race-condition.
var migrateMu sync.Mutex

// Migrate запускает миграции.
func (r *Repository) Migrate() error {
	migrateMu.Lock()
	defer migrateMu.Unlock()

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(goose.NopLogger())

	db := r.db.(*sql.DB)

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}

	return nil
}

// FindLimitation находит antibrut.Limitation.
// Если совпадений нет, вернет antibrut.ErrNotFound.
func (r *Repository) FindLimitation(ctx context.Context, c antibrut.LimitationCode) (*antibrut.Limitation, error) {
	var limitation antibrut.Limitation

	err := r.db.QueryRowContext(ctx, `
		SELECT code, max_attempts, interval_sec
		FROM limitations
		WHERE code = ?
	`, c).Scan(&limitation.Code, &limitation.MaxAttempts, &limitation.Interval)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = antibrut.ErrNotFound
		}

		return nil, fmt.Errorf("find limitation error: %w", err)
	}

	return &limitation, nil
}

// FindBucket находит antibrut.Bucket.
// Если совпадений нет, вернет antibrut.ErrNotFound.
func (r *Repository) FindBucket(
	ctx context.Context,
	c antibrut.LimitationCode,
	val string,
) (*antibrut.Bucket, error) {
	var bucket antibrut.Bucket

	err := r.db.QueryRowContext(ctx, `
		SELECT id, limitation_code, value, created_at
		FROM buckets
		WHERE limitation_code = ? AND value = ?
	`, c, val).Scan(&bucket.ID, &bucket.LimitationCode, &bucket.Value, &bucket.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = antibrut.ErrNotFound
		}

		return nil, fmt.Errorf("find bucket error: %w", err)
	}

	return &bucket, nil
}

// CreateBucket создает antibrut.Bucket.
func (r *Repository) CreateBucket(ctx context.Context, bucket *antibrut.Bucket) (*antibrut.Bucket, error) {
	bucket.CreatedAt = r.timeNow()

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO buckets (limitation_code, value, created_at)
		VALUES (?, ?, ?)
	`,
		bucket.LimitationCode,
		bucket.Value,
		bucket.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create bucket error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create bucket error: %w", err)
	}

	bucket.ID = antibrut.BucketID(id)

	return bucket, nil
}

// DeleteBuckets удаляет нужные antibrut.Bucket.
func (r *Repository) DeleteBuckets(ctx context.Context, filter antibrut.BucketFilter) (n int64, err error) {
	where, args := []string{"1 = 1"}, []any{}

	if filter.LimitationCode != "" {
		where, args = append(where, "limitation_code = ?"), append(args, filter.LimitationCode)
	}

	if filter.Value != "" {
		where, args = append(where, "value = ?"), append(args, filter.Value)
	}

	if !filter.CreatedAtTo.IsZero() {
		where, args = append(where, "created_at <= ?"), append(args, filter.CreatedAtTo)
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE
		FROM buckets
		WHERE `+strings.Join(where, " AND ")+`
		`,
		args...,
	)
	if err != nil {
		return 0, fmt.Errorf("delete buckets error: %w", err)
	}

	deletedCnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete buckets error: %w", err)
	}

	return deletedCnt, nil
}

// FindAttempts находит совпадающие antibrut.Attempt.
func (r *Repository) FindAttempts(ctx context.Context, filter antibrut.AttemptFilter) ([]*antibrut.Attempt, error) {
	where, args := []string{"1 = 1"}, []any{}

	if filter.BucketID != 0 {
		where, args = append(where, "bucket_id = ?"), append(args, filter.BucketID)
	}

	if !filter.CreatedAtFrom.IsZero() {
		where, args = append(where, "created_at >= ?"), append(args, filter.CreatedAtFrom)
	}

	if !filter.CreatedAtTo.IsZero() {
		where, args = append(where, "created_at <= ?"), append(args, filter.CreatedAtTo)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT 
		    id,
		    bucket_id,
		    created_at
		FROM attempts
		WHERE `+strings.Join(where, " AND ")+`
		`,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("find attempts error: %w", err)
	}
	defer rows.Close()

	attempts := make([]*antibrut.Attempt, 0)
	for rows.Next() {
		var attempt antibrut.Attempt
		if err := rows.Scan(
			&attempt.ID,
			&attempt.BucketID,
			&attempt.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("find attempts error: %w", err)
		}
		attempts = append(attempts, &attempt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find attempts error: %w", err)
	}

	return attempts, nil
}

// CreateAttempt создает antibrut.Attempt.
func (r *Repository) CreateAttempt(ctx context.Context, attempt *antibrut.Attempt) (*antibrut.Attempt, error) {
	attempt.CreatedAt = r.timeNow()

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO attempts (bucket_id, created_at)
		VALUES (?, ?)
	`,
		attempt.BucketID,
		attempt.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create attempt error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create attempt error: %w", err)
	}

	attempt.ID = antibrut.AttemptID(id)

	return attempt, nil
}

// FindIPRuleBySubnet находит antibrut.IPRule на основе подсети.
// Если совпадений нет, вернет antibrut.ErrNotFound.
func (r *Repository) FindIPRuleBySubnet(ctx context.Context, subnet antibrut.Subnet) (*antibrut.IPRule, error) {
	var rule antibrut.IPRule

	err := r.db.QueryRowContext(ctx, `
		SELECT id, type, subnet
		FROM ip_rules
		WHERE subnet = ?
	`, subnet).Scan(&rule.ID, &rule.Type, &rule.Subnet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = antibrut.ErrNotFound
		}

		return nil, fmt.Errorf("find ip rule error: %w", err)
	}

	return &rule, nil
}

// FindIPRulesByIP находит совпадения antibrut.IPRule на основе IP адреса.
func (r *Repository) FindIPRulesByIP(ctx context.Context, ip antibrut.IP) ([]*antibrut.IPRule, error) {
	ipParts := strings.Split(ip.String(), ".")
	ipWithoutLastOctet := strings.Join(ipParts[0:3], ".") + ".%"

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, type, subnet
		FROM ip_rules
		WHERE subnet LIKE ?
	`, ipWithoutLastOctet)
	if err != nil {
		return nil, fmt.Errorf("find ip rules by ip error: %w", err)
	}
	defer rows.Close()

	rules := make([]*antibrut.IPRule, 0)
	for rows.Next() {
		var rule antibrut.IPRule
		if err := rows.Scan(
			&rule.ID,
			&rule.Type,
			&rule.Subnet,
		); err != nil {
			return nil, fmt.Errorf("find ip rules by ip error: %w", err)
		}
		rules = append(rules, &rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find ip rules by ip error: %w", err)
	}

	return rules, nil
}

// CreateIPRule создает antibrut.IPRule.
func (r *Repository) CreateIPRule(ctx context.Context, ipRule *antibrut.IPRule) (*antibrut.IPRule, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO ip_rules (type, subnet)
		VALUES (?, ?)
	`,
		ipRule.Type,
		ipRule.Subnet,
	)
	if err != nil {
		return nil, fmt.Errorf("create ip rule error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create ip rule error: %w", err)
	}

	ipRule.ID = antibrut.IPRuleID(id)

	return ipRule, nil
}

// UpdateIPRule обновляет antibrut.IPRule.
func (r *Repository) UpdateIPRule(
	ctx context.Context,
	id antibrut.IPRuleID,
	upd *antibrut.IPRuleUpdate,
) (*antibrut.IPRule, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ip_rules 
		SET type = ?, subnet = ?
		WHERE id = ?
	`,
		upd.Type,
		upd.Subnet,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("update ip rule error: %w", err)
	}

	return &antibrut.IPRule{
		ID:     id,
		Type:   upd.Type,
		Subnet: upd.Subnet,
	}, nil
}

// DeleteIPRules удаляет совпадения из antibrut.IPRule.
func (r *Repository) DeleteIPRules(ctx context.Context, filter antibrut.IPRuleFilter) (int64, error) {
	where, args := []string{"1 = 1"}, []any{}

	if filter.Type != 0 {
		where, args = append(where, "type = ?"), append(args, filter.Type)
	}

	if filter.Subnet != "" {
		where, args = append(where, "subnet = ?"), append(args, filter.Subnet)
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE
		FROM ip_rules
		WHERE `+strings.Join(where, " AND ")+`
		`,
		args...,
	)
	if err != nil {
		return 0, fmt.Errorf("delete ip rules by subnet error: %w", err)
	}

	deletedCnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete ip rules by subnet error: %w", err)
	}

	return deletedCnt, nil
}
