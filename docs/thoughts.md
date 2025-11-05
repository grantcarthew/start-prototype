# Design Thoughts

- I want it to be able to use a prompt writer prompt to create role documents on the fly

## `start update` Command Ideas

### Potential Subcommands/Features

**Health Check (`start update check` or just `start update`):**
- Current vs latest version (from GitHub releases)
- Config validation (both global and local)
- Agent binary availability (all configured agents)
- Context document existence (verify files referenced in config)
- Broken reference detection (missing files, invalid commands)

**Update Operations:**
- `start update self` - Self-update the binary (GitHub releases)
- `start update assets` - Refresh any bundled templates/examples
- `start update check-agents` - Verify agent binaries and suggest updates

**Diagnostic Output Example:**
```
start v1.2.3 (latest: v1.2.5 - update available)
✓ Global config valid (~/.config/start/config.toml)
✓ Local config valid (.start/config.toml)
✓ 3/3 agents available (claude, cursor, aider)
⚠ Context document missing: PROJECT.md (referenced in config)
✓ All commands valid
```

**Additional Ideas:**
- `start update migrate` - Migrate config from old format to new (if schema changes)
- Interactive mode to fix issues found during health check
- Dry-run flag to see what would be updated
- Check for deprecated config options
