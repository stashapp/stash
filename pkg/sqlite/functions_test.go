package sqlite

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomFunctions(t *testing.T) {
	db, err := sql.Open(sqlite3Driver, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test regexp functions
	t.Run("regexp", func(t *testing.T) {
		var result bool
		err := db.QueryRow(`SELECT regexp('foo.*', 'seafood')`).Scan(&result)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("regexp_substr", func(t *testing.T) {
		var result string
		err := db.QueryRow(`SELECT regexp_substr('seafood', 'foo.')`).Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, "food", result)
	})

	t.Run("regexp_capture", func(t *testing.T) {
		var result string
		err := db.QueryRow(`SELECT regexp_capture('seafood', '(foo)(d)', 2)`).Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, "d", result)
	})

	t.Run("regexp_replace", func(t *testing.T) {
		var result string
		err := db.QueryRow(`SELECT regexp_replace('seafood', 'foo', 'bar')`).Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, "seabard", result)
	})

	// Test initcap
	t.Run("initcap", func(t *testing.T) {
		var result string
		err := db.QueryRow(`SELECT initcap('hello world')`).Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})

	// Test trim
	t.Run("trim", func(t *testing.T) {
		var result string
		err := db.QueryRow(`SELECT trim('  hello world  ')`).Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", result)
	})

	// Test uuid4
	t.Run("uuid4", func(t *testing.T) {
		var result string
		err := db.QueryRow(`SELECT uuid4()`).Scan(&result)
		assert.NoError(t, err)
		// check if it's a valid UUID
		_, err = regexp.Compile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
		assert.NoError(t, err)
	})
}
