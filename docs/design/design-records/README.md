# Design Record

This document indexes all design decisions for the `start` tool.

See [vision.md](../../vision.md) for the product vision and goals.

## Quick Reference

| DR | Title | Category | Date |
|----|-------|----------|------|
| [DR-001](./dr-001-toml-format.md) | TOML Configuration Format | Configuration | 2025-01-03 |
| [DR-002](./dr-002-config-merge.md) | Global + Local Config Merge | Configuration | 2025-01-03 |
| [DR-003](./dr-003-named-documents.md) | Named Documents for Context | Configuration | 2025-01-03 |
| [DR-004](./dr-004-agent-scope.md) | Agent Configuration Scope | Configuration | 2025-01-03 |
| [DR-005](./dr-005-role-configuration.md) | Role Configuration & Selection | Configuration | 2025-01-03 |
| [DR-006](./dr-006-cobra-cli.md) | CLI Command Structure (Cobra) | CLI Design | 2025-01-03 |
| [DR-007](./dr-007-placeholders.md) | Command Interpolation & Placeholders | Configuration | 2025-01-03 |
| [DR-008](./dr-008-file-handling.md) | Context File Detection & Handling | Runtime Behavior | 2025-01-03 |
| [DR-009](./dr-009-task-structure.md) | Task Structure, Agent Field & Placeholders | Tasks | 2025-01-03 |
| [DR-010](./dr-010-default-tasks.md) | Default Task Definitions | Tasks | 2025-01-03 |
| [DR-011](./dr-011-asset-distribution.md) | Asset Distribution & Update System | Distribution | 2025-01-03 |
| [DR-012](./dr-012-context-required.md) | Context Document Required Field | Configuration | 2025-01-04 |
| [DR-013](./dr-013-agent-templates.md) | Agent Templates from GitHub | Distribution | 2025-01-04 |
| [DR-014](./archive/dr-014-github-tree-api.md) | GitHub Tree API for Assets | Asset Management | 2025-01-06 |
| [DR-015](./archive/dr-015-atomic-updates.md) | Atomic Update Mechanism | Asset Management | 2025-01-06 |
| [DR-016](./archive/dr-016-asset-discovery.md) | Asset Discovery Strategy | Asset Management | 2025-01-06 |
| [DR-017](./dr-017-cli-reorganization.md) | CLI Command Reorganization | CLI Design | 2025-01-06 |
| [DR-018](./dr-018-init-update-integration.md) | Init/Update Command Integration | Asset Management | 2025-01-06 |
| [DR-019](./dr-019-task-loading.md) | Task Loading & Merging Algorithm | Tasks | 2025-01-06 |
| [DR-020](./dr-020-version-injection.md) | Binary Version Injection Strategy | Build & Distribution | 2025-01-06 |
| [DR-021](./dr-021-github-version-check.md) | GitHub Version Checking | Version Management | 2025-01-06 |
| [DR-022](./dr-022-asset-branch-strategy.md) | Asset Branch Strategy | Asset Management | 2025-01-06 |
| [DR-023](./archive/dr-023-asset-staleness-check.md) | Asset Staleness Checking | Asset Management | 2025-01-06 |
| [DR-024](./dr-024-doctor-exit-codes.md) | Doctor Exit Code System | CLI Design | 2025-01-06 |
| [DR-025](./dr-025-no-automatic-checks.md) | No Automatic Checks or Caching | CLI Design | 2025-01-06 |
| [DR-026](./dr-026-offline-behavior.md) | Offline Fallback & Network Unavailable | Asset Management | 2025-01-07 |
| [DR-027](./dr-027-security-trust-model.md) | Security & Trust Model for Assets | Asset Management | 2025-01-07 |
| [DR-028](./dr-028-shell-completion.md) | Shell Completion Support | CLI Design | 2025-01-07 |
| [DR-029](./dr-029-task-agent-field.md) | Task Agent Field | Tasks | 2025-01-07 |
| [DR-030](./dr-030-prefix-matching.md) | Prefix Matching for Commands | CLI Design | 2025-01-10 |
| [DR-031](./dr-031-catalog-based-assets.md) | Catalog-Based Asset Architecture | Asset Management | 2025-01-10 |
| [DR-032](./dr-032-asset-metadata-schema.md) | Asset Metadata Schema | Asset Management | 2025-01-10 |
| [DR-033](./dr-033-asset-resolution-algorithm.md) | Asset Resolution Algorithm | Asset Management | 2025-01-10 |
| [DR-034](./dr-034-github-catalog-api.md) | GitHub Catalog API Strategy | Asset Management | 2025-01-10 |
| [DR-035](./dr-035-interactive-browsing.md) | Interactive Asset Browsing | Asset Management | 2025-01-10 |
| [DR-036](./dr-036-cache-management.md) | Cache Management | Asset Management | 2025-01-10 |
| [DR-037](./dr-037-asset-updates.md) | Asset Update Mechanism | Asset Management | 2025-01-10 |

## By Category

### Configuration (DR-001 to DR-008, DR-012)

Core configuration structure and file handling:

- **[DR-001](./dr-001-toml-format.md)** - Use TOML for all configuration files
- **[DR-002](./dr-002-config-merge.md)** - Global + local config merge strategy
- **[DR-003](./dr-003-named-documents.md)** - Named document sections instead of arrays
- **[DR-004](./dr-004-agent-scope.md)** - Agents in both global and local configs
- **[DR-005](./dr-005-role-configuration.md)** - Role configuration with UTD pattern and selection precedence
- **[DR-007](./dr-007-placeholders.md)** - Single-brace placeholder system
- **[DR-008](./dr-008-file-handling.md)** - Relative paths and missing file handling
- **[DR-012](./dr-012-context-required.md)** - Required field and document order

