# Go Lines Counter

A lightweight command-line tool written in Go to recursively count lines of Go source code across a project.

Displays line counts for individual files, sorted in ascending order, followed by the total number of files and lines.

Built entirely with the **Go standard library**, with no external dependencies.

## Features

- Recursively scans directories for `.go` files.
- Displays line counts per file, sorted from smallest to largest.
- Excludes Go test files (`_test.go`) by default.
- Excludes vendor directories (`vendor`) by default.
- Excludes Protocol Buffers generated files (`.pb.go`) by default.
- Optionally ignores empty and whitespace-only lines.
- Supports excluding multiple directories by name or relative path.
- Automatically skips hidden directories such as `.git`.
- Supports specifying a custom root directory.
- Handles files without trailing newline characters.
- Uses buffered file reading for efficient processing.
- Requires no third-party dependencies.

## Installation

### Requirements

- Go 1.27.3 or later.
- Compatible with Linux, macOS, and Windows.

### Build from source

Clone the repository:

```bash
git clone https://github.com/wasmup/golines.git
cd golines
```

Build the executable:

```bash
go build -o golines .
```

Optionally install it into your Go binary directory:

```bash
go install .

# or
go install -trimpath -ldflags=-s
```


## Usage

```bash
golines [flags]
```

By default, the tool scans the current directory and its subdirectories, excluding test files, vendor directories, protobuf-generated files, and hidden directories.

### Command-line flags

| Flag | Default | Description |
|---|---|---|
| `-dir` | `.` | Root directory to scan |
| `-t` | `false` | Include `_test.go` files |
| `-v` | `false` | Include `vendor` directories |
| `-p` | `false` | Include `.pb.go` files |
| `-s` | `false` | Skip empty and whitespace-only lines |
| `-x` | None | Exclude directories by name or relative path |
| `-h` | — | Display help |

### Examples

**Count Go source lines in the current directory**

```bash
golines
```

**Exclude empty and whitespace-only lines**

```bash
golines -s
```

**Include test files**

```bash
golines -t
```

**Include vendor directories**

```bash
golines -v
```

**Include protobuf-generated files**

```bash
golines -p
```

**Include all Go files**

```bash
golines -t -v -p
```

**Scan a specific directory**

```bash
golines -dir ./internal
```

**Exclude multiple directories**

```bash
golines -x mocks generated tmp
```

**Exclude directories using relative paths**

```bash
golines -x internal/generated pkg/mocks
```

**Combine multiple options**

```bash
golines -dir ./src -s -t -x mocks generated
```

> **Important:** The `-x` flag must be the last flag in the command. Every argument following `-x` is interpreted as an excluded directory.

## Example output

```text
     12  internal/config.go
     38  internal/errors.go
     94  cmd/server/main.go
    127  internal/logger.go
    285  internal/database.go
    412  internal/handler.go
------------------------------------------------------------
Files: 6
Total lines: 968
```

Files are sorted by line count in ascending order. Files with identical line counts are sorted alphabetically by path.

## Directory exclusion

The `-x` flag supports two matching strategies.

### Exclude by directory name

```bash
golines -x mocks
```

Excludes all directories named `mocks`, regardless of their location:

```text
project/
├── internal/
│   ├── mocks/        # Excluded
│   └── service/
├── pkg/
│   └── mocks/        # Excluded
└── mocks/            # Excluded
```

### Exclude by relative path

```bash
golines -x internal/mocks
```

Excludes only `internal/mocks` and its descendants.

Other directories named `mocks` remain included.

Directory paths are evaluated relative to the root directory specified by `-dir`.

### Default exclusions

Unless explicitly enabled, the following are excluded:

| Type | Pattern | Enable with |
|---|---|---|
| Go tests | `*_test.go` | `-t` |
| Vendor directories | `vendor/` | `-v` |
| Protobuf files | `*.pb.go` | `-p` |
| Hidden directories | `.*/` | Always excluded |

## Line counting behavior

The tool counts physical lines rather than logical statements or expressions.

For example:

```go
package main

import "fmt"

func main() {
    // Print a message
    fmt.Println("Hello, World!")
}
```

Results:

| Mode | Lines counted |
|---|---:|
| Default | 8 |
| With `-s` | 6 |

When `-s` is enabled, lines containing only spaces, tabs, newlines, or other Unicode whitespace are ignored.

Comment-only lines are still counted.

## Performance

The tool uses:

- `filepath.WalkDir` for recursive directory traversal.
- `bufio.Reader` for buffered file reading.
- `unicode.IsSpace` for whitespace detection.
- `sort.Slice` for ordering results.
- `flag` for command-line argument parsing.

Each file is processed sequentially, without loading the entire file into memory.

Memory usage is primarily determined by the buffered reader and the number of file statistics retained for sorting.

## Limitations

- Counts physical lines, not logical lines of code.
- Does not exclude comments.
- Does not distinguish executable code from declarations.
- Hidden directories are always skipped, except when the hidden directory itself is selected as the root.
- `-x` must appear last.
- Symbolic links to directories are not followed.
- Symbolic links and other non-regular files are not counted.
- File reading and traversal errors terminate execution.

## License

The MIT License.

## Contributing

Contributions, bug reports, and feature requests are welcome. Open a GitHub issue or submit a pull request.
