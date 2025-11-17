# DR-001: Configuration File Format

- Date: 2025-01-03
- Status: Accepted
- Category: Configuration

## Problem

The tool needs a configuration format for storing settings, agents, roles, tasks, and contexts. The format must:

- Be human-editable (users will hand-edit their configs)
- Support comments (for inline documentation and guidance)
- Handle nested structures (agents with models, tasks with multiple fields, etc.)
- Be robust against formatting errors (users shouldn't break config with whitespace)

## Decision

Use TOML for all configuration files.

## Why

- Human-readable and editable
- No whitespace sensitivity (unlike YAML)
- Excellent Go support via BurntSushi/toml
- Supports comments and complex nested structures
- Used by similar tools (mise, Cargo)

## Trade-offs

Accept:

- Users must learn TOML syntax (less common than YAML/JSON)
- Requires external parsing library (BurntSushi/toml)
- Slightly more verbose than JSON for simple configs

Gain:

- Comments support (users can document their configs inline)
- No whitespace sensitivity (no YAML-style indentation errors)
- Nested structures for complex configuration
- Human-readable and approachable for hand-editing

## Alternatives

YAML:

- Pro: Widely known in DevOps community
- Pro: Native Go support (gopkg.in/yaml.v3)
- Con: Whitespace-sensitive, error-prone for hand-editing
- Con: Complex spec with surprising edge cases (Norway problem, etc.)
- Rejected: Error-prone editing outweighs familiarity benefits

JSON:

- Pro: Universal format, simple spec
- Pro: Excellent Go support (encoding/json)
- Con: No comments (dealbreaker - users can't document their config)
- Con: Less human-friendly (trailing commas errors, quoted keys required)
- Rejected: Lack of comments is unacceptable for user-editable config

Custom key-value format:

- Pro: Extremely simple, no learning curve
- Pro: Easy to parse without external library
- Con: Cannot handle nested structures (agents.claude.models)
- Con: No standard format, would need to design from scratch
- Rejected: Too limited for complex configuration needs
