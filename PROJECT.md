# Design Thoughts

This document is just my thoughts dumped for reference. Do not use this document as a concrete reference. It is most likely wrong.

## General Ideas

- ~~I want it to be able to use a prompt writer prompt to create role documents on the fly~~ → **RESOLVED**: Added to `docs/ideas/assets.md` as catalog assets (`roles/meta/role-writer.md` + `tasks/new-role.toml`)
- ~~Need an easy way to switch defaults~~ → **RESOLVED**: Already exists via `start config role default <name>` and `start config agent default <name>`
- ~~Need a config delete option if it does not exist (or remove), something like `start config agent rm xyz`~~ → **RESOLVED**: Already exists - `start config <type> remove` commands in design
- ~~**Dry-run flag**: Add `--dry-run` flag to preview aggregated context without calling the agent~~ → **RESOLVED**: Evolved into `start show` command - see `docs/cli/start-show.md`
  - Execution preview: `start show`, `start show task <name>`, `start show prompt <text>`
  - Content viewer: `start show role`, `start show context`, `start show agent`, `start show task`
- **Unified asset management**: Consider `start assets` subcommand to consolidate asset operations
  - `start assets browse` - Open GitHub catalog in browser (better discoverability than `start config <type> add`)
  - `start assets add <type> <name>` - Add from catalog to config (replaces `start config <type> add`)
  - `start assets update [name]` - Replace vague `start update` with clearer naming; optionally update specific asset
  - `start assets info <name>` - Show asset metadata (description, last updated, dependencies, source)
  - `start assets list` - Show all cached assets with status
  - `start assets clean` - Clear cache (vs manual rm -rf)
  - **Problem:** Current asset management fragmented across `start config <type> add` commands, not discoverable
  - **Benefit:** Single place for all asset operations, clearer commands, better UX
  - **Semantic separation:**
    - `start config` = Manage YOUR configuration (things you've defined/customized)
      - `new` - Create new custom asset
      - `edit` - Edit your asset
      - `remove` - Remove your asset
      - `test` - Test your asset
      - `list` - List your assets
    - `start assets` = Interact with the CATALOG (browse, add from GitHub, update cache)
      - `browse` - View catalog in browser
      - `add` - Add from catalog to config
      - `update` - Update cached assets
      - `info` - Show asset metadata
      - `list` - Show cached assets
  - **Migration:** `start config task add` → `start assets add task`
