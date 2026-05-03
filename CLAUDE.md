# Agentic CS Paper Makers — Claude Code Instructions

## Project Overview

This is a Go-based CLI tool with an embedded Python LLM bridge for orchestrating multi-agent academic paper writing workflows. The Go CLI manages state, checkpoints, schema validation, and cross-OS distribution. The Python bridge handles LLM API calls to Kimi/Moonshot, web search (arXiv, Semantic Scholar), and tool execution.

## Language Stack

| Layer | Language | Responsibility |
|---|---|---|
| CLI & Orchestrator | Go | Commands, state machine, checkpointing, schema validation, cross-OS binary distribution |
| LLM Bridge | Python (embedded) | API calls to Kimi/Moonshot, web search, tool execution |

**Rule of thumb:** If it touches the filesystem, CLI parsing, or state persistence → **Go**. If it talks to an LLM or makes HTTP requests → **Python**.

## Architecture Rules

### Go (`cmd/`, `internal/`, `pkg/`)
- Use `cobra` for CLI commands (we will add this dependency).
- All JSON artifacts must be validated against embedded JSON Schema before the next stage runs.
- Every agent stage writes a versioned artifact to `.paperflow/artifacts/` before proceeding.
- `Ctrl+C` must always trigger graceful abort: save checkpoint, print resume command.
- Cross-compilation target: Windows, macOS, Linux, ARM.

### Python (`python/`)
- The bridge communicates with Go via **stdin/stdout JSON** only.
- Each agent is a standalone module called by `bridge.py`.
- The Kimi/Moonshot client lives in `python/models/kimi_client.py` and handles retries, seed, and streaming.
- Web search tools live in `python/tools/` — arXiv API and Semantic Scholar API only. No scraping.

### Agent Profiles (`agents/`)
- Each agent has a `.md` profile file defining its role and system prompt.
- These are embedded into the Go binary via `//go:embed`.
- **Never edit agent profiles without explicit user sign-off on each.** The user wants to define them one-by-one with their check.

## Workflow Stages

The CLI is strictly sequential where order matters, parallel only for the 3 Critics:

1. `paperflow init` → scaffold project
2. `paperflow start` or `paperflow research` → Research Agent
3. `paperflow architect` → Socratic Architect (human gate)
4. `paperflow write` → Lead Writer
5. `paperflow review` → dispatch 3 Critics in parallel
6. `paperflow inbox` → human picks reviewers, loops until `stop`
7. `paperflow finalize` → assemble paper.md + references.bib + metadata.json

## Configuration Hierarchy

Highest to lowest:
1. CLI flags (`--model`, `--budget`)
2. Environment variables (`PAPERFLOW_API_KEY`, `PAPERFLOW_MODEL`)
3. Project-local `.paperflow/config.json`
4. User-global `~/.config/paperflow/config.json` (Windows: `%APPDATA%/paperflow/config.json`)

## Key Decisions (Locked)

- **Artifacts:** JSON files with strict schemas, stored in `.paperflow/artifacts/`
- **Output:** Markdown (`paper.md`) is canonical. Optional PDF via `pandoc + wkhtmltopdf` (no LaTeX).
- **Models:** Kimi K2.5 is default. Moonshot v1-32k for cheap mode. `--cheap` and `--production` flags swap tiers.
- **No Simulation Agent.** Papers are purely theoretical.
- **Single paper per workflow.** No multi-project workspaces in MVP.
- **Review model:** Inbox style — pick any reviewer in any order, respond, loop until `stop`.

## Testing & Quality

- Go code: standard `go test`.
- Python code: `pytest` in a virtual environment.
- Validate all JSON artifacts against schema before accepting them from the bridge.
- Every LLM call must be cached in `.paperflow/cache/` with hash-based lookup.

## Communication Rules

- When working on files, **always ask before creating or modifying each file** unless the user explicitly gave blanket permission.
- The user prefers **one-thing-at-a-time** interaction. Present one recommendation, wait for feedback, then move to the next.
- The user speaks English and Spanish only.
- Do **not** use `.cn` sites for reference or configuration.

## File Creation Policy

**NEVER create or modify multiple files in one go without explicit user permission.**

The correct flow is:
1. Propose the next file to create or modify
2. Ask for confirmation
3. Create/modify only that file
4. Report completion
5. Ask if ready for the next file

This rule is non-negotiable. The user explicitly requested it.

## When to Use Agents/Subagents

- For read-only codebase exploration across >3 files → use `Agent` with `subagent_type="explore"`
- For implementation planning before writing code → use `Agent` with `subagent_type="plan"`
- For focused coding tasks → use `Agent` with `subagent_type="coder"`
- For simple, single-file edits → do it directly, no agent needed
