Use AIDesigner for the request in $ARGUMENTS.

Workflow:
1. Gather design context before spending credits:
   - Read `DESIGN.md`, `.aidesigner/DESIGN.md`, or `docs/design.md` if present
   - Inspect theme files, tokens, fonts, shared components, and the target route/page
2. Write a compact internal design brief from that repo context:
   - platform and surface
   - product goal and main user action
   - visual language to preserve, evolve, or reset
   - important repo patterns, constraints, and content types
   - concrete typography/tokens only if preserving the current repo aesthetic
3. Split the task into:
   - a visual reference prompt for AIDesigner
   - a detailed implementation spec you keep local
4. Rewrite $ARGUMENTS into a broad visual reference prompt before calling AIDesigner.
   - Focus on product type, audience, UX priorities, overall feel, and non-negotiable constraints.
   - Let AIDesigner choose the specific layout, supporting sections, composition, and most stylistic details.
   - Avoid exact section orders, card counts, micro-copy, and detailed per-element placement unless the user explicitly requested them.
   - Do not pass exhaustive content outlines, tables, command lists, API fields, or other documentation detail dumps into the prompt.
   - If the user gave a very specific page spec, compress it into a smaller visual brief for AIDesigner and save the full spec for local implementation.
   - If matching an existing aesthetic, it is fine to mention concrete repo colors, fonts, and tokens.
   - If this is a new design or revamp, keep styling guidance broad and do not hard-code exact colors or gradients unless the user explicitly asked for them.
5. Prefer the connected `aidesigner` MCP server for `whoami`, `get_credit_status`, `generate_design`, and `refine_design`.
   For v1 MCP usage, keep calls prompt-driven and do not pass `mode` or `url`.
6. After a successful MCP generation or refinement, persist the returned HTML into a local run:
   - Prefer piping the HTML into `npx -y @aidesigner/agent-skills capture --prompt "<final prompt>" --transport mcp --remote-run-id "<run-id>"`
   - If piping is awkward, write the HTML to `.aidesigner/mcp-latest.html` and run `npx -y @aidesigner/agent-skills capture --html-file .aidesigner/mcp-latest.html --prompt "<final prompt>" --transport mcp --remote-run-id "<run-id>"`
7. If MCP is unavailable or needs re-auth:
   - If `AIDESIGNER_API_KEY` is already set, use `npx -y @aidesigner/agent-skills generate --prompt "<final prompt>"`
   - Otherwise stop and show setup instructions:
     - Recommended: `npx -y @aidesigner/agent-skills init`, then open Claude Code, run `/mcp`, and sign into `aidesigner`
     - Cross-repo: `npx -y @aidesigner/agent-skills init --scope user`
     - Fallback: set `AIDESIGNER_API_KEY` and retry
8. After every successful run, use the preview image and run `npx -y @aidesigner/agent-skills adopt --id <run-id>`.
9. Treat the AIDesigner output as the primary visual reference, not a literal page spec.
10. Apply the full user requirements during repo implementation using local content, behaviors, and component structure. Do not ship raw standalone HTML into app code.
