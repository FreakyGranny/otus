package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"

	"strconv"
	"strings"
	"unicode"
)

const escapeRune = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(line string) (string, error) {
	var result strings.Builder
	var counter int
	var currentRune rune
	runeSlice := []rune(line)
	rLen := len(runeSlice)

	for i := 0; i < rLen; i++ {
		currentRune = runeSlice[i]
		counter = 1

		if unicode.IsDigit(currentRune) {
			return "", ErrInvalidString
		}
		if currentRune == escapeRune {
			if i+1 >= rLen {
				return "", ErrInvalidString
			}
			currentRune = runeSlice[i+1]

			if !unicode.IsDigit(currentRune) && currentRune != escapeRune {
				return "", ErrInvalidString
			}
			i++
		}
		if i+1 < rLen && unicode.IsDigit(runeSlice[i+1]) {
			counter, _ = strconv.Atoi(string(runeSlice[i+1]))
			i++
		}
		result.WriteString(strings.Repeat(string(currentRune), counter))
	}
	return result.String(), nil
}
