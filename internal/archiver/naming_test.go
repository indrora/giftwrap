package archiver

import (
	"testing"
)

func TestArchiveName(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		targetName    string
		variant       string
		format        ArchiveFormat
		singleDefault bool
		want          string
		wantErr       bool
	}{
		{
			name:          "single default target omits target name",
			projectName:   "myproject",
			targetName:    "default",
			variant:       "linux/amd64",
			format:        FormatTarGz,
			singleDefault: true,
			want:          "myproject-linux-amd64.tar.gz",
		},
		{
			name:          "multiple targets include target name",
			projectName:   "myproject",
			targetName:    "foo",
			variant:       "linux/amd64",
			format:        FormatTarGz,
			singleDefault: false,
			want:          "myproject-foo-linux-amd64.tar.gz",
		},
		{
			name:          "zip format",
			projectName:   "myapp",
			targetName:    "default",
			variant:       "windows/amd64",
			format:        FormatZip,
			singleDefault: true,
			want:          "myapp-windows-amd64.zip",
		},
		{
			name:          "tar.zst format with target name",
			projectName:   "tool",
			targetName:    "release",
			variant:       "darwin/arm64",
			format:        FormatTarZst,
			singleDefault: false,
			want:          "tool-release-darwin-arm64.tar.zst",
		},
		{
			name:        "invalid variant missing slash",
			projectName: "myproject",
			targetName:  "default",
			variant:     "linuxamd64",
			format:      FormatTarGz,
			wantErr:     true,
		},
		{
			name:        "empty variant",
			projectName: "myproject",
			targetName:  "default",
			variant:     "",
			format:      FormatTarGz,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArchiveName(tt.projectName, tt.targetName, tt.variant, tt.format, tt.singleDefault)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArchiveName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArchiveName() = %q, want %q", got, tt.want)
			}
		})
	}
}
