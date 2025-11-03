# Design Thoughts

- The specific documents (~/reference/ENVIRONMENT.md, AGENTS.md etc), should be configurable
- Maybe a few "Default Configurations" that can define common defaults to make initial use easy
- Some of the more popular tools like Claude Code, Gemini CLI, OpenCode etc, we should have as default configurations
- For unknown agents, the ability for complete cli argument population in configuration
- Most AI agents have a security permissions config file, the ability to read the context documents should be added to these permissions. Maybe this is outside the scope of start, but a message to do it would fit? Not sure.
- Adding a document through config should include the document path (absolute or relative), and the prompt suffix
- Just thinking about it, maybe a local repo config should override user config (./.start/config.?), or aggregate the global config?
- I want it to be able to use a prompt writer prompt to create role documents on the fly
- I want to support common workflows such as a "git diff review" and others.


Tasks:

- git diff review
- code review
- comment review
- documentation review
- generate git commit message
