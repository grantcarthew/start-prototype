# Documentation Review

You are performing a comprehensive documentation reconciliation review to identify "lint" - inconsistencies, deprecated references, stale content, and artifacts from iterative design.

## Scope

**Include:**

- All markdown files in docs/
- AGENTS.md
- README.md

**Exclude:**

- ddd/ directory (separate template project)
- PROJECT.md (design scratchpad, intentionally messy)
- Archive directories are okay to check but flag as low priority

## Critical Issues to Find

Run these exact commands and analyze results:

### 1. Broken Cross-References (CRITICAL)

```bash
# Check for old directory paths
rg '<any-discovered-old-dir>' docs/

# Check for broken DR links
rg '\[DR-[0-9]+\]' docs/ -n
```

**What to look for:**

- References to `design/decisions/` should be `design/design-records/`
- DR numbers that don't exist
- Incorrect relative paths

### 2. Deprecated Command References (HIGH)

```bash
# Find deprecated commands
rg 'start config .* add' docs/
```

**What to check:**

- Are references in historical/archive files? (OK - mark as intentional)
- Are they in migration docs showing old→new? (OK - mark as intentional)
- Are they in active DRs or user docs? (FIX - update to `start assets add`)

**Context matters:**

- DR-041, start-assets.md, start-assets-add.md showing comparisons = OK
- Archive files = OK
- Ideas/brainstorm files = OK
- Active DRs used as current examples = NOT OK

### 3. Stale Placeholders (MEDIUM)

```bash
# Find incomplete content
rg 'TODO|TBD|to be written|FIXME' docs/ -i -n
```

**What to check:**

- Is this in an archive? (Low priority)
- Is it actually incomplete or just a note?
- Should it reference an existing DR instead?

## Report Format

Provide a detailed report with:

**For each issue found:**

- File path and line number
- Issue type (broken link, deprecated ref, stale TODO, etc.)
- Exact problematic text (quote the line)
- Suggested fix (be specific)
- Priority: CRITICAL, HIGH, MEDIUM, or LOW

**Priority levels:**

- CRITICAL: Wrong information that would confuse/mislead users
- HIGH: Deprecated references in active docs, broken links
- MEDIUM: Stale TODOs, minor inconsistencies
- LOW: Archive files, cosmetic issues

**Summary at end:**

- Total issues by priority
- Which files need attention
- Verification commands to run after fixes

## Important Notes

1. **Historical context is OK**: Archive files and migration docs SHOULD mention old commands
2. **Be specific**: Always include file path and line number
3. **Suggest fixes**: Don't just identify problems, say how to fix them
4. **Verify patterns**: The commands above are the perfected search patterns - use them exactly
5. **Context matters**: Same text can be OK in one place (DR-041 explaining deprecation) but wrong in another (DR-033 using deprecated command as current example)

Begin the review now. Run each command and analyze the results carefully.
