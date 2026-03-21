package version

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
)

// mkv is a test helper that parses a semver string or panics.
func mkv(s string) *semver.Version {
	v, err := semver.NewVersion(s)
	if err != nil {
		panic(err)
	}
	return v
}

// ----------------------------------------------------------------------------
// Version.String
// ----------------------------------------------------------------------------

func TestVersionString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    Version
		want string
	}{
		{
			name: "exact clean tag without v prefix",
			v:    Version{Version: mkv("1.2.3")},
			want: "1.2.3",
		},
		{
			name: "exact clean tag with v prefix",
			v:    Version{Version: mkv("v1.2.3")},
			want: "v1.2.3",
		},
		{
			name: "commits ahead clean",
			v:    Version{Version: mkv("1.2.3"), Commits: 4, Hash: "abcdef0"},
			want: "1.2.3+4.gabcdef0",
		},
		{
			name: "exact tag dirty",
			v:    Version{Version: mkv("1.2.3"), Dirty: true},
			want: "1.2.3+dirty",
		},
		{
			name: "commits ahead and dirty",
			v:    Version{Version: mkv("1.2.3"), Commits: 2, Hash: "1234567", Dirty: true},
			want: "1.2.3+2.g1234567.dirty",
		},
		{
			name: "pre-release tag exact",
			v:    Version{Version: mkv("1.0.0-alpha.1")},
			want: "1.0.0-alpha.1",
		},
		{
			name: "pre-release tag with commits",
			v:    Version{Version: mkv("1.0.0-alpha.1"), Commits: 1, Hash: "deadbee"},
			want: "1.0.0-alpha.1+1.gdeadbee",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.v.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Version.IsExact
// ----------------------------------------------------------------------------

func TestIsExact(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    Version
		want bool
	}{
		{"clean exact tag", Version{Version: mkv("1.0.0")}, true},
		{"commits ahead", Version{Version: mkv("1.0.0"), Commits: 1}, false},
		{"dirty exact tag", Version{Version: mkv("1.0.0"), Dirty: true}, false},
		{"commits ahead and dirty", Version{Version: mkv("1.0.0"), Commits: 3, Dirty: true}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.v.IsExact(); got != tc.want {
				t.Errorf("IsExact() = %v, want %v", got, tc.want)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Version.RCString
// ----------------------------------------------------------------------------

func TestRCString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    Version
		want string
	}{
		{
			name: "exact clean tag — unchanged",
			v:    Version{Version: mkv("v1.2.3")},
			want: "v1.2.3",
		},
		{
			name: "commits ahead clean — patch bumped",
			v:    Version{Version: mkv("v1.2.3"), Commits: 5},
			want: "v1.2.4+rc5",
		},
		{
			name: "commits ahead and dirty",
			v:    Version{Version: mkv("v1.2.3"), Commits: 3, Dirty: true},
			want: "v1.2.4+rc3.dirty",
		},
		{
			name: "dirty only — no patch bump",
			v:    Version{Version: mkv("v1.2.3"), Dirty: true},
			want: "v1.2.3+dirty",
		},
		{
			name: "patch zero — bumps to 1",
			v:    Version{Version: mkv("v2.1.0"), Commits: 1},
			want: "v2.1.1+rc1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.v.RCString(); got != tc.want {
				t.Errorf("RCString() = %q, want %q", got, tc.want)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// shortHash
// ----------------------------------------------------------------------------

func TestShortHash(t *testing.T) {
	t.Parallel()

	cases := []struct {
		raw  string // 40-char hex for plumbing.NewHash
		want string
	}{
		{"abcdef0123456789abcdef0123456789abcdef01", "abcdef0"},
		{"1234567000000000000000000000000000000000", "1234567"},
	}

	for _, tc := range cases {
		t.Run(tc.raw[:7], func(t *testing.T) {
			t.Parallel()
			if got := shortHash(plumbing.NewHash(tc.raw)); got != tc.want {
				t.Errorf("shortHash(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// commitsSince – same-hash early-return path
// ----------------------------------------------------------------------------

func TestCommitsSince_SameHash(t *testing.T) {
	t.Parallel()

	h := plumbing.NewHash("abcdef0123456789abcdef0123456789abcdef01")

	// When base == head the function returns 0 immediately without touching the
	// repository, so passing nil is safe here.
	count, err := commitsSince(nil, h, h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("commitsSince with identical hashes = %d, want 0", count)
	}
}
