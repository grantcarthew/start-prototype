# Minimal Example

The simplest possible configuration to get started with `start`.

## What's Here

This example shows only **global configuration** (`~/.config/start/`):

- `config.toml` - Basic settings
- `agents.toml` - Smith test agent
- `roles.toml` - Simple assistant role
- `contexts.toml` - One required context
- `tasks.toml` - Basic help task

## Installation

```bash
# Copy to your global config directory
mkdir -p ~/.config/start
cp examples/minimal/global/* ~/.config/start/

# Verify
start doctor
```

## Local Configuration

Local configs (`./.start/`) follow the **exact same structure** as global:

```bash
# In your project root
mkdir -p .start
cd .start

# Create the same files:
# - config.toml
# - agents.toml
# - roles.toml
# - contexts.toml
# - tasks.toml
```

**Local configs override global** - use them for project-specific settings.

See the `complete/` and `real-world/` examples for local config patterns.

## Next Steps

1. **Test it**: `start --help`
2. **Run a task**: `start task help "explain Go interfaces"`
3. **Customize**: Edit `~/.config/start/*.toml` files
4. **Learn more**: See `examples/README.md`
