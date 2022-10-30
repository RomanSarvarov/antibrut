package sqlite

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/romsar/antibrut/mock/sqlite"
)

type mocks struct {
	db *mock.Database
}

func newFakeRepository(t *testing.T) (*Repository, mocks) {
	t.Helper()

	db := mock.NewDatabase(t)
	repo, err := New(":memory:?_foreign_keys=on")
	require.NoError(t, err)

	repo.db = db

	return repo, mocks{
		db: db,
	}
}

func TestRepository_Close(t *testing.T) {
	repo, m := newFakeRepository(t)

	m.db.On("Close").Return(nil).Once()

	require.NoError(t, repo.Close())
}
