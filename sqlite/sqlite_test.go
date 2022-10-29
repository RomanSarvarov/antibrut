package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/clock"
	mock "github.com/romsar/antibrut/mock/sqlite"
)

type mocks struct {
	db *mock.Database
}

func newFakeRepository(t *testing.T) (*Repository, mocks) {
	t.Helper()

	db := mock.NewDatabase(t)
	repo, err := New(":memory:")
	require.NoError(t, err)

	repo.db = db

	return repo, mocks{
		db: db,
	}
}

func setupRepository(t *testing.T) (*Repository, *sql.DB) {
	t.Helper()

	repo, err := New(":memory:")
	require.NoError(t, err)

	err = repo.Migrate()
	require.NoError(t, err)

	db := repo.db.(*sql.DB)

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	return repo, db
}

func TestRepository_Close(t *testing.T) {
	repo, m := newFakeRepository(t)

	m.db.On("Close").Return(nil).Once()

	require.NoError(t, repo.Close())
}

func TestRepository_FindLimitation(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		repo, _ := setupRepository(t)

		limit, err := repo.FindLimitation(ctx, "foobar")
		require.Nil(t, limit)
		require.ErrorIs(t, err, antibrut.ErrNotFound)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		wantLimit := &antibrut.Limitation{
			Code:        "foo",
			MaxAttempts: 10,
			Interval:    clock.NewDurationFromTimeDuration(1 * time.Minute),
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60), ('bar', 10, 60);
		`)
		require.NoError(t, err)

		gotLimit, err := repo.FindLimitation(ctx, "foo")
		require.Equal(t, wantLimit, gotLimit)
		require.NoError(t, err)
	})
}

func TestRepository_FindBucket(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		repo, _ := setupRepository(t)

		limit, err := repo.FindBucket(ctx, "foo", "bar")
		require.Nil(t, limit)
		require.ErrorIs(t, err, antibrut.ErrNotFound)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		tm := clock.NewFromTime(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))

		wantBucket := &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "bar",
			CreatedAt:      tm,
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(limitation_code, value, created_at) 
			VALUES('foo', 'bar', ?), ('foo', 'abc', 0);
		`, tm)
		require.NoError(t, err)

		gotBucket, err := repo.FindBucket(ctx, "foo", "bar")
		require.NoError(t, err)

		require.NotNil(t, gotBucket.ID)
		require.Equal(t, wantBucket.LimitationCode, gotBucket.LimitationCode)
		require.Equal(t, wantBucket.Value, gotBucket.Value)
		require.Equal(t, wantBucket.CreatedAt, gotBucket.CreatedAt)
	})
}

func TestRepository_CreateBucket(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		tm := clock.NewFromTime(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))

		wantBucket := &antibrut.Bucket{
			LimitationCode: "foo",
			Value:          "bar",
			CreatedAt:      tm,
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
		`, tm)
		require.NoError(t, err)

		gotBucket, err := repo.CreateBucket(ctx, wantBucket)
		require.NoError(t, err)

		require.NotNil(t, gotBucket.ID)
		require.Equal(t, wantBucket.LimitationCode, gotBucket.LimitationCode)
		require.Equal(t, wantBucket.Value, gotBucket.Value)
		require.Equal(t, wantBucket.CreatedAt, gotBucket.CreatedAt)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM buckets WHERE id = ?`, gotBucket.ID)
		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 1, count)
	})
}

