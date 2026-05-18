# Domain Docs

## Layout

Single-context. The following files live at the repo root:

- `CONTEXT.md` — project domain language, key concepts, glossary, architecture overview.
- `docs/adr/` — architectural decision records (ADRs), one per decision.

## Consumer rules

- Skills (`improve-codebase-architecture`, `diagnose`, `tdd`) read `CONTEXT.md` first to learn domain terminology.
- ADRs in `docs/adr/` are consulted when decisions need to align with past choices.
- If `CONTEXT.md` does not exist, the skill should ask the user before proceeding with domain-heavy work.