### CLI Design (DR-006, DR-017, DR-024, DR-025, DR-028, DR-030)

Command-line interface structure:

- **[DR-006](./dr-006-cobra-cli.md)** - Cobra framework with subcommands and global flags
- **[DR-017](./dr-017-cli-reorganization.md)** - Configuration under `start config`, execution at top level
- **[DR-024](./dr-024-doctor-exit-codes.md)** - Simple binary exit codes (0 = healthy, 1 = issues)
- **[DR-025](./dr-025-no-automatic-checks.md)** - No automatic checks or result caching
- **[DR-028](./dr-028-shell-completion.md)** - Shell completion for bash/zsh/fish
- **[DR-030](./dr-030-prefix-matching.md)** - Enable prefix matching for all commands

### Tasks (DR-009, DR-010, DR-019, DR-029)

Task configuration and loading:

- **[DR-009](./dr-009-task-structure.md)** - Task structure with role/agent fields and {instructions}/{command} placeholders
- **[DR-010](./dr-010-default-tasks.md)** - Four interactive review tasks as defaults
- **[DR-019](./dr-019-task-loading.md)** - Task loading now includes cache resolution
- **[DR-029](./dr-029-task-agent-field.md)** - Optional agent field for task-specific agent preference

### Asset Management (DR-011, DR-013 to DR-016, DR-018, DR-022, DR-023, DR-026, DR-027, DR-031 to DR-037)

Asset distribution and updates:

- **[DR-011](./dr-011-asset-distribution.md)** - GitHub-fetched assets with update system
- **[DR-013](./dr-013-agent-templates.md)** - Fetch agent configs from GitHub
- **[DR-014](./archive/dr-014-github-tree-api.md)** - SHA-based caching for incremental updates (Superseded by DR-031)
- **[DR-015](./archive/dr-015-atomic-updates.md)** - Atomic install with rollback capability (Superseded by DR-031)
- **[DR-016](./archive/dr-016-asset-discovery.md)** - Each feature checks its own directory (Superseded by DR-031)
- **[DR-018](./dr-018-init-update-integration.md)** - Init and update share implementation
- **[DR-022](./dr-022-asset-branch-strategy.md)** - Assets from main branch (not releases)
- **[DR-023](./archive/dr-023-asset-staleness-check.md)** - GitHub commit comparison with no caching (Superseded by DR-031)
- **[DR-026](./dr-026-offline-behavior.md)** - Network-only, no manual installation, graceful degradation
- **[DR-027](./dr-027-security-trust-model.md)** - Trust GitHub HTTPS, no signatures, no pinning
- **[DR-031](./dr-031-catalog-based-assets.md)** - Catalog-Based Asset Architecture
- **[DR-032](./dr-032-asset-metadata-schema.md)** - Asset Metadata Schema
- **[DR-033](./dr-033-asset-resolution-algorithm.md)** - Asset Resolution Algorithm
- **[DR-034](./dr-034-github-catalog-api.md)** - GitHub Catalog API Strategy
- **[DR-035](./dr-035-interactive-browsing.md)** - Interactive Asset Browsing
- **[DR-036](./dr-036-cache-management.md)** - Cache Management
- **[DR-037](./dr-037-asset-updates.md)** - Asset Update Mechanism

### Build & Distribution (DR-011, DR-013, DR-020)

How the tool and its content are distributed:

- **[DR-011](./dr-011-asset-distribution.md)** - Assets fetched from GitHub, updateable without releases
- **[DR-013](./dr-013-agent-templates.md)** - Agent configurations as GitHub content
- **[DR-020](./dr-020-version-injection.md)** - Binary version via ldflags at build time

### Runtime Behavior (DR-008)

How the tool behaves during execution:

- **[DR-008](./dr-008-file-handling.md)** - Path resolution and missing file handling

### Version Management (DR-020, DR-021)

Version tracking and update checking:

- **[DR-020](./dr-020-version-injection.md)** - Binary version via ldflags at build time
- **[DR-021](./dr-021-github-version-check.md)** - GitHub Releases API check with no caching

## Key Patterns

### Unified Template Design (UTD)

The `file`, `command`, `prompt` pattern is used consistently across:
- Context documents
- Roles
- Task configurations

See individual DRs for details.

### Replacement vs Merge

- **Merge:** Contexts (DR-003), Agents (DR-004)
- **Replace:** Roles (DR-005), Tasks (DR-019)

### Local Precedence

When both global and local configs define the same item, local wins:
- Contexts: Combined (global + local), local overrides by name
- Agents: Combined (global + local), local overrides by name
- Roles: Local completely replaces global (for same role name)
- Tasks: Local completely replaces global (for same task name)

## References

- [Vision](../../vision.md) - Product vision and goals
- [Config Reference](../../config.md) - Complete configuration specification
- [Unified Template Design](../unified-template-design.md) - UTD pattern details
- [Tasks](../../tasks.md) - Task configuration details

## Contributing

When adding a new design decision:

1. Create a new file: `dr-XXX-short-name.md`
2. Use the template format (see existing DRs)
3. Update this index with the new DR in the Quick Reference table
4. Add to the appropriate category section
5. Cross-reference related DRs
