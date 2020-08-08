package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	name     string
	input    string
	expected string
}

func TestReadFirstLine(t *testing.T) {
	for _, tst := range [...]test{
		{
			name:     "simple case",
			input:    "val",
			expected: "val",
		},
		{
			name:     "with new line",
			input:    "val\nvar",
			expected: "val",
		},
		{
			name:     "with EOF",
			input:    "",
			expected: "",
		},
		{
			name:     "with tab",
			input:    "val\t\t",
			expected: "val",
		},
		{
			name:     "with space at the end",
			input:    "val  ",
			expected: "val",
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			result, err := readFirstLine(strings.NewReader(tst.input))
			require.NoError(t, err)
			require.Equal(t, tst.expected, result)
		})
	}
}

func TestReplaceZeroByte(t *testing.T) {
	result, err := readFirstLine(bytes.NewReader([]byte{'A', 'B', 0x00, 'C'}))
	require.NoError(t, err)
	require.Equal(t, "AB\nC", result)
}

func TestReadDir(t *testing.T) {
	env, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	expectedEnv := make(Environment)
	expectedEnv["BAR"] = "bar"
	expectedEnv["FOO"] = "   foo\nwith new line"
	expectedEnv["HELLO"] = "\"hello\""
	expectedEnv["UNSET"] = ""

	require.Equal(t, expectedEnv, env)
}

func TestReadDirWrongPath(t *testing.T) {
	env, err := ReadDir("/no_dir")
	require.Error(t, err)

	emptyEnv := make(Environment)
	require.Equal(t, emptyEnv, env)
}

func TestReadDirWithNestedDir(t *testing.T) {
	env, err := ReadDir("./testdata")
	require.NoError(t, err)

	expectedEnv := make(Environment)
	expectedEnv["echo.sh"] = "#!/usr/bin/env bash"

	require.Equal(t, expectedEnv, env)
}
