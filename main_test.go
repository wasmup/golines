package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	tests := []struct {
		name, content string
		all, nonempty int
	}{
		{"empty", "", 0, 0},
		{"final newline", "one\n\ntwo\n", 3, 2},
		{"no final newline", "one\n \t\nlast", 3, 2},
		{"unicode whitespace", "\u2003\n\u00a0", 2, 0},
		{"invalid UTF-8", "\xff\n", 1, 1},
		{"long line", strings.Repeat(" ", 100000) + "x", 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.go")
			if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}
			for _, mode := range []struct {
				skip bool
				want int
			}{{false, tt.all}, {true, tt.nonempty}} {
				got, err := countLines(path, mode.skip)
				if err != nil || got != mode.want {
					t.Errorf("skip=%v: got %d, %v; want %d", mode.skip, got, err, mode.want)
				}
			}
		})
	}
	if _, err := countLines(filepath.Join(t.TempDir(), "missing"), false); err == nil {
		t.Error("expected an error for a missing file")
	}
}

func TestCLIHelper(t *testing.T) {
	if os.Getenv("GOLINES_TEST_HELPER") != "1" {
		return
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = append([]string{os.Args[0]}, os.Args[3:]...)
	main()
	os.Exit(0)
}

func TestCLI(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string]string{
		"a.go": "a\n", "b.go": "b\n\n", "a_test.go": "t\n",
		"gen.pb.go": "p\n", "vendor/v.go": "v\n", ".hidden/h.go": "h\n",
		"internal/generated/g.go": "g\n", "other/generated/g.go": "g\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	// A dangling symlink must not prevent other files from being counted.
	if err := os.Symlink("missing", filepath.Join(root, "broken.go")); err != nil {
		t.Logf("symlink unavailable: %v", err)
	}
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.Command(os.Args[0], append([]string{"-test.run=^TestCLIHelper$", "--", "-dir", root}, args...)...)
		cmd.Env = append(os.Environ(), "GOLINES_TEST_HELPER=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command failed: %v\n%s", err, out)
		}
		return string(out)
	}
	out := run("-s", "-x", "internal/generated")
	if !strings.Contains(out, "Files: 3\nTotal lines: 3\n") || strings.Contains(out, "internal/generated") || strings.Contains(out, "broken.go") {
		t.Errorf("unexpected default output:\n%s", out)
	}
	if strings.Index(out, filepath.Join(root, "a.go")) > strings.Index(out, filepath.Join(root, "b.go")) {
		t.Errorf("equal line counts are not sorted by path:\n%s", out)
	}
	out = run("-t", "-v", "-p", "-x", "generated")
	if !strings.Contains(out, "Files: 5\nTotal lines: 6\n") || strings.Contains(out, ".hidden") {
		t.Errorf("unexpected included output:\n%s", out)
	}
}
