# Agentic CS Paper Makers — Claude Code Instructions

## Project Overview

This is a Go-only CLI tool for orchestrating multi-agent academic paper writing workflows. Go handles everything: CLI, state machine, checkpointing, schema validation, LLM API calls to Kimi/Moonshot, web search (arXiv, Semantic Scholar), tool execution, and cross-OS binary distribution.

## Language Stack

| Layer | Language | Responsibility |
|---|---|---|
| Everything | Go | Commands, state machine, checkpointing, schema validation, LLM API calls, web search, tool execution, cross-OS binary distribution |

**Rule of thumb:** It's all Go. No Python, no external runtime dependency. The binary ships alone.

## Architecture Rules

### Go (`cmd/`, `internal/`, `pkg/`)
- Use `cobra` for CLI commands. One command per file under `internal/cli/`.
- All JSON artifacts must be validated against embedded JSON Schema before the next stage runs.
- Every agent stage writes a versioned artifact to `.paperflow/artifacts/` before proceeding.
- `Ctrl+C` must always trigger graceful abort: save checkpoint, print resume command.
- Cross-compilation target: Windows, macOS, Linux, ARM.
- LLM calls go directly from `internal/llm/client.go` to the Kimi/Moonshot REST API.
- Web search tools live in `internal/tools/` — arXiv and Semantic Scholar HTTP clients. No scraping.
- The tool-use loop (LLM returns function call → Go executes → result fed back) lives in `internal/llm/client.go`.

### Agent Profiles (`agents/`)
- Each agent has a `.md` profile file defining its role and system prompt.
- These are embedded into the Go binary via `//go:embed`.
- **Never edit agent profiles without explicit user sign-off on each.** The user wants to define them one-by-one with their check.

## AI Skills (`.skills/`)

Model-agnostic system prompts for working on this codebase. Work with any LLM — Claude, Kimi, GPT, Codex.

| Skill | When to use |
|---|---|
| `.skills/coder.md` | Writing code file by file |
| `.skills/reviewer.md` | Reviewing changes before committing |
| `.skills/architect.md` | Design decisions and tradeoff evaluation |

See `.skills/README.md` for usage instructions per model.

## Coding Conventions

### Conventional Commits
All commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) spec:

```
<type>(scope): <description>

Types: feat, fix, refactor, test, docs, chore, ci
Examples:
  feat(llm): add kimi client with exponential backoff
  fix(inbox): handle empty review round on resume
  chore(deps): add cobra dependency
```

**Never add Co-Authored-By or any Claude/Anthropic attribution to commit messages.** No exceptions.

### Style: Use the Oldest File in the Directory as a Guide
When writing or modifying a file in any directory, read the **oldest existing file** in that directory first and use it as a style reference: package name, import grouping, error handling pattern, naming style, receiver names. Prefer consistency with what's already there, but use judgment — it's a suggestion, not a hard constraint.

### Tests
**Never modify a test to make it pass.** If a test fails, fix the root cause in the production code first. Only update the test itself if the test is genuinely wrong (wrong expectation, stale contract). When in doubt, ask.

## Workflow Stages

The CLI is strictly sequential where order matters, parallel only for the 3 Critics:

1. `paperflow init` → scaffold project
2. `paperflow research` → Research Agent
3. `paperflow architect` → Socratic Architect (human gate)
4. `paperflow write` → Lead Writer
5. `paperflow review` → dispatch 3 Critics in parallel
6. `paperflow inbox` → human picks reviewers, loops until `stop`
7. `paperflow finalize` → assemble paper.md + references.bib + metadata.json

## Build Order

See `docs/codingplan.ini` for the phased build order. Never start a phase until all phases above it are done.

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

- Go code: `go test ./...` must pass before any commit.
- Validate all JSON artifacts against schema before the next stage runs.
- Every LLM call must be cached in `.paperflow/cache/` with hash-based lookup.

## Communication Rules

- When working on files, **always ask before creating or modifying each file** unless the user explicitly gave blanket permission.
- The user prefers **one-thing-at-a-time** interaction. Present one recommendation, wait for feedback, then move to the next.
- The user speaks English and Spanish only.
- Do **not** use `.cn` sites for reference or configuration.
- **NEVER use emojis** — not in code, comments, docs, commit messages, CLI output, or chat responses.

## File Creation Policy

**NEVER create or modify multiple files in one go without explicit user permission.**

The correct flow is:
1. Propose the next file to create or modify
2. Ask for confirmation
3. Create/modify only that file
4. Report completion
5. Ask if ready for the next file

This rule is non-negotiable. The user explicitly requested it.
