package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Repository struct {
	db database
}

type database struct {
	*sql.DB
}

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

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, errors.Wrap(err, "foreign keys pragma error")
	}

	s := &Repository{
		db: database{DB: db},
	}

	return s, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Migrate() error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return errors.Wrap(err, "migrate error")
	}

	goose.SetBaseFS(migrationsFS)

	if err := goose.Up(r.db.DB, "migrations"); err != nil {
		return errors.Wrap(err, "migrate error")
	}

	return nil
}

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

func (r *Repository) DeleteIPRuleBySubnet(ctx context.Context, subnet antibrut.Subnet) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE 
		FROM ip_rules
		WHERE subnet = ?
	`,
		subnet,
	)
	if err != nil {
		return 0, errors.Wrap(err, "delete ip rules by subnet error")
	}

	deletedCnt, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "delete ip rules by subnet error")
	}

	return deletedCnt, nil

	return result.RowsAffected()
}

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

func (r *Repository) DeleteBuckets(ctx context.Context, filter antibrut.BucketFilter) (n int64, err error) {
	where, args := []string{"1 = 1"}, []any{}

	if filter.LimitationCode != "" {
		where, args = append(where, "limitation_code = ?"), append(args, filter.LimitationCode)
	}

	if filter.Value != "" {
		where, args = append(where, "value = ?"), append(args, filter.Value)
	}

	if !filter.DateTo.IsZero() {
		where, args = append(where, "created_at <= ?"), append(args, filter.DateTo)
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
