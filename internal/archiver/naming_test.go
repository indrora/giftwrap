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

func TestArchiveNameFromTemplate(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		data    ArchiveNameData
		want    string
		wantErr bool
	}{
		{
			name: "version in stem",
			tmpl: "{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}",
			data: ArchiveNameData{
				Name: "mytool", Version: "v1.2.3",
				OS: "linux", Arch: "amd64", Format: "tar.gz",
			},
			want: "mytool-v1.2.3-linux-amd64.tar.gz",
		},
		{
			name: "target included",
			tmpl: "{{.Name}}-{{.Target}}-{{.OS}}-{{.Arch}}",
			data: ArchiveNameData{
				Name: "mytool", Target: "release",
				OS: "darwin", Arch: "arm64", Format: "tar.gz",
			},
			want: "mytool-release-darwin-arm64.tar.gz",
		},
		{
			name: "semver components",
			tmpl: "{{.Name}}-{{.Major}}.{{.Minor}}.{{.Patch}}-{{.OS}}-{{.Arch}}",
			data: ArchiveNameData{
				Name: "tool", Major: 2, Minor: 0, Patch: 1,
				OS: "windows", Arch: "amd64", Format: "zip",
			},
			want: "tool-2.0.1-windows-amd64.zip",
		},
		{
			name: "extension always appended even when format not in template",
			tmpl: "{{.Name}}-{{.OS}}",
			data: ArchiveNameData{
				Name: "app", OS: "linux", Arch: "amd64", Format: "tar.zst",
			},
			want: "app-linux.tar.zst",
		},
		{
			name:    "invalid template syntax",
			tmpl:    "{{.Name}-{{.OS}}",
			data:    ArchiveNameData{Name: "app", OS: "linux", Format: "tar.gz"},
			wantErr: true,
		},
		{
			name:    "unknown field causes error",
			tmpl:    "{{.Nonexistent}}",
			data:    ArchiveNameData{Format: "tar.gz"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArchiveNameFromTemplate(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArchiveNameFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArchiveNameFromTemplate() = %q, want %q", got, tt.want)
			}
		})
	}
}
