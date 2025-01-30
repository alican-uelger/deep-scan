# DeepScan CLI

## Overview

Deep Scan is a powerful command-line tool designed to search for files based on filenames, paths, content, and encrypted secrets (SOPS). It works seamlessly in both local fileMatch systems and GitLab repositories.

## Features

- **Fast and flexible fileMatch scanning** in local directories or GitLab organizations.
- **Advanced filtering** by filename, directory path, content, and regex patterns.
- **SOPS decryption support** for handling encrypted files.
- **Exclusion filters** to refine search results.
- **Configurable** via CLI flags and environment variables.

## Installation

### Build from Source

Ensure you have Go installed, then run:

```sh
make build
```

This will generate the binary in `./dist/local/deep-scan`.

## Usage

### Global Flags

These flags apply to all commands:

```sh
-l, --log-level  Set the log level (DEBUG, INFO, WARN, ERROR) [default: INFO]
```

## Commands

### Local Filesystem Scanning (`os`)

#### Description
Scans a specified directory for matching files.

#### Example

```sh
deep-scan os search -d /my/project
```

#### Flags

```sh
-d, --dir  The root directory to scan [default: "."]
```

---

### GitLab Scanning (`gitlab`)

#### Description
Scans a GitLab organization for matching files.

#### Environment Variables

Before scanning GitLab, set the required environment variables:

```sh
export GITLAB_HOST=https://gitlab.com
export GITLAB_TOKEN=your_access_token
```

#### Example

```sh
deep-scan gitlab search -o my-org
```

#### Flags

```sh
-o, --g-org  The GitLab organization to scan
```

---

### Search Filters

#### Filename Filters

```sh
-n, --name               Search for files with specific names (exact match)
    --name-contains      Search for files with names containing this string
    --name-regex         Search for files with names matching this regex
```

#### Path Filters

```sh
-p, --path               Search in specific directories (exact match)
    --path-contains      Search in directories containing this string
    --path-regex         Search in directories matching this regex
```

#### Content Filters

```sh
-c, --content            Search for files containing specific content
    --content-regex      Search for files containing content matching this regex
```

#### SOPS Filters

```sh
-s, --sops               Search for SOPS-encrypted files
    --sops-key           Specify SOPS keys for decryption
```

#### Exclusion Filters

```sh
    --exclude-name           Exclude files with specific names (exact match)
    --exclude-name-contains  Exclude files with names containing this string
    --exclude-path           Exclude specific directories (exact match)
    --exclude-path-contains  Exclude directories containing this string
    --exclude-content        Exclude files containing specific content
```

---

## Development

### Run Code Quality Checks

```sh
make check
```

### Run Tests

#### End-to-End Tests

```sh
make test-e2e
```

#### Integration Tests

```sh
make test-int
```

#### Unit Tests

```sh
make test-unit
```

### Generate Mocks

```sh
make mocks
```

## License

This project is licensed under the MIT License.

