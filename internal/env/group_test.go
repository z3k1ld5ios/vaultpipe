package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroup_NilInput_ReturnsError(t *testing.T) {
	g := NewGroup(nil)
	_, err := g.Apply(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestGroup_NoPrefix_AllInDefault(t *testing.T) {
	g := NewGroup(map[string]string{})
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := g.Apply(env)
	require.NoError(t, err)
	assert.Len(t, res.Groups["default"], 2)
}

func TestGroup_MatchesPrefixedKeys(t *testing.T) {
	g := NewGroup(map[string]string{"DB_": "database"})
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "prod",
	}
	res, err := g.Apply(env)
	require.NoError(t, err)
	assert.Len(t, res.Groups["database"], 2)
	assert.Contains(t, res.Groups["database"], "DB_HOST")
	assert.Contains(t, res.Groups["database"], "DB_PORT")
	assert.Len(t, res.Groups["default"], 1)
	assert.Contains(t, res.Groups["default"], "APP_ENV")
}

func TestGroup_LongestPrefixWins(t *testing.T) {
	g := NewGroup(map[string]string{
		"CACHE_":      "cache",
		"CACHE_REDIS_": "redis",
	})
	env := map[string]string{
		"CACHE_REDIS_HOST": "redis.local",
		"CACHE_MEMCACHED":  "mem.local",
	}
	res, err := g.Apply(env)
	require.NoError(t, err)
	assert.Contains(t, res.Groups["redis"], "CACHE_REDIS_HOST")
	assert.Contains(t, res.Groups["cache"], "CACHE_MEMCACHED")
}

func TestGroup_EmptyEnv_ReturnsEmptyResult(t *testing.T) {
	g := NewGroup(map[string]string{"DB_": "database"})
	res, err := g.Apply(map[string]string{})
	require.NoError(t, err)
	assert.Empty(t, res.Groups)
}

func TestGroupResult_Keys_SortedAlphabetically(t *testing.T) {
	g := NewGroup(map[string]string{"X_": "x"})
	env := map[string]string{"X_C": "3", "X_A": "1", "X_B": "2"}
	res, err := g.Apply(env)
	require.NoError(t, err)
	keys := res.Keys("x")
	assert.Equal(t, []string{"X_A", "X_B", "X_C"}, keys)
}

func TestGroupResult_Keys_MissingGroup_ReturnsNil(t *testing.T) {
	g := NewGroup(map[string]string{})
	res, err := g.Apply(map[string]string{"FOO": "bar"})
	require.NoError(t, err)
	assert.Nil(t, res.Keys("nonexistent"))
}

func TestGroup_MultipleGroups_AllPopulated(t *testing.T) {
	g := NewGroup(map[string]string{
		"DB_":    "database",
		"CACHE_": "cache",
		"AUTH_":  "auth",
	})
	env := map[string]string{
		"DB_URL":      "postgres://",
		"CACHE_TTL":   "300",
		"AUTH_SECRET": "s3cr3t",
		"APP_NAME":    "vaultpipe",
	}
	res, err := g.Apply(env)
	require.NoError(t, err)
	assert.Len(t, res.Groups, 4)
	assert.Contains(t, res.Groups["database"], "DB_URL")
	assert.Contains(t, res.Groups["cache"], "CACHE_TTL")
	assert.Contains(t, res.Groups["auth"], "AUTH_SECRET")
	assert.Contains(t, res.Groups["default"], "APP_NAME")
}
