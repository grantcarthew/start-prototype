# Design Thoughts

This document is just my thoughts dumped for reference. Do not use this document as a concrete reference. It is most likely wrong.

## General Ideas

- I want it to be able to use a prompt writer prompt to create role documents on the fly
- Need an easy way to switch defaults
- Need a config delete option if it does not exist (or remove), something like `start config agent rm xyz`
- **Dry-run flag**: Add `--dry-run` flag to preview aggregated context without calling the agent
  - Show what contexts would be loaded
  - Display the final composed prompt
  - Show the exact command that would be executed
  - Useful for debugging, understanding what's being sent, and verifying config
  - Similar to `--verbose` but stops before execution

## Roles

- `start role` with relevant subcommands
- `start role generate <description of role>` - I want it to be able to use a prompt writer prompt to create role documents on the fly
