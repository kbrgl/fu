package matchers

import (
	"os"
	"regexp"
	"strings"

	"github.com/kbrgl/fuzzy"
)

// FileMatcher is an interface providing a Match method that can match on a
// file.
type FileMatcher interface {
	Match(os.FileInfo) bool
}

// FuzzyMatcher is a FileMatcher that fuzzy-matches on the filename.
type FuzzyMatcher struct {
	pattern string
}

// NewFuzzyMatcher returns a new FuzzyMatcher.
func NewFuzzyMatcher(pattern string) *FuzzyMatcher {
	return &FuzzyMatcher{pattern: pattern}
}

// Match fuzzy matches on fi.Name().
func (f FuzzyMatcher) Match(fi os.FileInfo) bool {
	return fuzzy.MatchFold(fi.Name(), f.pattern)
}

// ExactMatcher is a FileMatcher that checks whether the filename is equal to
// some value.
type ExactMatcher struct {
	expected string
}

// NewExactMatcher returns a new ExactMatcher.
func NewExactMatcher(expected string) *ExactMatcher {
	return &ExactMatcher{expected: expected}
}

// Match exact matches on fi.Name().
func (e ExactMatcher) Match(fi os.FileInfo) bool {
	return fi.Name() == e.expected
}

// SuffixMatcher is a FileMatcher that checks whether the filename has some
// suffix.
type SuffixMatcher struct {
	suffix string
}

// NewSuffixMatcher returns a SuffixMatcher that performs a check against the
// provided suffix.
func NewSuffixMatcher(suffix string) *SuffixMatcher {
	return &SuffixMatcher{suffix: suffix}
}

// Match matches on the suffix of fi.Name().
func (s SuffixMatcher) Match(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), s.suffix)
}

// PrefixMatcher is a FileMatcher that checks whether the filename has some
// prefix.
type PrefixMatcher struct {
	prefix string
}

// NewPrefixMatcher returns a PrefixMatcher that performs a check against the
// provided prefix.
func NewPrefixMatcher(prefix string) *PrefixMatcher {
	return &PrefixMatcher{prefix: prefix}
}

// Match matches on the prefix of fi.Name().
func (p PrefixMatcher) Match(fi os.FileInfo) bool {
	return strings.HasPrefix(fi.Name(), p.prefix)
}

// REMatcher is a FileMatcher that checks whether the filename matches a regexp
// pattern.
type REMatcher struct {
	pattern *regexp.Regexp
}

// NewREMatcher returns a new REMatcher that performs a check against the
// provided pattern.
func NewREMatcher(pattern string) (*REMatcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &REMatcher{pattern: re}, nil
}

// Match regex-matches on fi.Name().
func (r REMatcher) Match(fi os.FileInfo) bool {
	return r.pattern.MatchString(fi.Name())
}

// PermMatcher is a FileMatcher that checks against the provided permissions.
type PermMatcher struct {
	perm os.FileMode
}

// NewPermMatcher returns a new PermMatcher that performs a check against the
// provided FileMode.
func NewPermMatcher(perm os.FileMode) *PermMatcher {
	return &PermMatcher{perm: perm}
}

// Match compares fi's Permissions to the permissions set on p.
func (p PermMatcher) Match(fi os.FileInfo) bool {
	return fi.Mode().Perm()&p.perm != 0
}

// DirMatcher is a FileMatcher that filters dirs.
type DirMatcher struct {
}

// NewDirMatcher returns a new DirMatcher.
func NewDirMatcher() *DirMatcher {
	return &DirMatcher{}
}

// Match matches dirs.
func (d DirMatcher) Match(fi os.FileInfo) bool {
	return fi.IsDir()
}

// AllMatcher is a FileMatcher matches everything.
type AllMatcher struct {
}

// NewAllMatcher returns a new AllMatcher.
func NewAllMatcher() *AllMatcher {
	return &AllMatcher{}
}

// Match always returns true.
func (a AllMatcher) Match(_ os.FileInfo) bool {
	return true
}
