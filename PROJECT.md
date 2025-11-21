# Project: CLI Documentation Review and Correction

**Status:** Restarting Review
**Date Started:** 2025-01-11
**Previous Phase:** Catalog Design (completed, see PROJECT-backlog.md)

## Overview

Comprehensive review and correction of all CLI documentation in `docs/cli/` to ensure accuracy, consistency, and alignment with design decisions before implementation begins.

## Important

I am reviewing the documents. Your task is to fix issues defined by me during this review process. It is an interactive session.

We will be working through this step by step. If you see something that should be addressed, bring it up in the discussion.

As we fix issues, I will be needing you to ensure the docs/design/design-records/\* document that relates to the changes is in-sync with our updates.

This repository is not large yet. Before you do anything, run a `lsd --tree` to get a list of the documents.

REPEAT: This is an interactive session. Do not change anything without approval from me!

## Objectives

1. Review all documentation files (3 root docs + 17 CLI commands) for accuracy
2. Fix inconsistencies, outdated information, and errors
3. Ensure alignment with design decisions (DRs)
4. Verify examples, flags, and usage patterns are correct
5. Prepare clean, accurate documentation for implementation phase
6. Design records will be updated as needed based on changes to other documents

## Review Process

1. I will identify and issue
2. We will discuss it and decide on the fix
3. You will fix the issue in the document
4. You will review related documents in docs/cli/ and docs/design/design-records/ to make sure there are no inconsistencies
5. You will update related documents
6. Add the fix to the bottom of this document in the ## Fixed section
7. Ask me to commit the changes
8. Next issue

## Documents to Review

Active document: `docs/cli/start-config-context.md`

### Root Documentation

- [x] `docs/config.md` - Configuration reference
- [x] `docs/tasks.md` - Task-specific documentation
- [x] `docs/vision.md` - Product vision and goals

### Main Commands

- [x] `docs/cli/start.md` - Main entry point, interactive sessions
- [x] `docs/cli/start-prompt.md` - Prompt composition and execution
- [x] `docs/cli/start-task.md` - Task execution

### Asset Commands

- [x] `docs/cli/start-assets.md` - Asset management overview
- [x] `docs/cli/start-assets-add.md` - Add assets from catalog
- [x] `docs/cli/start-assets-browse.md` - Browse catalog in browser
- [x] `docs/cli/start-assets-search.md` - Search catalog
- [x] `docs/cli/start-assets-info.md` - Show asset information

### Configuration Commands

- [x] `docs/cli/start-config.md` - Configuration management overview
- [x] `docs/cli/start-config-agent.md` - Agent configuration
- [x] `docs/cli/start-config-context.md` - Context configuration
- [ ] `docs/cli/start-config-role.md` - Role configuration
- [ ] `docs/cli/start-config-task.md` - Task configuration

## Design Alignment

Ensure all documentation aligns with these key design decisions:

- **DR-031**: Catalog-based asset architecture
- **DR-032**: Asset metadata schema (.meta.toml files)
- **DR-033**: Asset resolution (local → global → cache → GitHub)
- **DR-034**: GitHub API strategy (Tree API + raw.githubusercontent.com)
- **DR-035**: Interactive browsing (numbered selection)
- **DR-036**: Cache management (invisible, manual delete only)
- **DR-037**: Update mechanism (manual, SHA-based)

## Success Criteria

Documentation review is complete when:

- [ ] All 20 documents (3 root + 17 CLI) reviewed and corrected
- [ ] Examples are accurate and tested conceptually
- [ ] Flags and options are consistent across commands
- [ ] Multi-file config structure correctly documented
- [ ] Catalog behavior (lazy loading, browsing) accurately described
- [ ] Design decisions referenced where relevant
- [ ] No contradictions between documents
- [ ] Design records updated to reflect changes
- [ ] Ready for implementation phase

## Notes

- This review may identify gaps requiring new design decisions
- Some issues may require updates to design documents (docs/design/)
- Focus on correctness over completeness - better to have accurate docs than comprehensive but wrong docs
- Track design questions separately for resolution before implementation

## Fixed

