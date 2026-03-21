package runner

import "github.com/indrora/giftwrap/internal"

// splitCommand splits a shell-like command string into the executable and its
// arguments. It handles single-quoted, double-quoted, and backslash-escaped
// tokens so that e.g. `echo "hello world"` yields ("echo", ["hello world"]).
// Returns an error if the string contains an unclosed quote or a trailing backslash.
func splitCommand(s string) (string, []string, error) {
	tokens, err := internal.SplitArgs(s)
	if err != nil {
		return "", nil, err
	}
	if len(tokens) == 0 {
		return "", nil, nil
	}
	return tokens[0], tokens[1:], nil
}
