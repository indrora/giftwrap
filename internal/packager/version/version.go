package version

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Version wraps a parsed semver tag with additional git state.
type Version struct {
	*semver.Version        // Major(), Minor(), Patch(), Prerelease(), etc.
	Commits         int    // commits ahead of the tagged commit (0 = exact tag)
	Hash            string // short (7-char) SHA of HEAD
	Dirty           bool   // true if the working tree has uncommitted changes
}

// GetVersion opens the git repository at repoPath (pass "." for the current
// working directory) and returns a Version derived from the most recent
// semver-compatible tag reachable from HEAD.
//
// DetectDotGit is enabled, so repoPath can be any subdirectory of the repo.
func GetVersion(repoPath string) (Version, error) {
	repo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return Version{}, fmt.Errorf("version: open repo at %q: %w", repoPath, err)
	}

	head, err := repo.Head()
	if err != nil {
		return Version{}, fmt.Errorf("version: resolve HEAD: %w", err)
	}

	sv, tagHash, err := latestSemverTag(repo)
	if err != nil {
		return Version{}, err
	}

	commits, err := commitsSince(repo, tagHash, head.Hash())
	if err != nil {
		return Version{}, err
	}

	dirty, err := isDirty(repo)
	if err != nil {
		return Version{}, err
	}

	return Version{
		Version: sv,
		Commits: commits,
		Hash:    shortHash(head.Hash()),
		Dirty:   dirty,
	}, nil
}

// String returns a semver-compatible string.  When the repo is exactly on a
// tag and clean it is identical to the raw tag (e.g. "1.2.3").  Otherwise
// build metadata is appended:
//
//	1.2.3+4.gabcdef0        – 4 commits ahead, clean
//	1.2.3+4.gabcdef0.dirty  – 4 commits ahead, dirty working tree
//	1.2.3+dirty             – exactly on tag but dirty
func (v Version) String() string {
	s := v.Version.Original() // preserve the original tag string (e.g. "v1.2.3")

	var meta string
	if v.Commits > 0 {
		meta = fmt.Sprintf("%d.g%s", v.Commits, v.Hash)
	}
	if v.Dirty {
		if meta != "" {
			meta += ".dirty"
		} else {
			meta = "dirty"
		}
	}
	if meta != "" {
		s += "+" + meta
	}
	return s
}

// IsExact returns true when HEAD sits exactly on the tagged commit and the
// working tree has no uncommitted changes.  Useful as a release gate.
func (v Version) IsExact() bool {
	return v.Commits == 0 && !v.Dirty
}

// latestSemverTag iterates all refs/tags, picks the highest valid semver tag,
// and returns both the parsed version and the commit hash it resolves to.
func latestSemverTag(repo *git.Repository) (*semver.Version, plumbing.Hash, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, plumbing.ZeroHash, fmt.Errorf("version: list tags: %w", err)
	}

	var (
		latest     *semver.Version
		latestHash plumbing.Hash
	)

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		sv, err := semver.NewVersion(ref.Name().Short())
		if err != nil {
			return nil // silently skip non-semver tags
		}

		commitHash, err := resolveTagToCommit(repo, ref)
		if err != nil {
			return err
		}

		if latest == nil || sv.GreaterThan(latest) {
			latest = sv
			latestHash = commitHash
		}
		return nil
	})
	if err != nil {
		return nil, plumbing.ZeroHash, err
	}
	if latest == nil {
		return nil, plumbing.ZeroHash, fmt.Errorf("version: no semver tags found in repository")
	}
	return latest, latestHash, nil
}

// resolveTagToCommit dereferences both lightweight and annotated tags to the
// underlying commit hash.
func resolveTagToCommit(repo *git.Repository, ref *plumbing.Reference) (plumbing.Hash, error) {
	// TagObject returns an error for lightweight tags, so we use that to branch.
	if tag, err := repo.TagObject(ref.Hash()); err == nil {
		commit, err := tag.Commit()
		if err != nil {
			return plumbing.ZeroHash, fmt.Errorf("version: dereference annotated tag %q: %w", ref.Name().Short(), err)
		}
		return commit.Hash, nil
	}
	// Lightweight tag – the ref hash is already the commit hash.
	return ref.Hash(), nil
}

// commitsSince counts the commits reachable from head that are not reachable
// from base (i.e. commits made after the tagged commit).
func commitsSince(repo *git.Repository, base, head plumbing.Hash) (int, error) {
	if base == head {
		return 0, nil
	}

	iter, err := repo.Log(&git.LogOptions{From: head})
	if err != nil {
		return 0, fmt.Errorf("version: git log: %w", err)
	}
	defer iter.Close()

	count := 0
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash == base {
			return storer.ErrStop
		}
		count++
		return nil
	})
	if err != nil && err != storer.ErrStop {
		return 0, fmt.Errorf("version: walking commits: %w", err)
	}
	return count, nil
}

// isDirty reports whether the working tree contains uncommitted changes.
func isDirty(repo *git.Repository) (bool, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("version: open worktree: %w", err)
	}
	status, err := wt.Status()
	if err != nil {
		return false, fmt.Errorf("version: worktree status: %w", err)
	}
	return !status.IsClean(), nil
}

// shortHash returns the first 7 characters of a git hash, matching the
// default behaviour of `git describe`.
func shortHash(h plumbing.Hash) string {
	s := h.String()
	if len(s) > 7 {
		return s[:7]
	}
	return s
}
