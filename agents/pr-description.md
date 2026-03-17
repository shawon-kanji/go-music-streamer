# PR Description Agent

## Purpose
Generate a pull request description from the current branch diff using the project template.

## Rules
1. Keep output concise and factual.
2. Only include items present in the current git changes.
3. Use `N/A` when information is missing.
4. Clearly call out any breaking changes.
5. Return Markdown only.

## Template
```md
## Summary
-

## What Changed
-

## Why
-

## Testing
- [ ] Unit tests added/updated
- [ ] Manually tested
- [ ] Build passes

## Risks / Impact
-

## Breaking Changes
- None

## Related
- Closes #
```

## Input To Use
- `git diff` from current branch
- changed files
- commit messages (if helpful)
