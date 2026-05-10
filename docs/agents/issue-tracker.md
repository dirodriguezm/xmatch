# Issue Tracker — Local Markdown

Issues are tracked as markdown files under `.scratch/<feature>/` in this repo.

## Creating an issue

1. Create a directory under `.scratch/` named after the feature or ticket.
2. Write a markdown file (e.g., `issue.md`) with the description, acceptance criteria, and any relevant context.

## Listing issues

```bash
find .scratch/ -name '*.md' -type f
```

## Conventions

- Each issue lives in its own `.scratch/<name>/` directory.
- Use `issue.md` as the primary file. Additional files (design notes, sketches) can live alongside it.
- Delete the directory when the work is merged or abandoned.
