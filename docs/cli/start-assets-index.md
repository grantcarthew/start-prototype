# start assets index

## Name

start assets index - Generate asset catalog index

## Synopsis

```bash
start assets index [flags]
```

## Description

Scans the current directory structure for asset metadata files (`.meta.toml`) and generates a unified `assets/index.csv` file. This command is intended for maintainers of the asset catalog repository.

It performs the following operations:

1. Validates the repository structure (must be a git repo with an `assets/` directory).
2. Recursively scans `assets/` for `.meta.toml` files.
3. Parses metadata (description, tags, creation date).
4. Calculates Git blob SHAs for version tracking.
5. Generates a sorted `assets/index.csv` file.

This index file is used by the `start assets search` and `start assets update` commands to perform fast, bandwidth-efficient operations.

## Requirements

- Current working directory must be the root of the asset catalog repository.
- `git` binary must be installed and accessible.
- `assets/` directory must exist.

## Output

The command generates or overwrites `assets/index.csv` with the following columns:

- **type**: Asset type (tasks, roles, agents, contexts)
- **category**: Directory name (e.g., git-workflow, general)
- **name**: Asset directory name
- **description**: Short description from `.meta.toml`
- **tags**: Comma-separated tags from `.meta.toml`
- **bin**: (Agents only) Binary name required
- **sha**: Git blob SHA of the main asset file
- **size**: File size in bytes
- **created**: Creation date (ISO 8601)
- **updated**: Last modification date (ISO 8601)

## Examples

### Generate Index

```bash
$ cd ~/projects/start-catalog
$ start assets index

Validating repository structure...
✓ Git repository detected
✓ Assets directory found

Scanning assets/...
Found 46 assets

Sorting assets (type → category → name)...
Writing index to assets/index.csv...

✓ Generated index with 46 assets
Updated: assets/index.csv

Ready to commit:
  git add assets/index.csv
  git commit -m "Regenerate catalog index"
```

## Flags

**--verify**
: Verify the existing index against the file system without modifying it. Returns exit code 1 if out of sync. Useful for CI/CD pipelines.

## Exit Codes

**0** - Success
**1** - Verification failed (with --verify) or general error

## See Also

- start-assets(1) - Asset management overview
