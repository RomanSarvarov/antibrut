package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"io"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Repository предоставляет API для работы с хранилищем.
type Repository struct {
	db database
}

// database декларирует методы для работы с БД.
type database interface {
	io.Closer

	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// New создает и открывает подключение к БД.
func New(dsn string) (*Repository, error) {
	if dsn == "" {
		return nil, errors.New("dsn required")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database connection")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "cannot ping database")
	}

	if _, err := db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return nil, errors.Wrap(err, "enable wal error")
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, errors.Wrap(err, "foreign keys pragma error")
	}

	s := &Repository{
		db: db,
	}

	return s, nil
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
		return errors.Wrap(err, "migrate error")
	}

	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(goose.NopLogger())

	db := r.db.(*sql.DB)

	if err := goose.Up(db, "migrations"); err != nil {
		return errors.Wrap(err, "migrate error")
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

		return nil, errors.Wrap(err, "find limitation error")
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

		return nil, errors.Wrap(err, "find bucket error")
	}

	return &bucket, nil
}

// CreateBucket создает antibrut.Bucket.
func (r *Repository) CreateBucket(ctx context.Context, bucket *antibrut.Bucket) (*antibrut.Bucket, error) {
	bucket.CreatedAt = clock.Now()

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO buckets (limitation_code, value, created_at)
		VALUES (?, ?, ?)
	`,
		bucket.LimitationCode,
		bucket.Value,
		bucket.CreatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create bucket error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "create bucket error")
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
		return 0, errors.Wrap(err, "delete buckets error")
	}

	deletedCnt, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "delete buckets error")
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
		return nil, errors.Wrap(err, "find attempts error")
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
			return nil, errors.Wrap(err, "find attempts error")
		}
		attempts = append(attempts, &attempt)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "find attempts error")
	}

	return attempts, nil
}

// CreateAttempt создает antibrut.Attempt.
func (r *Repository) CreateAttempt(ctx context.Context, attempt *antibrut.Attempt) (*antibrut.Attempt, error) {
	attempt.CreatedAt = clock.Now()

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO attempts (bucket_id, created_at)
		VALUES (?, ?)
	`,
		attempt.BucketID,
		attempt.CreatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create attempt error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "create attempt error")
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

		return nil, errors.Wrap(err, "find ip rule error")
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
		return nil, errors.Wrap(err, "find ip rules by ip error")
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
			return nil, errors.Wrap(err, "find ip rules by ip error")
		}
		rules = append(rules, &rule)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "find ip rules by ip error")
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
		return nil, errors.Wrap(err, "create ip rule error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "create ip rule error")
	}

	ipRule.ID = antibrut.IPRuleID(id)

	return ipRule, nil
}

// UpdateIPRule обновляет antibrut.IPRule.
func (r *Repository) UpdateIPRule(
	ctx context.Context,
	id antibrut.IPRuleID,
	ipRule *antibrut.IPRule,
) (*antibrut.IPRule, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ip_rules 
		SET type = ?, subnet = ?
		WHERE id = ?
	`,
		ipRule.Type,
		ipRule.Subnet,
		id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "update ip rule error")
	}

	ipRule.ID = id

	return ipRule, nil
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
		return 0, errors.Wrap(err, "delete ip rules by subnet error")
	}

	deletedCnt, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "delete ip rules by subnet error")
	}

	return deletedCnt, nil
}
