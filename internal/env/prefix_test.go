package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply_AddsPrefix(t *testing.T) {
	pf := NewPrefixFilter("APP_")
	result := pf.Apply(map[string]string{"DB_HOST": "localhost", "PORT": "5432"})
	assert.Equal(t, "localhost", result["APP_DB_HOST"])
	assert.Equal(t, "5432", result["APP_PORT"])
	assert.Len(t, result, 2)
}

func TestApply_EmptyPrefix_ReturnsUnchanged(t *testing.T) {
	pf := NewPrefixFilter("")
	input := map[string]string{"KEY": "val"}
	result := pf.Apply(input)
	assert.Equal(t, input, result)
}

func TestStrip_RemovesPrefixedKeys(t *testing.T) {
	pf := NewPrefixFilter("APP_")
	input := map[string]string{"APP_HOST": "localhost", "OTHER": "ignored"}
	result := pf.Strip(input)
	assert.Equal(t, "localhost", result["HOST"])
	assert.NotContains(t, result, "OTHER")
	assert.Len(t, result, 1)
}

func TestStrip_EmptyPrefix_ReturnsAll(t *testing.T) {
	pf := NewPrefixFilter("")
	input := map[string]string{"A": "1", "B": "2"}
	result := pf.Strip(input)
	assert.Equal(t, input, result)
}

func TestHasPrefix_Match(t *testing.T) {
	pf := NewPrefixFilter("VAULT_")
	assert.True(t, pf.HasPrefix("VAULT_SECRET"))
	assert.False(t, pf.HasPrefix("OTHER_KEY"))
}

func TestHasPrefix_EmptyPrefix_AlwaysTrue(t *testing.T) {
	pf := NewPrefixFilter("")
	assert.True(t, pf.HasPrefix("ANYTHING"))
	assert.True(t, pf.HasPrefix(""))
}

func TestApply_PrefixNormalisedToUpper(t *testing.T) {
	pf := NewPrefixFilter("app_")
	result := pf.Apply(map[string]string{"KEY": "v"})
	assert.Contains(t, result, "APP_KEY")
}
