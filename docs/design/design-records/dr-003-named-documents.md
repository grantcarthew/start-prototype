# DR-003: Named Documents for Context

- Date: 2025-01-03
- Status: Accepted
- Category: Configuration

## Problem

Context documents need to be configurable at both global and local scopes. The configuration format must support:

- Overriding specific context documents in local config (replace global definition)
- Adding new context documents in local config (extend global)
- Clear identification of which document is being configured
- Merge-friendly structure that works with the global/local merge strategy

## Decision

Use named TOML sections for context documents instead of arrays.

## Why

Names allow precise targeting:

- Local config can override specific documents by using the same name
- Local config can add new documents by using different names
- Arrays cannot be targeted for override (would have to replace entire array)
- Enables both override pattern (same name) and add pattern (new name)
- More explicit and readable than array indices

Merge-friendly:

- Same name in local overrides global (replace)
- New name in local adds to global (extend)
- Works naturally with TOML merge behavior

## Trade-offs

Accept:

- Users must choose meaningful names for contexts
- Slightly more verbose than simple array syntax
- Names must be unique within a scope

Gain:

- Precise control over which contexts to override
- Can override one context without affecting others
- Clear, self-documenting configuration
- Natural merge behavior with global/local scopes

## Structure

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."

[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
```

Example merge behavior:

Global config:
```toml
[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
```

Local config:
```toml
# Override: Replace global "project" context
[contexts.project]
file = "~/multi-repo/BIG-PROJECT.md"
prompt = "Read {file} for project context."

# Add: New context not in global
[contexts.vision]
file = "./docs/vision.md"
prompt = "Read {file} for product vision."
```

Result after merge:
- `context.project` uses local definition (overridden)
- `context.vision` added from local (new)

## Alternatives

Array-based structure:

```toml
contexts = [
  { name = "environment", file = "~/reference/ENVIRONMENT.md", prompt = "..." },
  { name = "project", file = "./PROJECT.md", prompt = "..." }
]
```

- Pro: Simpler syntax for simple cases
- Pro: All contexts in one structure
- Con: Cannot override specific array items during merge
- Con: Local config must replace entire array to override one item
- Con: No way to "add" to array without duplicating global entries
- Rejected: Incompatible with override/extend merge pattern

Separate files per context:

- Pro: Maximum flexibility, each context is independent file
- Pro: Very easy to override (replace file)
- Con: Too many small files to manage
- Con: Harder to see all contexts at once
- Con: Excessive file system overhead
- Rejected: Too granular, over-engineering the problem