- `docs/cli/start.md`: Removed "alias" terminology from --model flag, clarified resolution as exact match → prefix match → passthrough
- `docs/cli/start-prompt.md`: Updated --model flag to match start.md (removed "alias|name", added resolution order)
- **Short flags added**: Added `-a` (--agent), `-r` (--role), `-m` (--model) short flags across all CLI docs
- **Version flag corrected**: Changed `-v` from --verbose to --version across all CLI docs (--verbose has no short form)
- `docs/cli/start.md`: Added short flags -a, -r, -m; moved -v to --version; removed -v from --verbose
- `docs/cli/start-prompt.md`: Added short flags -a, -r, -m; added --version with -v
- `docs/cli/start-config.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-agent.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-task.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-role.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-context.md`: Moved -v to --version; removed from --verbose
- `docs/design/design-records/dr-006-cobra-cli.md`: Updated Global Flags section with complete flag list including short forms
- `docs/design/design-records/dr-028-shell-completion.md`: Updated flag completion examples to include all short flags
- **--context flag**: Decided not to add --context/-c flag; current config-driven design is sufficient
- **Missing file behavior**: Standardized wording across all docs - missing context files always generate warnings and are skipped
- `docs/cli/start.md`: Updated context behavior description and execution flow; changed output examples from ✗ to ⚠ for missing files
- `docs/cli/start-config.md`: Changed from "silently skipped" to "generate warnings and are skipped"
- `docs/cli/start-config-context.md`: Updated test output example for missing files
- `docs/cli/start-task.md`: Updated execution flow for missing file handling
- **Runtime behavior**: Missing files show `⚠ context-name file-path (not found, skipped)` - warnings, not errors
- `docs/design/design-records/dr-008-file-handling.md`: Updated to reflect warning behavior instead of "silently skipped"; updated output examples to use ⚠ symbol; added rationale about catching config errors
- **Flag value prefix matching**: Implemented intelligent prefix matching for --agent, --role, --task, and --model flags with ambiguity detection and TTY-aware interactive selection
- `docs/design/design-records/dr-038-flag-value-resolution.md`: Created new DR defining two-phase resolution (exact → prefix), short-circuit evaluation, ambiguity handling (interactive/error), and passthrough for --model
- `docs/cli/start.md`: Updated --agent, --role, and --model flag descriptions with full resolution algorithm and examples
- `docs/cli/start-prompt.md`: Updated --agent, --role, and --model flag descriptions to reference DR-038
- `docs/cli/start-task.md`: Updated task resolution section and name argument to describe prefix matching behavior
- `docs/design/design-records/dr-033-asset-resolution-algorithm.md`: Added DR-038 to related decisions (prefix matching extends exact match)
- **UTD Placeholder Design**: Enhanced placeholder system with four new placeholders for better flexibility
- New placeholders: `{file}` (path), `{file_contents}` (contents), `{command}` (string), `{command_output}` (output)
- Pattern: Short form = reference/source, long form = result/contents
- `docs/design/unified-template-design.md`: Updated all placeholder definitions and examples with new four-placeholder system
- `docs/design/design-records/dr-007-placeholders.md`: Updated UTD Pattern Placeholders section with detailed descriptions and use cases for all four placeholders
- `docs/config.md`: Updated all UTD field descriptions and examples throughout roles, contexts, and tasks sections
- `docs/tasks.md`: Updated task prompt placeholders and all code examples
- `docs/cli/start-prompt.md`: Examples already correct (using {file} for paths)
- `docs/cli/start-task.md`: Updated execution flow and placeholder documentation
- `docs/cli/start-config-context.md`: Updated field descriptions and placeholder section
- `docs/cli/start-config-role.md`: Updated field descriptions and placeholder section
- `docs/cli/start-config-task.md`: Updated field descriptions and placeholder section
- `docs/design/design-records/dr-009-task-structure.md`: Updated decision summary, field descriptions, placeholder definitions, execution flow, and all examples
- **Contexts as lazy-loadable assets**: Confirmed contexts can be lazy-loaded from GitHub catalog (config templates, not content)
- `docs/design/design-records/dr-031-catalog-based-assets.md`: Added `contexts/` to cache structure; removed contexts from Future Considerations
- `docs/config.md`: Updated Asset Resolution & Lazy-Loading section to include contexts as downloadable asset type
- `docs/design/design-records/dr-032-asset-metadata-schema.md`: Added contexts to asset types list
- `docs/design/design-records/dr-022-asset-branch-strategy.md`: Added contexts to asset types in three locations (content list, benefits, rationale)
- **Command name corrections**: Fixed all instances of `start agent` to `start config agent` (there is no `start agent` command)
- `docs/config.md`: Changed `start agent list` to `start config agent list`
- `docs/cli/start.md`: Changed `start agent list` to `start config agent list`
- `docs/cli/start-doctor.md`: Changed `start agent remove` to `start config agent remove` and `start agent test` to `start config agent test`
- `docs/cli/start-config-agent.md`: Changed `start agent list` to `start config agent list` and two instances of `start agent remove` to `start config agent remove`
- **Local-only config support**: Fixed documentation to reflect that local config can be used without global config
- `docs/cli/start-config.md`: Changed "Agents: (none - agents must be defined in global config)" to "Agents: (none configured)" (line 178)
- `docs/cli/start-config.md`: Fixed Configuration Merge Behavior section - agents, roles, and tasks can all be in both global and local configs with merge behavior (lines 889-893)
- `docs/cli/start.md`: Fixed "Local config only" section - removed error message, clarified that local-only config is valid (no global config required)
- **Man page references**: Fixed all `start-agent(1)` references to `start-config-agent(1)` in See Also sections
- `docs/cli/start-doctor.md`: Changed `start-agent(1)` to `start-config-agent(1)`
- `docs/cli/start-init.md`: Changed `start-agent(1)` to `start-config-agent(1)`
- `docs/cli/start-config.md`: Changed `start-agent(1)` to `start-config-agent(1)`
- `docs/cli/start.md`: Changed `start-agent(1)` to `start-config-agent(1)`
- `docs/cli/start-prompt.md`: Changed `start-agent(1)` to `start-config-agent(1)`
- **TOML link added**: Added link to <https://toml.io/> in Overview section of docs/config.md
- **Settings merge clarification**: Changed "Settings: Local values override global values" to "Settings: Merged per-field, local overrides global for same field" to clarify it's per-field merge, not whole-section replacement
- `docs/config.md`: Updated Merge behavior section (line 23)
- `docs/cli/start-config.md`: Updated Configuration Merge Behavior section (line 889)
- **Path choice rationale**: Added explanation to DR-002 for why `./.start/` is used (not `./.config/start/`)
- `docs/design/design-records/dr-002-config-merge.md`: Added "Path choices" section to Rationale explaining `./.start/` follows project-level tool convention (like `.vscode/`, `.github/`), while `.config/` is for user-level configs
- **Missing settings documentation**: Added complete documentation for all settings fields in [settings] section
- `docs/config.md`: Added field documentation for `default_role` (uses first role if not specified)
- `docs/config.md`: Added field documentation for `asset_download` (enable/disable auto-download from GitHub, default true)
- `docs/config.md`: Added field documentation for `asset_repo` (GitHub repository for assets, default "grantcarthew/start")
- `docs/config.md`: Added field documentation for `asset_path` (cache directory, default "~/.config/start/assets")
- `docs/config.md`: Updated Validation section to include all settings
- `docs/config.md`: Updated Example section to show all settings
- `docs/config.md`: Updated Complete Example to show all settings in global config
- **Settings field name correction**: Fixed incorrect field name "verbosity" to "log_level" in examples
- `docs/cli/start-config.md`: Changed three instances of `verbosity` to `log_level` (lines 901, 913, 924)
- **Incomplete settings blocks**: Added missing `asset_download` field to complete settings examples
- `docs/ideas/catalog-based-assets.md`: Added `asset_download = true` to settings block (line 358)
- `PROJECT-backlog.md`: Added `asset_download = true` to settings block (line 259)
- `docs/design/design-records/dr-031-catalog-based-assets.md`: Added missing basic settings (`log_level`, `shell`, `command_timeout`) to complete the example (lines 80-82)
- **New command: start config role default**: Added documentation for setting/showing default role
- `docs/cli/start-config-role.md`: Added `start config role default [name]` to Synopsis (line 16)
- `docs/cli/start-config-role.md`: Added "default" to role management operations list (line 31)
- `docs/cli/start-config-role.md`: Added complete section documenting the default subcommand with examples (lines 904-1025)
- **Asset resolution for defaults**: Fixed default_agent and default_role to reflect lazy loading (assets don't need to be pre-defined in config sections)
- `docs/config.md`: Changed default_agent from "must match agent defined in [agents] section" to "resolved using asset resolution algorithm" (line 238)
- `docs/config.md`: Changed default_role from "must match role defined in [roles] section" to "resolved using asset resolution algorithm" (line 246)
- `docs/tasks.md`: Updated agent field description to reflect asset resolution algorithm (line 75)
- `docs/tasks.md`: Updated agent validation to reflect lazy loading from catalog (lines 84-86)
- `docs/tasks.md`: Updated role field description to reflect asset resolution algorithm (line 91)
- `docs/tasks.md`: Updated role validation to reflect lazy loading from catalog (lines 107-109)
- `docs/cli/start-config-task.md`: Changed role field from "defined in [roles.<name>] section" to "resolved via asset resolution algorithm" (line 75)
- **Removed github_token_env setting**: Removed unnecessary github_token_env setting - hardcoded to use GITHUB_TOKEN env var (industry standard)
- `docs/config.md`: Removed github_token_env field documentation and all examples
- `docs/design/design-records/dr-031-catalog-based-assets.md`: Removed github_token_env from settings example
- `docs/design/design-records/dr-034-github-catalog-api.md`: Removed github_token_env setting, added note about hardcoded GITHUB_TOKEN env var
- `docs/cli/start-assets-update.md`: Removed github_token_env from settings example, added note about GITHUB_TOKEN
- `docs/ideas/catalog-based-assets.md`: Removed github_token_env from settings, updated GitHub Token section to clarify hardcoded env var
- `PROJECT-backlog.md`: Removed github_token_env from settings example
- **Agent {prompt} placeholder validation**: Changed from required (error) to recommended (warning) - "works but warns" design
- `docs/config.md`: Changed command field description from "Must contain" to "Should contain {prompt} placeholder" (line 349)
- `docs/config.md`: Updated validation from error to warning: "Command doesn't contain {prompt} - composed prompt won't be passed to agent" (line 390)
- `docs/config.md`: Updated Validation Rules section - changed to "should contain {prompt} placeholder (warns if missing)" (line 1005)
- `docs/cli/start-config-agent.md`: Changed command field description to "Should contain {prompt} placeholder" (line 55)
- `docs/cli/start-config-agent.md`: Updated validation checklist to "Should contain (warns if missing)" (line 225)
- `docs/cli/start-config-agent.md`: Changed add command example from error to warning with Continue prompt (lines 366-375)
- `docs/cli/start-config-agent.md`: Updated test validation description to "checked for {prompt} (warns if missing)" (line 423)
- `docs/cli/start-config-agent.md`: Removed "required" from all test output examples (lines 447, 467, 490)
- `docs/cli/start-config-agent.md`: Changed test error output to warning (⚠) for missing {prompt} (line 509, 593)
- `docs/cli/start-config-agent.md`: Updated edit validation note to "should contain (warns if missing)" (line 653)
- `docs/cli/start-config-agent.md`: Changed edit command example from error to warning with Continue prompt (lines 769-778)
- `docs/cli/start-config-agent.md`: Removed "(required)" from verbose placeholder analysis output (line 542)
- `docs/cli/start-config.md`: Changed edit output from error (✗) to warning (⚠) for missing {prompt} (line 333)
- `docs/cli/start-config.md`: Changed validate output - agent section shows warnings instead of errors for missing {prompt} (lines 641-644)
- **Model terminology: alias → name**: Renamed "model alias" to "model name" throughout documentation for clarity
- `docs/config.md`: Changed "Model alias" to "Model name" in default_model field and [agents.<name>.models] section (user already updated)
- `docs/cli/start.md`: Changed "Model alias:" to "Model name:" in output example (line 220)
- `docs/cli/start-config-agent.md`: Changed "Model alias" to "Model name" in field descriptions (lines 67, 70)
- `docs/cli/start-config-agent.md`: Changed all "Model alias:" prompts to "Model name:" in examples (5 occurrences)
- `docs/cli/start-config-agent.md`: Changed "full name and alias" to "full name and model name" (lines 94-95)
- `docs/cli/start-config-agent.md`: Changed "alias name + full model name" to "model name + full model identifier" (line 232)
- `docs/cli/start-config-agent.md`: Changed "Model aliases defined" to "Model names defined" in validation (line 425)
- `docs/cli/start-config-agent.md`: Changed "model aliases" to "model names" in preferences (line 1246)
- `docs/cli/start-config-agent.md`: Renamed section "Model Aliases" to "Model Names" and updated content (lines 1280-1287)
- `docs/cli/start-config.md`: Changed "Model alias validation" to "Model name validation" (line 700)
- `docs/cli/start-init.md`: Changed "model aliases" to "model names" in three locations (lines 205, 224, 642, 648)
- `docs/design/design-records/dr-004-agent-scope.md`: Changed "alias" to "name" throughout Model Alias Behavior section (lines 19, 27-32, 38, 76)
- `docs/design/design-records/dr-006-cobra-cli.md`: Changed --model flag description from "Model alias or full model name" to "Model name or full model identifier" (line 48)
- `docs/design/design-records/dr-007-placeholders.md`: Changed {model} description from "after alias resolution" to "after name resolution" (line 16)
- **Model identifier terminology**: Standardized "full model name" to "full model identifier" for clarity
- `docs/cli/start-config-agent.md`: Changed "full name and model name" to "full identifier and model name" (lines 94-95)
- `docs/cli/start-config-agent.md`: Changed all "Full model name:" prompts to "Full model identifier:" in examples (3 occurrences)
- `docs/cli/start-config-agent.md`: Changed validation text from "full model name" to "full model identifier" (line 798)
- `docs/cli/start-config-agent.md`: Changed "alias name" to "model name" and "aliases" to "model names" in validation section (lines 797, 799)
- `docs/cli/start-config-agent.md`: Changed `--model <full-name>` to `--model <full-identifier>` (line 1264)
- `docs/cli/start-prompt.md`: Changed "Use full model name:" to "Use full model identifier:" (line 187)
- **Final alias cleanup**: Fixed remaining model alias references found via comprehensive search
- `docs/design/design-records/dr-004-agent-scope.md`: Changed historical update note "user-defined aliases" to "user-defined model names" (line 69)
- `docs/cli/start.md`: Changed debug output "Default alias:" to "Default model name:" (line 271)
- **{role_file} placeholder clarity**: Enhanced documentation to clearly explain temp file behavior for UTD roles
- `docs/config.md`: Expanded {role_file} description to explain simple vs UTD role behavior and temp file creation/cleanup (lines 904-907)
- `docs/cli/start-config-agent.md`: Clarified {role_file} description to mention temp file for UTD roles (line 1297)
- **Markdown rendering fix**: Fixed broken GitHub Markdown rendering caused by incorrect backtick nesting in toml blocks containing nested code blocks
- `docs/config.md`: Changed outer code fence to 5 backticks for two toml blocks with nested diff blocks (lines 799-813, 831-853)
- `docs/cli/start-task.md`: Fixed line 224 from 4 to 3 backticks (closing diff block) and upgraded outer fence to 5 backticks (lines 206-232)
- `docs/cli/start-task.md`: Fixed line 254 from 4 to 3 backticks (closing bash block)
- `docs/cli/start-config-task.md`: Upgraded outer fence to 5 backticks for two toml blocks with nested diff blocks (lines 36-61, 420-442)
- `docs/design/design-records/dr-009-task-structure.md`: Upgraded outer fence to 5 backticks for two toml blocks with nested diff blocks (lines 13-36, 293-310)
- `docs/design/design-records/dr-010-default-tasks.md`: Upgraded outer fence to 5 backticks for toml block with nested diff block (lines 49-75)
- **Removed "documents array" references**: Cleaned up outdated references to removed feature across all documentation
- `docs/config.md`: Removed "There is no `documents` array" from Context Inclusion section (line 825)
- `docs/tasks.md`: Changed "There is no `documents` array in task configuration. Instead:" to "Context inclusion behavior:" (line 169)
- `docs/cli/start-config-task.md`: Removed "No `documents` array needed" from 3 locations (lines 92, 1375, 1390)
- `docs/cli/start-task.md`: Removed "No `documents` array needed" from 3 locations (lines 175, 243, 641)
- `docs/design/design-records/dr-009-task-structure.md`: Removed "There is **no `documents` array**" statement and comment; reframed "Why no documents array:" to "Why automatic context inclusion:" (lines 105, 129, 445)
- **Fixed incorrect file path examples**: Corrected {file} placeholder examples to show proper path expansion
- `docs/config.md`: Changed `/Users/username/ref/ENV.md` to `/Users/username/reference/ENVIRONMENT.md` (line 924)
- `docs/design/design-records/dr-007-placeholders.md`: Changed example from `~/ref/ENV.md` to `~/reference/ENVIRONMENT.md` and output from `/Users/username/ref/ENV.md` to `/Users/username/reference/ENVIRONMENT.md` (line 90)
- **Fixed configuration file references**: Corrected all references to use specific config files (tasks.toml, agents.toml, etc.) instead of generic config.toml
- `docs/tasks.md`: Changed global/local tasks paths from config.toml to tasks.toml (lines 22-23, 198, 202)
- `docs/tasks.md`: Updated code example comments to use tasks.toml (lines 355, 413, 419)
- `docs/config.md`: Changed agent scope documentation from config.toml to agents.toml (lines 401-402)
- `docs/cli/start-config-task.md`: Updated task source references from config.toml to tasks.toml in list output examples (lines 148, 155, 162)
- `docs/cli/start-config-task.md`: Updated task source labeling from config.toml to tasks.toml (lines 1330-1331)
- `docs/cli/start-config-task.md`: Updated test output config file reference from config.toml to tasks.toml (line 713)
- `docs/cli/start-config.md`: Changed agent configuration scope from config.toml to agents.toml (lines 934-935)
- `docs/cli/start-config.md`: Fixed validation output examples to show multi-file structure instead of single config.toml (lines 604-696)
- `docs/cli/start.md`: Updated execution flow to reference config directories instead of config.toml files (lines 110-111)
- `docs/cli/start.md`: Updated verbose output to show multi-file loading (lines 214-215)
- `docs/cli/start.md`: Updated debug output to show individual config files being loaded (lines 258-265)
- `docs/cli/start.md`: Updated Files section to describe config directories instead of individual config.toml files (lines 460-464)
- `docs/cli/start.md`: Updated error handling to reference config directories (lines 559, 569)
- `docs/cli/start-task.md`: Updated verbose output to show multi-file loading (lines 394-395)
- `docs/cli/start-task.md`: Updated task discovery to reference tasks.toml instead of config.toml (lines 724, 731-732)
- `docs/cli/start-config-context.md`: Changed all context management references from config.toml to contexts.toml (lines 20, 246-247, 318, 364, 404, 714, 897, 899, 938, 1026, 1047, 1086, 1219, 1222, 1258, 1263)
- `docs/cli/start-config-role.md`: Changed --local flag description from config.toml to roles.toml (line 153)
- `docs/cli/start-config-agent.md`: Changed agent test output config file from config.toml to agents.toml (line 533)
- `docs/cli/start-config-agent.md`: Changed agent list output config file from config.toml to agents.toml (line 1187)
- `docs/cli/start-config-agent.md`: Changed local agents reference from config.toml to agents.toml (line 1248)
- `docs/cli/start-config-agent.md`: Changed invalid syntax error message from config.toml to agents.toml (line 1229)
- `docs/cli/start-config-context.md`: Changed backup failure error messages from config.toml to contexts.toml (lines 521, 1139)
- `docs/cli/start-config-context.md`: Changed invalid syntax error message from config.toml to contexts.toml (line 1244)
- `docs/cli/start-config-task.md`: Changed invalid syntax error message from config.toml to tasks.toml (line 1300)
- `docs/cli/start-config-role.md`: Changed TOML comment from "config.toml (or tasks.toml...)" to "roles.toml" (line 237)
- **Additional config.toml reference corrections**: Fixed remaining single-file config references to reflect multi-file structure
- `docs/design/unified-template-design.md`: Changed security warning from `./.start/config.toml` to `./.start/` (config files can execute commands, not just config.toml) (line 546)
- `docs/cli/start-doctor.md`: Changed simple output from individual config.toml files to config directories (lines 87-88)
- `docs/cli/start-doctor.md`: Updated verbose output to show multi-file config validation with all 5 files (lines 192-207)
- `docs/design/design-records/dr-019-task-loading.md`: Changed task source paths from config.toml to tasks.toml (lines 16-17)
- `docs/ideas/catalog-based-assets.md`: Changed resolution order from config.toml to tasks.toml for task resolution (lines 339-340)
- **Fixed nested code fence in tasks.md**: Upgraded outer fence to 5 backticks for toml block containing nested diff block (broke GitHub Markdown rendering)
- `docs/tasks.md`: Changed outer fence from 3 to 5 backticks for "Task with Role Reference" example (lines 255-273)
- **Corrected "Default Tasks" section**: Changed terminology and explanation to reflect catalog-based lazy-loading architecture
- `docs/tasks.md`: Renamed "Default Tasks" to "Available Tasks"; removed incorrect "embedded in binary" statement; explained GitHub catalog auto-download mechanism (lines 429-450)
- **Added missing --local flag documentation**: Added --local flag to start and start-prompt commands (controls where downloaded assets are added)
- `docs/cli/start.md`: Added --local flag description after --asset-download (lines 100-101)
- `docs/cli/start-prompt.md`: Added --asset-download and --local flag descriptions (lines 67-71)
- **Clarified task listing behavior**: Fixed misleading "available tasks" wording; clarified `start task` lists only configured tasks, not GitHub catalog
- `docs/cli/start-task.md`: Changed "all available tasks (default and custom)" to "all configured tasks" with pointer to catalog browsing (line 148)
- **Fixed "Default Tasks" in start-task.md**: Applied same fix as tasks.md - changed to "Available Tasks" with catalog explanation
- `docs/cli/start-task.md`: Renamed "Default Tasks" to "Available Tasks"; removed "ships with" language; explained GitHub catalog auto-download mechanism (lines 711-733)
- **Removed remaining "Default tasks" reference**: Removed misleading bullet point about default tasks in global config
- `docs/tasks.md`: Removed "Default tasks (cr, gdr, ct, dr)" from Global tasks description (line 200)
- **Replaced start init [scope] with --local flag**: Changed positional scope argument to --local flag for consistency with rest of CLI
- `docs/cli/start-init.md`: Changed synopsis from `start init [scope]` to `start init [flags]`; replaced location choice with `--local` flag; updated all examples and behavior documentation
- `docs/cli/start-config.md`: Changed `start init local` to `start init --local` (line 483)
- `docs/design/design-records/dr-026-offline-behavior.md`: Changed `start init global` to `start init` (line 65)
- `docs/design/design-records/dr-017-cli-reorganization.md`: Updated synopsis to show `[--local]` flag (line 43)
- **Enhanced init command with interactive/automatic modes**: Added three interaction levels controlled by --local and --force flags
- `docs/cli/start-init.md`: Added interactive mode (asks location), partially interactive (--local skips location), and fully automatic (--force skips all prompts)
- **Added bin field to agent TOML structure**: New required field for auto-detection and DRY command templates
- `docs/design/design-records/dr-013-agent-templates.md`: Added bin field specification, {bin} placeholder usage, and updated init behavior to use index.csv
- `docs/design/design-records/dr-007-placeholders.md`: Added {bin} placeholder to global placeholders; updated all agent examples to include bin field
- `docs/design/design-records/dr-039-catalog-index.md`: Added `bin` column to index.csv schema; updated generation algorithm to extract bin from agent TOML files
- `docs/config.md`: Added bin field documentation to [agents.<name>] section; updated all agent examples; updated validation rules for required {bin} and {model} placeholders
- `docs/cli/start-init.md`: Updated agent detection to use index.csv → binary detection → lazy TOML download workflow; updated GitHub Catalog Details section
- **Validation changes**: Made {bin} and {model} placeholders required (errors), {prompt} remains recommended (warns)
- **Deprecated command cleanup in start-init.md**: Fixed 2 references from `start config agent add` to `start assets add` (lines 254, 831)
- **Removed shell-specific command -v references**: Made documentation implementation-agnostic for Go implementation
- `docs/cli/start-init.md`: Changed "command -v" to "checking if bin is executable" throughout agent detection sections
- `docs/config.md`: Changed bin field from "auto-detection via command -v" to "auto-detection in PATH" → "Binary path or name to execute"
- `docs/design/design-records/dr-013-agent-templates.md`: Changed "command -v detection" to "Binary path or name for agent detection"
- `docs/design/design-records/dr-024-doctor-exit-codes.md`: Changed "Binary not found in PATH" to "Binary not found" throughout
- `docs/cli/start.md`: Changed "in your PATH" to "installed and available"
- `docs/cli/start-config.md`: Changed "Binary not found in PATH" to "Binary not found"
- `docs/cli/start-config-agent.md`: Changed "not in your PATH" to "not available"
- `docs/cli/start-config-task.md`: Changed error message from "not found in PATH" to "not found"
- `PROJECT.md`: Changed changelog "command -v" to "binary detection"
- **Binary detection language made implementation-agnostic**: Removed PATH-only assumptions; bin field can be name, relative path, or absolute path
- `docs/config.md`: Updated bin field to document support for binary name, relative path, or absolute path
- `docs/cli/start-init.md`: Changed "only finds binaries in PATH" to "only finds executables that are discoverable"
- `docs/design/design-records/dr-024-doctor-exit-codes.md`: Changed "binary in PATH" to "binary is discoverable"
- **Agent asset file correction**: Fixed documentation showing agents have 3 files; agents only have 2 (.toml and .meta.toml, no .md file)
- `docs/cli/start-init.md`: Corrected agent asset downloads from 3 files to 2 files (removed .md reference) - agents don't use UTD pattern
- **Updated PROJECT.md command structure**: Removed obsolete `start-update.md` reference; added 7 asset commands (start-assets.md and subcommands); updated document count from 14 to 20
- **Fixed all `start update` to `start assets update` references**: Updated 18 design records and 4 CLI docs to use current command name
- Living (non-archived) design records updated: dr-010, dr-011, dr-018, dr-021, dr-022, dr-024, dr-025, dr-026, dr-027, dr-030, dr-031, dr-034, dr-036, dr-037
- CLI docs updated: start-config-task.md, start-doctor.md, start-config-role.md
- Preserved historical references in archive files, DR-017 (deprecation record), and DR-041 (migration record)
- **Asset resolution clarification**: Updated `docs/config.md` to clarify that assets are lazy-loaded to the global cache (`~/.config/start/assets/{type}/`) and added to global config by default, or local config with `-l`/`--local`.
- **Path resolution simplification**: Updated `docs/config.md` to align with DR-008 - relative paths always resolve to working directory, regardless of config scope.
- **Task context inclusion**: Updated `docs/tasks.md` to generalize "required context" description.
- **Task lazy-loading details**: Updated `docs/tasks.md` to explicitly state tasks download to global cache and add to global/local config based on flags.
- **Lazy-loading clarification**: Updated `docs/cli/start.md` to include contexts in asset types, simplify cache path, and clarify global vs local configuration targeting with `--local`.
- **Local flag clarification**: Updated `docs/cli/start-prompt.md` to explicitly state that assets download to global cache but config entry goes to local scope when `--local` is used.
- **Task discovery update**: Removed "embedded in binary" tasks reference; updated task discovery to include cache and GitHub catalog with correct precedence order (Local → Global → Cache → GitHub)
- `docs/cli/start-task.md`: Updated Task Discovery section (lines 711-733)
- **{model} placeholder scope**: Removed `{model}` from task prompt placeholder lists and examples (only available in agent commands)
- `docs/cli/start-task.md`: Removed `{model}` from Execution Flow step 8 and Task Placeholders section
- `docs/cli/start-config-task.md`: Removed `{model}` from Placeholders section and Invalid Template error example
- **Role reference clarification**: Renamed "System Prompt File Location" to "Role Reference" in task docs to clarify tasks reference roles (which have files) rather than owning files directly
- `docs/cli/start-task.md`: Renamed and updated section (lines 248-255)
- **Quiet flag consistency**: Updated `--quiet` flag behavior for tasks to suppress output (like other commands) instead of being ignored
- `docs/cli/start-task.md`: Updated Quiet Flag Behavior section and Flags description (lines 128-130, 188-195)
- **Missing Asset Restoration (DR-042)**: Implemented logic to automatically restore missing asset files (prompts, role definitions) from the GitHub catalog if they are referenced in config but missing from disk.
- `docs/design/design-records/dr-042-missing-asset-restoration.md`: Created new DR defining the restoration logic (intercept missing file -> check asset path -> download from catalog).
- `docs/design/design-records/dr-033-asset-resolution-algorithm.md`: Updated to include "Recursive Resolution & Content Restoration" section, linking resolution to file restoration.
- **start-assets-browse.md updates**: Added missing exit code and clarified browser interaction
- `docs/cli/start-assets-browse.md`: Added `(or Ctrl+Click)` to browser failure output
- `docs/cli/start-assets-browse.md`: Added Exit Code 1 for general errors
- `docs/cli/start-assets-search.md`: Updated exit codes to distinguish between errors (1) and no matches (2)
- `docs/cli/start-assets-search.md`: Updated scripting section to reflect new exit codes
- `docs/cli/start-config-agent.md`: Removed incorrect "global config" qualifier from agent definition
- `docs/cli/start-config-agent.md`: Updated bin field description to include index.csv usage
- `docs/cli/start-config-agent.md`: Clarified local agent management (CLI supported)
- `docs/design/design-records/dr-007-placeholders.md`: Clarified placeholder requirements ({prompt} recommended/warning, others recommended)
- **start-config-agent.md updates**:
- Corrected backup filenames in examples (`config.YYYY...` → `agents.YYYY...`) for `new`, `edit`, `remove` commands
- Added scope selection step to `new` command interactive flow
- Verified `default` command correctly uses `config.YYYY...` backups (modifies settings)
- **start-config-context.md updates**:
- Corrected backup filenames in examples (`config.YYYY...` → `contexts.YYYY...`)
- Updated UTD examples to use new placeholders (`{command_output}`, `{file_contents}`)
- **start-config-role.md updates**:
- Corrected backup filenames in examples (`config.YYYY...` → `roles.YYYY...`)
- Updated UTD examples to use `{file_contents}` placeholder instead of `{file}` where content is included
