# DeepScan CLI 🔥

## Description

DeepScan is a powerful command-line tool designed to search for files based on filenames, paths, content, and encrypted secrets (SOPS). It works seamlessly in both local file systems and GitLab repositories.

## Features

- **OS File Search**: Search local directories based on filename, path, or content.
- **SOPS Decrypted Search**: Handle and search within SOPS-encrypted files.
- **GitHub Search**: Scan GitHub organizations for matching files.
- **GitLab Search**: Scan GitLab organizations for matching files.
- **Advanced Filtering**: Apply filters for filenames, paths, content, and regex patterns.
- **Exclusion Filters**: Refine search results by excluding specific files or directories.

## Installation

### Build from Source

Ensure you have Go installed, then run:

```sh
go build
```

## Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `os search` | Scans a specified directory for matching files. | `-d, --dir` The root directory to scan [default: "."] |
| `gitlab search` | Scans a GitLab organization for matching files. | `-o, --org` The GitLab organization to scan |
| `github search` | Scans a GitHub organization for matching files. | `-o, --org` The GitHub organization to scan |

### Global Flags

These flags apply to all commands:

```sh
-l, --log-level  Set the log level (DEBUG, INFO, WARN, ERROR) [default: INFO]
```

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
    --sops-only          Search for files that are only SOPS-encrypted
    --sops-key           Search for files encrypted with a specific key
```

#### Exclusion Filters

```sh
    --exclude-name           Exclude files with specific names (exact match)
    --exclude-name-contains  Exclude files with names containing this string
    --exclude-path           Exclude specific directories (exact match)
    --exclude-path-contains  Exclude directories containing this string
    --exclude-content        Exclude files containing specific content
```


#### Output

```sh
    --output                Output to file (JSON, YAML)
```
    

## Examples

### Local Filesystem Scanning

```sh
deep-scan os search -d /my/project
```

### GitLab Scanning

Set the required environment variables:

```sh
export GITLAB_TOKEN=your_access_token
```

Optionally, you can set the GitLab host (default is `gitlab.com`):

```sh
export GITLAB_HOST=https://gitlab.com
```

Then run:

```sh
deep-scan gitlab search -o my-org
```

### GitHub Scanning

Set the required environment variables:

```sh
export GITHUB_TOKEN=your_access_token
```

Optionally, you can set the GitHub host (default is `github.com`):

```sh
export GITHUB_HOST=https://github.com
```

Then run:

```sh
deep-scan github search -o my-org
```