func TestRepository_DeleteBuckets(t *testing.T) {
	t.Parallel()

	t.Run("delete by limitation code", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60), ('bar', 10, 60);
			INSERT INTO buckets(limitation_code, value, created_at) 
			VALUES ('foo', '10', 0), ('bar', '10', 0);
		`)
		require.NoError(t, err)

		n, err := repo.DeleteBuckets(ctx, antibrut.BucketFilter{
			LimitationCode: "foo",
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM buckets`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		row = db.QueryRow(`SELECT COUNT(*) FROM buckets WHERE limitation_code = 'foo'`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("delete by value", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(limitation_code, value, created_at) 
			VALUES ('foo', '10', 0), ('foo', '20', 0);
		`)
		require.NoError(t, err)

		n, err := repo.DeleteBuckets(ctx, antibrut.BucketFilter{
			Value: "10",
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM buckets`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		row = db.QueryRow(`SELECT COUNT(*) FROM buckets WHERE value = '10'`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("delete by created at to", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(limitation_code, value, created_at) 
			VALUES ('foo', '10', '2022-05-01 12:00:00+00:00'), ('foo', '20', '2022-05-02 12:00:00+00:00');
		`)
		require.NoError(t, err)

		createdAtTo := time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC)

		n, err := repo.DeleteBuckets(ctx, antibrut.BucketFilter{
			CreatedAtTo: clock.NewFromTime(createdAtTo),
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM buckets`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		row = db.QueryRow(`SELECT COUNT(*) FROM buckets WHERE created_at = '2022-05-01 12:00:00+00:00'`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestRepository_FindAttempts(t *testing.T) {
	t.Run("find by bucket id", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(id, limitation_code, value, created_at) 
			VALUES (5, 'foo', '10', 0), (6, 'foo', '10', 0);
			INSERT INTO attempts(bucket_id, created_at) 
			VALUES (5, '2022-05-01 12:00:00+00:00'), (6, '2022-05-01 13:00:00+00:00');
		`)
		require.NoError(t, err)

		attempts, err := repo.FindAttempts(ctx, antibrut.AttemptFilter{
			BucketID: 5,
		})
		require.NoError(t, err)
		require.Len(t, attempts, 1)

		require.Equal(t, antibrut.BucketID(5), attempts[0].BucketID)
		require.Equal(
			t,
			clock.NewFromTime(time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC)),
			attempts[0].CreatedAt,
		)
	})

	t.Run("find by created at from", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(id, limitation_code, value, created_at) 
			VALUES (5, 'foo', '10', 0);
			INSERT INTO attempts(bucket_id, created_at) 
			VALUES (5, '2022-05-01 12:00:00+00:00'), (5, '2022-05-01 13:00:00+00:00');
		`)
		require.NoError(t, err)

		attempts, err := repo.FindAttempts(ctx, antibrut.AttemptFilter{
			CreatedAtFrom: clock.NewFromTime(time.Date(2022, 5, 1, 13, 0, 0, 0, time.UTC)),
		})
		require.NoError(t, err)
		require.Len(t, attempts, 1)

		require.Equal(t, antibrut.BucketID(5), attempts[0].BucketID)
		require.Equal(
			t,
			clock.NewFromTime(time.Date(2022, 5, 1, 13, 0, 0, 0, time.UTC)),
			attempts[0].CreatedAt,
		)
	})

	t.Run("find by created at to", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(id, limitation_code, value, created_at) 
			VALUES (5, 'foo', '10', 0);
			INSERT INTO attempts(bucket_id, created_at) 
			VALUES (5, '2022-05-01 12:00:00+00:00'), (5, '2022-05-01 13:00:00+00:00');
		`)
		require.NoError(t, err)

		attempts, err := repo.FindAttempts(ctx, antibrut.AttemptFilter{
			CreatedAtTo: clock.NewFromTime(time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC)),
		})
		require.NoError(t, err)
		require.Len(t, attempts, 1)

		require.Equal(t, antibrut.BucketID(5), attempts[0].BucketID)
		require.Equal(
			t,
			clock.NewFromTime(time.Date(2022, 5, 1, 12, 0, 0, 0, time.UTC)),
			attempts[0].CreatedAt,
		)
	})
}

func TestRepository_CreateAttempt(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		tm := clock.NewFromTime(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))

		wantAttempt := &antibrut.Attempt{
			BucketID:  5,
			CreatedAt: tm,
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO limitations(code, max_attempts, interval_sec) 
			VALUES('foo', 10, 60);
			INSERT INTO buckets(id, limitation_code, value, created_at) 
			VALUES(5, 'foo', 'bar', 0);
		`, tm)
		require.NoError(t, err)

		gotAttempt, err := repo.CreateAttempt(ctx, wantAttempt)
		require.NoError(t, err)

		require.NotNil(t, gotAttempt.ID)
		require.Equal(t, wantAttempt.BucketID, gotAttempt.BucketID)
		require.Equal(t, wantAttempt.CreatedAt, gotAttempt.CreatedAt)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM attempts WHERE id = ?`, gotAttempt.ID)
		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 1, count)
	})
}

func TestRepository_FindIPRuleBySubnet(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		repo, _ := setupRepository(t)

		ipRule, err := repo.FindIPRuleBySubnet(ctx, "192.168.5.0/26")
		require.Nil(t, ipRule)
		require.ErrorIs(t, err, antibrut.ErrNotFound)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		wantIPRule := &antibrut.IPRule{
			Type:   10,
			Subnet: "192.168.5.0/26",
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO ip_rules(type, subnet) 
			VALUES(10, '192.168.5.0/26'), (20, '192.168.5.0/24');
		`)
		require.NoError(t, err)

		gotIPRule, err := repo.FindIPRuleBySubnet(ctx, "192.168.5.0/26")
		require.NoError(t, err)
		require.NotNil(t, gotIPRule.ID)
		require.Equal(t, wantIPRule.Type, gotIPRule.Type)
		require.Equal(t, wantIPRule.Subnet, gotIPRule.Subnet)
	})
}

func TestRepository_FindIPRulesByIP(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		repo, _ := setupRepository(t)

		ipRule, err := repo.FindIPRulesByIP(ctx, "192.168.5.15")
		require.NoError(t, err)
		require.Len(t, ipRule, 0)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		wantIPRule := &antibrut.IPRule{
			Type:   10,
			Subnet: "192.168.5.0/26",
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO ip_rules(type, subnet) 
			VALUES(10, '192.168.5.0/26'), (30, '192.168.6.0/26');
		`)
		require.NoError(t, err)

		gotIPRules, err := repo.FindIPRulesByIP(ctx, "192.168.5.15")
		require.NoError(t, err)
		require.Len(t, gotIPRules, 1)
		require.Equal(t, wantIPRule.Type, gotIPRules[0].Type)
		require.Equal(t, wantIPRule.Subnet, gotIPRules[0].Subnet)
	})
}

func TestRepository_CreateIPRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		wantIPRule := &antibrut.IPRule{
			Type:   10,
			Subnet: "192.168.5.0/26",
		}

		gotIPRule, err := repo.CreateIPRule(ctx, wantIPRule)
		require.NoError(t, err)

		require.NotNil(t, gotIPRule.ID)
		require.Equal(t, wantIPRule.Type, gotIPRule.Type)
		require.Equal(t, wantIPRule.Subnet, gotIPRule.Subnet)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM ip_rules WHERE id = ?`, gotIPRule.ID)
		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 1, count)
	})
}

