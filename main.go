package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type fileStat struct {
	path  string
	lines int
}

func main() {
	root := flag.String("dir", ".", "root directory")
	includeTests := flag.Bool("t", false, "include _test.go files")
	includeVendor := flag.Bool("v", false, "include vendor directories")
	includeProto := flag.Bool("p", false, "include .pb.go files")
	skipEmpty := flag.Bool("s", false, "skip whitespace-only lines")

	var excluded []string

	// Extract -x and its directory arguments before standard flag parsing.
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "-x" {
			excluded = append(excluded, args[i+1:]...)
			args = args[:i]
			break
		}
	}

	if err := flag.CommandLine.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	excludedNames := make(map[string]bool)
	excludedPaths := make(map[string]bool)

	for _, dir := range excluded {
		dir = filepath.Clean(dir)

		if filepath.Dir(dir) == "." {
			excludedNames[dir] = true
		} else {
			excludedPaths[dir] = true
		}
	}

	var files []fileStat
	total := 0

	err := filepath.WalkDir(*root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == *root {
				return nil
			}

			if d.Name() == "vendor" && !*includeVendor {
				return filepath.SkipDir
			}

			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			rel, err := filepath.Rel(*root, path)
			if err != nil {
				return err
			}

			if excludedNames[d.Name()] || excludedPaths[filepath.Clean(rel)] {
				return filepath.SkipDir
			}

			return nil
		}

		name := d.Name()

		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if !*includeTests && strings.HasSuffix(name, "_test.go") {
			return nil
		}
		if !*includeProto && strings.HasSuffix(name, ".pb.go") {
			return nil
		}

		n, err := countLines(path, *skipEmpty)
		if err != nil {
			return err
		}

		files = append(files, fileStat{path: path, lines: n})
		total += n
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].lines == files[j].lines {
			return files[i].path < files[j].path
		}
		return files[i].lines < files[j].lines
	})

	for _, f := range files {
		fmt.Printf("%7d  %s\n", f.lines, f.path)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Files: %d\nTotal lines: %d\n", len(files), total)
}

func countLines(path string, skipEmpty bool) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	lines := 0
	hasData := false
	hasNonSpace := false

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		hasData = true
		if !unicode.IsSpace(r) {
			hasNonSpace = true
		}
		if r == '\n' {
			if !skipEmpty || hasNonSpace {
				lines++
			}
			hasData, hasNonSpace = false, false
		}
	}
	if hasData && (!skipEmpty || hasNonSpace) {
		lines++
	}

	return lines, nil
}
