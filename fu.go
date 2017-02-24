package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/kbrgl/fu/matchers"
	"github.com/kbrgl/fu/shallowradix"
	"github.com/mattn/go-isatty"
	"github.com/stretchr/powerwalk"
)

const (
	// Version is the program version
	Version = "1.1.2"
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	fz := kingpin.Flag("fuzzy", "Use fuzzy search").Short('f').Bool()
	re := kingpin.Flag("regexp", "Use regexp-based search").Short('r').Bool()
	sfx := kingpin.Flag("suffix", "Use suffix-based search").Short('s').Bool()
	pfx := kingpin.Flag("prefix", "Use prefix-based search").Short('p').Bool()
	dir := kingpin.Flag("dir", "Show only directories").Short('d').Bool()
	perm := kingpin.Flag("perm", "Filter by Unix permissions").Short('m').Int()
	pll := kingpin.Flag("parallel",
		"Walk directories in parallel, may result in substantial speedups "+
			"for directories with many files").
		Short('c').
		Bool()

	query := kingpin.Arg("query", "Search query").Required().String()
	paths := kingpin.Arg("paths", "Paths to search").Default(".").ExistingDirs()
	kingpin.Version(Version)
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	walk := filepath.Walk
	if *pll {
		runtime.GOMAXPROCS(runtime.NumCPU())
		walk = powerwalk.Walk
	}

	var (
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
	default:
		nm = matchers.NewExactMatcher(*query)
	}

	ms := make([]matchers.FileMatcher, 0, 3)
	ms = append(ms, nm)

	if *perm != 0 {
		ms = append(ms, matchers.NewPermMatcher(os.FileMode(*perm)))
	}
	if *dir {
		ms = append(ms, matchers.NewDirMatcher())
	}

	s := shallowradix.New()

	for _, path := range *paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			fail(err)
		}
		s.Insert(fmt.Sprintf("%s%c", abs, os.PathSeparator))
	}

	// number of files found
	var found, traversed int
	start := time.Now()
	for _, path := range s.Prefixes() {
		walk(path, func(path string, fi os.FileInfo, err error) error {
			traversed++

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return nil
			}

			matches := true
			for _, matcher := range ms {
				if !matcher.Match(fi) {
					matches = false
				}
			}
			if /* the file */ matches {
				fmt.Println(path)
				found++
			}
			return nil
		})
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("\nTraversed %d files in %s, found %d matches.\n",
			traversed,
			time.Since(start),
			found)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