func TestRepository_UpdateIPRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		existIPRule := &antibrut.IPRule{
			ID:     5,
			Type:   10,
			Subnet: "192.168.5.0/26",
		}

		wantIPRule := &antibrut.IPRule{
			Type:   20,
			Subnet: "192.168.10.0/26",
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO ip_rules(id, type, subnet) 
			VALUES(?, ?, ?);
		`, existIPRule.ID, existIPRule.Type, existIPRule.Subnet)
		require.NoError(t, err)

		gotIPRule, err := repo.UpdateIPRule(ctx, existIPRule.ID, wantIPRule)
		require.NoError(t, err)

		require.Equal(t, existIPRule.ID, gotIPRule.ID)
		require.Equal(t, wantIPRule.Type, gotIPRule.Type)
		require.Equal(t, wantIPRule.Subnet, gotIPRule.Subnet)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM ip_rules WHERE id = ?`, gotIPRule.ID)
		err = row.Scan(&count)
		require.NoError(t, err)

		require.Equal(t, 1, count)
	})
}

func TestRepository_DeleteIPRules(t *testing.T) {
	t.Parallel()

	t.Run("delete by type", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO ip_rules(type, subnet) 
			VALUES (10, '192.168.5.0/24'), (20, '192.168.10.0/24');
		`)
		require.NoError(t, err)

		n, err := repo.DeleteIPRules(ctx, antibrut.IPRuleFilter{
			Type: 10,
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM ip_rules`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		row = db.QueryRow(`SELECT COUNT(*) FROM ip_rules WHERE type = 10`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("delete by subnet", func(t *testing.T) {
		ctx := context.Background()

		repo, db := setupRepository(t)

		_, err := db.ExecContext(ctx, `
			INSERT INTO ip_rules(type, subnet) 
			VALUES (10, '192.168.5.0/24'), (20, '192.168.10.0/24');
		`)
		require.NoError(t, err)

		n, err := repo.DeleteIPRules(ctx, antibrut.IPRuleFilter{
			Subnet: "192.168.5.0/24",
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), n)

		var count int
		row := db.QueryRow(`SELECT COUNT(*) FROM ip_rules`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		row = db.QueryRow(`SELECT COUNT(*) FROM ip_rules WHERE subnet = '192.168.5.0/24'`)
		err = row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
