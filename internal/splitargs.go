package internal

import (
	"errors"
	"strings"
)

// SplitArgs splits a shell-like string into a slice of tokens.
// It handles single-quoted, double-quoted, and backslash-escaped tokens,
// so e.g. `-ldflags "-X main.Version=1.0"` yields ["-ldflags", "-X main.Version=1.0"].
// Returns an error if the string contains an unclosed quote or a trailing backslash.
func SplitArgs(s string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inSingle := false
	inDouble := false
	quoted := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '\\' && !inSingle:
			if i+1 >= len(s) {
				return nil, errors.New("trailing backslash in command")
			}
			i++
			current.WriteByte(s[i])
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			quoted = true
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			quoted = true
		case (ch == ' ' || ch == '\t') && !inSingle && !inDouble:
			if current.Len() > 0 || quoted {
				tokens = append(tokens, current.String())
				current.Reset()
				quoted = false
			}
		default:
			current.WriteByte(ch)
		}
	}

	if inSingle {
		return nil, errors.New("unclosed single quote in command")
	}
	if inDouble {
		return nil, errors.New("unclosed double quote in command")
	}

	if current.Len() > 0 || quoted {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}
