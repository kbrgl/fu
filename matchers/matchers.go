package matchers

import (
	"os"
	"regexp"
	"strings"

	"github.com/kbrgl/fuzzy"
)

// FileMatcher is an interface providing a Match method that checks whether a
// file satisfies some constraint.
type FileMatcher interface {
	Match(os.FileInfo) bool
}

// FuzzyMatcher fuzzy-matches the filename.
type FuzzyMatcher struct {
	pattern string
}

// NewFuzzyMatcher returns a new FuzzyMatcher.
func NewFuzzyMatcher(pattern string) *FuzzyMatcher {
	return &FuzzyMatcher{pattern: pattern}
}

// Match fuzzy-matches the filename.
func (f FuzzyMatcher) Match(fi os.FileInfo) bool {
	return fuzzy.MatchFold(fi.Name(), f.pattern)
}

// ExactMatcher checks the filename for exact equality with the expected value.
type ExactMatcher struct {
	expected string
}

// NewExactMatcher returns a new ExactMatcher.
func NewExactMatcher(expected string) *ExactMatcher {
	return &ExactMatcher{expected: expected}
}

// Match checks that the filename is exactly equal to the expected value.
func (e ExactMatcher) Match(fi os.FileInfo) bool {
	return fi.Name() == e.expected
}

// SuffixMatcher checks whether the filename has some suffix.
type SuffixMatcher struct {
	suffix string
}

// NewSuffixMatcher returns a SuffixMatcher that checks a string for the provided
// suffix.
func NewSuffixMatcher(suffix string) *SuffixMatcher {
	return &SuffixMatcher{suffix: suffix}
}

// Match matches on the suffix of fi.Name().
func (s SuffixMatcher) Match(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), s.suffix)
}

// PrefixMatcher checks whether the filename has some prefix.
type PrefixMatcher struct {
	prefix string
}

// NewPrefixMatcher returns a PrefixMatcher that checks a string for the provided
// prefix.
func NewPrefixMatcher(prefix string) *PrefixMatcher {
	return &PrefixMatcher{prefix: prefix}
}

// Match matches on the prefix of fi.Name().
func (p PrefixMatcher) Match(fi os.FileInfo) bool {
	return strings.HasPrefix(fi.Name(), p.prefix)
}

// REMatcher checks whether the filename matches a regexp pattern.
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

// PermMatcher checks against the provided permissions.
type PermMatcher struct {
	perm os.FileMode
}

// NewPermMatcher returns a new PermMatcher that performs a check against the
// provided FileMode.
func NewPermMatcher(perm os.FileMode) *PermMatcher {
	return &PermMatcher{perm: perm}
}

// Match compares fi's permissions to the permissions set on p.
func (p PermMatcher) Match(fi os.FileInfo) bool {
	return fi.Mode().Perm()&p.perm != 0
}

// DirMatcher allows only dirs.
type DirMatcher struct {
}

// NewDirMatcher returns a new DirMatcher.
func NewDirMatcher() *DirMatcher {
	return &DirMatcher{}
}

// Match returns true for dirs.
func (d DirMatcher) Match(fi os.FileInfo) bool {
	return fi.IsDir()
}

// AllMatcher allows everything.
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

// SubstringMatcher checks the provided string for a substring.
type SubstringMatcher struct {
	// substring to look for
	substring string
}

// NewSubstringMatcher returns a new SubstringMatcher.
func NewSubstringMatcher(substring string) *SubstringMatcher {
	return &SubstringMatcher{substring: substring}
}

// Match searches the filename for a given substring and returns true if it is
// present.
func (s SubstringMatcher) Match(fi os.FileInfo) bool {
	return strings.Contains(fi.Name(), s.substring)
}
