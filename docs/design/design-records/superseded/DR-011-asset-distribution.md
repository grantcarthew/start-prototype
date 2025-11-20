# DR-011: Asset Distribution and Update System

- Date: 2025-01-03, Updated 2025-01-06, Superseded 2025-01-17
- Status: Superseded by DR-031 (Catalog-Based Assets)
- Category: Distribution

## Superseded

This design record is superseded by DR-031 (Catalog-Based Assets). The bulk download model is replaced by catalog-driven, on-demand asset loading.

The catalog system replaces the bulk download approach because:

- NO bulk downloads during `start init` - catalog queries GitHub on-demand
- NO `asset-version.toml` tracking file - filesystem IS the state
- Lazy loading: assets downloaded on first use via `start task <name>` or `start assets add`
- Per-asset updates: `start assets update` uses SHA comparison for individual cached assets
- Interactive discovery: `start assets add` (TUI browser) or `start assets search "query"`
- Cache is transparent: `~/.config/start/assets/` automatically populated as needed

The bulk download model (download entire library upfront) is incompatible with the catalog model (query and download on-demand).

## Asset Installation Location

```
~/.config/start/
├── config.toml              # User's global config
├── asset-version.toml       # Track asset library version
└── assets/                  # Downloaded asset library
    ├── agents/
    ├── roles/
    ├── tasks/
    └── examples/
```

## Distribution

- Assets stored in GitHub repository (`/assets` directory)
- Downloaded on-demand (not embedded in binary)
- Updateable without new release
- `start init` performs initial download
- `start assets update` refreshes asset library
- Network required for download (can work offline after initial setup)

## Asset Usage Patterns

**Agent templates:**

- Located in `~/.config/start/assets/agents/`
- Used during `start assets add` to pre-fill configurations
- User selects template, values are copied to `config.toml`

**Role files:**

- Located in `~/.config/start/assets/roles/`
- Referenced in config: `file = "~/.config/start/assets/roles/code-reviewer.md"`
- Updates flow automatically when `start assets update` is run

**Task definitions:**

- Located in `~/.config/start/assets/tasks/`
- Shown as templates in `start config task list`
- Users explicitly add them to config to activate

**Example configs:**

- Located in `~/.config/start/assets/examples/`
- Reference only, not automatically loaded
- Users manually copy sections to their config

## Rationale

- Assets updateable without binary release
- New agent configs, roles, tasks available immediately
- Network dependency acceptable (one-time per update)
- Offline work after initial download
- Separation: binary vs content
- Users control update timing (not forced)

## Related Decisions

- [DR-013](./dr-013-agent-templates.md) - Agent configuration distribution
- [DR-014](./dr-014-github-tree-api.md) - GitHub download strategy
- [DR-015](./dr-015-atomic-updates.md) - Atomic update mechanism
