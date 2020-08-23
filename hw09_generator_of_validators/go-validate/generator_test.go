package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSupportedType(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		m := make(map[string]string)
		val, ok := getSupportedType("string", m)
		require.True(t, ok)
		require.Equal(t, "string", val)
	})
	t.Run("custom type", func(t *testing.T) {
		m := make(map[string]string)
		ct := "MyCustomType"
		m[ct] = "int"
		val, ok := getSupportedType(ct, m)
		require.True(t, ok)
		require.Equal(t, "int", val)
	})
	t.Run("unknown type", func(t *testing.T) {
		m := make(map[string]string)
		_, ok := getSupportedType("someType", m)
		require.False(t, ok)			
	})
	t.Run("hard way", func(t *testing.T) {
		m := make(map[string]string)
		m["MyCustomType"] = "int"
		m["SomeOtherType"] = "XType"
		_, ok := getSupportedType("SomeOtherType", m)
		require.False(t, ok)	
	})
}

func TestParseTag(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		tags, err := parseTag(`validate:"min:18"`)
		expect := make(map[string]string)
		expect["min"] = "18"

		require.NoError(t, err)
		require.Equal(t, expect, tags)
	})
	t.Run("multiple expr", func(t *testing.T) {
		tags, err := parseTag(`validate:"min:18|in:20,30"`)
		expect := make(map[string]string)
		expect["min"] = "18"
		expect["in"] = "20,30"

		require.NoError(t, err)
		require.Equal(t, expect, tags)
	})
	t.Run("no validate tag", func(t *testing.T) {
		tags, err := parseTag(`db:"pk"`)
		expect := make(map[string]string)

		require.NoError(t, err)
		require.Equal(t, expect, tags)
	})
	t.Run("wrong expr", func(t *testing.T) {
		tags, err := parseTag(`validate:"some=18"`)
		require.Error(t, err)
		require.Nil(t, tags)
	})		
}
