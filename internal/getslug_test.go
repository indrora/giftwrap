package internal

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello world", "hello-world"},
		{"Hello World", "hello-world"},
		{"UPPER CASE", "upper-case"},
		{"already-slugified", "already-slugified"},
		{"hello_world", "hello-world"},
		{"hello   world", "hello-world"},
		{"foo123bar", "foo123bar"},
		{"foo 123 bar", "foo-123-bar"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"special!@#chars", "special-chars"},
		{"", ""},
		{"café", "cafe"},
		{"foo--bar", "foo-bar"},
		{"hello 🦊", "hello"},
		{"this 🐍 hisses", "this-hisses"},
		{"naïve", "naive"},
		{"résumé", "resume"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
