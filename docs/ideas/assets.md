# Asset Ideas

This document tracks ideas for assets to add to the catalog (roles, tasks, contexts, etc.).

Assets are stored in the GitHub catalog and downloaded on-demand via lazy-loading. This file helps us plan what to build when implementation begins.

---

## Roles

### meta/role-writer.md

**Purpose**: A specialized role for writing role documents from natural language descriptions.

**Use case**: When users want to generate a custom role document, they can use this meta-role to have the AI write it for them.

**Content**: System prompt instructing the AI to:

- Understand role requirements from user descriptions
- Write effective system prompts
- Follow best practices for AI instruction design
- Output Markdown-formatted role documents
- Consider clarity, specificity, and actionability

**Related task**: `new-role` (see Tasks section below)

**Example usage**:

```bash
start task new-role "security auditor focused on OWASP top 10"
start task new-role "Go expert with focus on concurrency patterns"
```

---

## Tasks

### new-role

**Purpose**: Generate a new role document from a natural language description.

**Configuration**:

```toml
[tasks.new-role]
alias = "nr"
description = "Generate a new role document from description"
role = "role-writer"
prompt = """
Create a role document for: {instructions}

The role should be written in markdown and follow best practices for AI system prompts.
Include clear instructions, scope, and behavioral guidelines.
Output only the role content, ready to save to a file.
"""
```

**Benefits**:

- Uses existing task/role system (no new commands needed)
- Demonstrates system composability
- Users can customize the meta-role or task
- Can pipe output to file: `start task new-role "..." > roles/custom.md`

**Related role**: `meta/role-writer` (see Roles section above)

---

## Contexts

<!-- Add context asset ideas here -->

---

## Templates

<!-- Add template ideas here (if we have template assets) -->

---

## Notes

- All assets should follow the catalog structure and naming conventions
- Assets should be self-documenting with clear descriptions
- Consider dependencies between assets (e.g., tasks that require specific roles)
- Test assets work offline after initial download
