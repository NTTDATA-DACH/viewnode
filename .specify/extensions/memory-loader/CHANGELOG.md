# Changelog

## [1.0.0] - 2026-04-20

### Added
- Initial release
- `speckit.memory-loader.load` command to read all `.specify/memory/*.md` files
- `before_*` hooks for specify, plan, tasks, implement, clarify, checklist, and analyze lifecycle commands
- Graceful degradation when memory directory is missing or files are unreadable
