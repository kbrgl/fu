package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/color"
	"github.com/kbrgl/fu/matchers"
	"github.com/kbrgl/fu/shallowradix"
	isatty "github.com/mattn/go-isatty"
	"github.com/stretchr/powerwalk"
)

const (
	// Version is the program version
	Version = "2.2.0"
)

var (
	fz        = kingpin.Flag("fuzzy", "Use fuzzy search").Short('f').Bool()
	re        = kingpin.Flag("regexp", "Use regexp-based search").Short('r').Bool()
	sfx       = kingpin.Flag("suffix", "Use suffix-based search (short flag 'a' is short for 'after')").Short('a').Bool()
	pfx       = kingpin.Flag("prefix", "Use prefix-based search (short flag 'b' is short for 'before')").Short('b').Bool()
	substring = kingpin.Flag("substring",
		"Use substring-based search allowing the query to be at any position in the filename").
		Short('s').Bool()
	dir     = kingpin.Flag("dirs", "Show only directories").Short('d').Bool()
	perm    = kingpin.Flag("perm", "Filter by Unix permissions").Short('m').Int()
	seq     = kingpin.Flag("seq", "Search directories sequentially").Bool()
	older   = kingpin.Flag("older", "Filter by age (modification time)").Short('o').Duration()
	younger = kingpin.Flag("younger", "Filter by age (modification time)").Short('y').Duration()
	exclude = kingpin.Flag("exclude", "Excludes files matching the filters").Short('e').Bool()
	_       = kingpin.Flag("parallel",
		"[DEPRECATED: see --seq] Walk directories in parallel, may result in substantial speedups "+
			"for directories with many files").
		Short('c').Bool()

	query = kingpin.Arg("query", "Search query").Required().String()
	paths = kingpin.Arg("paths", "Paths to search").Default(".").ExistingDirs()
)

func init() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version(Version)
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()
}

func main() {
	walk := getWalker()
	ms := getMatchers()

	// number of files found / traversed
	var found, traversed uint64
	start := time.Now()
	pathPrefixes := getPathPrefixes(*paths)
	for _, path := range pathPrefixes {
		_ = walk(path, func(path string, fi os.FileInfo, err error) error {
			atomic.AddUint64(&traversed, 1)

			if err != nil {
				if isatty.IsTerminal(os.Stderr.Fd()) {
					errorf(err)
				}
				return nil
			}

			matches := true
			for _, matcher := range ms {
				if matcher.Match(fi) == *exclude {
					matches = false
					break
				}
			}

			if /* the file */ matches {
				fmt.Println(path)
				atomic.AddUint64(&found, 1)
			}

			return nil
		})
	}

	// Print message to standard error, since we don't want it to appear when
	// redirecting standard output.
	fmt.Fprintf(os.Stderr,
		"\nTraversed %d files in %s, found %d matches.\n",
		traversed,
		time.Since(start),
		found)
}

func getWalker() func(string, filepath.WalkFunc) error {
	walk := filepath.Walk
	if !*seq {
		runtime.GOMAXPROCS(runtime.NumCPU())
		walk = powerwalk.Walk
	}
	return walk
}

func getMatchers() []matchers.FileMatcher {
	// Filename matchers can't be stacked (there can only be one), so check the
	// flags and choose one.
	// There's currently no way to fail if multiple filename matchers are provided,
	// so the ordering of cases here is pretty much arbitrary.
	nm := getNameMatcher()

	ms := make([]matchers.FileMatcher, 0, 3)
	ms = append(ms, nm)

	// Append stackable matchers to the matchers slice
	if *perm != 0 {
		ms = append(ms, matchers.NewPermMatcher(os.FileMode(*perm)))
	}
	if *dir {
		ms = append(ms, matchers.NewDirMatcher())
	}
	if *older != 0 {
		ms = append(ms, matchers.NewAgeOlderMatcher(*older))
	}
	if *younger != 0 {
		ms = append(ms, matchers.NewAgeYoungerMatcher(*younger))
	}

	return ms
}

func getNameMatcher() matchers.FileMatcher {
	var (
		// name matcher
		nm  matchers.FileMatcher
		err error
	)
	switch {
	case *re:
		nm, err = matchers.NewREMatcher(*query)
		if err != nil {
			fail(err)
		}
	case *sfx:
		nm = matchers.NewSuffixMatcher(*query)
	case *pfx:
		nm = matchers.NewPrefixMatcher(*query)
	case *fz:
		nm = matchers.NewFuzzyMatcher(*query)
	case *substring:
		nm = matchers.NewSubstringMatcher(*query)
	default:
		if len(*query) == 0 {
			nm = matchers.NewAllMatcher()
		} else {
			nm = matchers.NewExactMatcher(*query)
		}
	}

	return nm
}

func getPathPrefixes(paths []string) []string {
	rdx := shallowradix.New()

	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			fail(err)
		}
		rdx.Insert(fmt.Sprintf("%s%c", abs, os.PathSeparator))
	}

	return rdx.Prefixes()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func errorf(err error) {
	if isatty.IsTerminal(os.Stderr.Fd()) {
		color.Set(color.FgRed)
		defer color.Unset()
	}
	fmt.Fprintf(os.Stderr, "%v\n", err)
}
