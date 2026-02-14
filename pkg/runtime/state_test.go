package runtime

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		_ = NewForTest()
		require.True(t, true)
	})
}

func TestUUIDFunction(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test_uuid (
		id TEXT NOT NULL DEFAULT (uuid()),
		name TEXT NOT NULL
	)`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO test_uuid (name) VALUES ('alice')`)
	require.NoError(t, err)

	var id string
	err = db.QueryRow(`SELECT id FROM test_uuid WHERE name = 'alice'`).Scan(&id)
	require.NoError(t, err)

	_, err = uuid.Parse(id)
	require.NoError(t, err, "expected a valid UUID, got: %s", id)
}
