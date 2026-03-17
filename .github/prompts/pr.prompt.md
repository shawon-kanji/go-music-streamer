---
agent: pr-generator
---
Generate a pull request description from the current repository changes and return only the final Markdown PR description using the exact template.

If a section has no data, write `N/A`.
Do not include explanations before or after the template output.
Output in markdown format always.
Wrap the entire output in a fenced markdown code block using:
```markdown
...content...
```