# govips

Go bindings for libvips, the fast image processing library.

**Read AGENTS.md first.** It has the setup, commands, layout, cgo/memory rules, testing conventions, and known local quirks. This file only adds Claude-specific workflow.

## Dev Flow

Flow: worktree
- All code changes happen in worktrees, never on main
- Use /dev to start work (creates worktree automatically)
- Use /stage to wrap up (prepares clean commit for landing)
- Review and land via wtr (ff-only merge)

## Issues

Tracked on GitHub: https://github.com/davidbyttow/govips/issues
