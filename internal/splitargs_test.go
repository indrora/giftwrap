package internal

import (
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single token",
			input: "-trimpath",
			want:  []string{"-trimpath"},
		},
		{
			name:  "multiple flags",
			input: "-trimpath -v",
			want:  []string{"-trimpath", "-v"},
		},
		{
			name:  "quoted value with spaces",
			input: `-ldflags "-X main.Version=1.0 -s -w"`,
			want:  []string{"-ldflags", "-X main.Version=1.0 -s -w"},
		},
		{
			name:  "single-quoted value",
			input: `-ldflags '-X main.Version=1.0'`,
			want:  []string{"-ldflags", "-X main.Version=1.0"},
		},
		{
			name:  "backslash-escaped space",
			input: `-ldflags -X\ main.Version=1.0`,
			want:  []string{"-ldflags", "-X main.Version=1.0"},
		},
		{
			name:  "tabs as separators",
			input: "-trimpath\t-v",
			want:  []string{"-trimpath", "-v"},
		},
		{
			name:  "empty quoted string",
			input: `""`,
			want:  []string{""},
		},
		{
			name:    "unclosed double quote",
			input:   `-ldflags "unclosed`,
			wantErr: true,
		},
		{
			name:    "unclosed single quote",
			input:   `-ldflags 'unclosed`,
			wantErr: true,
		},
		{
			name:    "trailing backslash",
			input:   `-trimpath \`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitArgs(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SplitArgs(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("SplitArgs(%q) = %v, want %v", tt.input, got, tt.want)
			}
			for i, tok := range got {
				if tok != tt.want[i] {
					t.Errorf("token[%d] = %q, want %q", i, tok, tt.want[i])
				}
			}
		})
	}
}
