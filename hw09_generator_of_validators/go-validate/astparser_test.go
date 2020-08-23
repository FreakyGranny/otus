package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSimple(t *testing.T) {
	tSource :=
		`package xxx

		type Req struct {
			ID   int    #validate:"len:5"#
			Body string
		}
		`
	vSource, err := parseSource(strings.NewReader(strings.Replace(tSource, "#", "`", -1)))

	require.NoError(t, err)
	require.Equal(t, vSource.Package, "xxx")
	require.Equal(t, len(vSource.Structs), 1)
}

func TestParseNoTags(t *testing.T) {
	tSource :=
		`package xxx

		type Req struct {
			ID   int
			Body string
		}
		`
	vSource, err := parseSource(strings.NewReader(tSource))

	require.NoError(t, err)
	require.Equal(t, vSource.Package, "xxx")
	require.Equal(t, len(vSource.Structs), 0)
}

func TestParseBadFile(t *testing.T) {
	tSource :=
		`package xxx

		type Req struct {
			ID   int    syntax error
			Body string
		}
		`
	vSource, err := parseSource(strings.NewReader(tSource))

	require.Error(t, err)
	require.Nil(t, vSource)
}

func TestParseWithCustomTypes(t *testing.T) {
	tSource :=
		`package xxx

		type myInt int
		`
	vSource, err := parseSource(strings.NewReader(tSource))

	require.NoError(t, err)
	x, ok := vSource.Types["myInt"]
	require.True(t, ok)
	require.Equal(t, x, "int")
}

func TestParseArray(t *testing.T) {
	tSource :=
		`package xxx

		type Req struct {
			ID   []int      #validate:"len:5"#
			Body [][]string #validate:"len:5"#
		}
		`
	vSource, err := parseSource(strings.NewReader(strings.Replace(tSource, "#", "`", -1)))

	require.NoError(t, err)
	require.Equal(t, vSource.Package, "xxx")
	require.Equal(t, len(vSource.Structs), 1)
	
	field := vSource.Structs["Req"][0]
	tooHardField := vSource.Structs["Req"][1]
	require.True(t, field.IsArray)
	require.Equal(t, field.Type, "int")
	require.True(t, tooHardField.IsArray)
	require.Equal(t, tooHardField.Type, "")

}
