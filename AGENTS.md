# AGENTS.md — Agentic CS Paper Makers

## Project Purpose

Agentic CS Paper Makers is a CLI-first academic paper generation workflow targeting computer science and game theory research. It orchestrates multiple specialized LLM agents through a Go-only CLI. The output targets preprint servers (Zenodo) and professional networks (LinkedIn), with academic rigor as a quality aspiration rather than a journal-submission guarantee.

## Technology Stack

### Core
- **Go 1.26** — CLI, state machine, checkpointing, schema validation, LLM API calls, web search, tool execution, cross-OS binary distribution
- **JSON** — All inter-agent artifacts use versioned JSON with strict schema validation

### Dependencies
- `github.com/spf13/cobra v1.9.0` — CLI commands

### Output
- Canonical: **Markdown** (`paper.md`) with YAML frontmatter, LaTeX math blocks
- Optional PDF: delegated to external `pandoc` + `wkhtmltopdf` (not bundled)
- Bibliography: auto-generated `references.bib` from placeholder citations

## Architecture Principles

1. **Go-only:** No Python, no external runtime. The binary ships alone.
2. **LLM calls are plain HTTP:** `internal/llm/kimi_client.go` calls the Kimi/Moonshot REST API directly.
3. **Tool-use loop in Go:** LLM returns a function call → Go executes the HTTP tool → result fed back into conversation. Lives in `internal/llm/kimi_client.go`.
4. **Immutable artifacts:** Every stage produces a versioned JSON file in `.paperflow/artifacts/`. Never mutate in place.
5. **Checkpoint survival:** `Ctrl+C` at any point must gracefully save state. `paperflow resume` must always work.
6. **Schema-first:** No artifact moves to the next stage without passing JSON Schema validation.

## Directory Conventions

```
.
├── cmd/paperflow/          # Go entry point
├── internal/               # Go private packages
│   ├── cli/                # One file per command + flags.go + help.go
│   ├── config/             # Config hierarchy
│   ├── state/              # Checkpoint, resume, rollback
│   ├── agents/             # Orchestrator, runner, registry
│   ├── artifacts/          # Schema definitions, validator, store
│   ├── llm/                # Kimi/Moonshot client, tool-use loop, cache, cost tracking
│   ├── tools/              # arXiv and Semantic Scholar HTTP clients
│   ├── inbox/              # Interactive reviewer response loop logic
│   └── export/             # Finalize and export logic
├── pkg/embed/              # go:embed directives for agent profiles + schemas
├── agents/                 # Agent profile .md files (embedded into binary)
├── schemas/                # JSON Schema files (embedded into binary)
├── .skills/                # Model-agnostic system prompts (coder, reviewer, architect)
├── .vscode/                # VSCode workspace settings and extension recommendations
├── docs/PRD.md             # Full product spec
├── codingplan.ini          # Phased build order
├── go.mod, go.sum
├── Makefile
├── README.md               # User-facing docs
└── CLAUDE.md               # Claude Code specific instructions
```

## AI Skills (`.skills/`)

Model-agnostic system prompts for working on this codebase. Work with Claude, Kimi, GPT, Codex, or any LLM.

| Skill | When to use |
|---|---|
| `coder.md` | Writing code file by file |
| `reviewer.md` | Reviewing changes before committing |
| `architect.md` | Design decisions and tradeoff evaluation |

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

### Go
- Use `cobra` for subcommands. One command per file under `internal/cli/`.
- Return errors, do not panic. Every function that can fail returns `(T, error)`.
- All file paths use `filepath.Join`, never string concatenation.
- Configuration loading: `flags > env > local JSON > global JSON > defaults`.
- State writes are atomic: write to temp file, then rename.

### Style: Use the Oldest File in the Directory as a Guide
When writing or modifying a file in any directory, read the **oldest existing file** in that directory first and use it as a style reference: package name, import grouping, error handling pattern, naming style, receiver names. Prefer consistency with what's already there, but use judgment — it's a suggestion, not a hard constraint.

### JSON Schemas
- Schemas live in `schemas/` and are embedded in the Go binary.
- Every artifact has a `version` field for forward compatibility.
- Schemas validate structure, required fields, and basic types. Not business logic.

## Testing Standards

- `go test ./...` must pass before any commit.
- Every schema change needs a corresponding test fixture (valid and invalid examples).
- Cache and checkpoint tests must verify crash recovery.
- **Never modify a test to make it pass.** If a test fails, fix the root cause in the production code first. Only update the test itself if the test is genuinely wrong (wrong expectation, stale contract). When in doubt, ask.

## Communication Rules

- **NEVER use emojis** — not in code, comments, docs, commit messages, CLI output, or agent responses.

## File Creation Policy

**This is enforced for all agents working on this project:**

- Create or modify **one file at a time**.
- Ask for explicit user confirmation before each file.
- Wait for the user to approve before proceeding to the next file.
- Exception: trivial single-line fixes in a known file may be done directly, but announce them.

## Workflow When Working on This Project

1. **Read the PRD** before implementing any feature. The PRD is the source of truth.
2. **Check `codingplan.ini`** for the current phase and build order.
3. **Read existing code** before modifying. Use exploration agents if you need to scan >3 files.
4. **Propose, then execute.** Say what file you are about to create/modify, wait for confirmation.
5. **Test after writing.** Run `go test ./...` and report results.
6. **Update docs** if you change behavior described in PRD, README, or this file.

## Multi-Agent Workflow (Runtime, Not Development)

This project implements a multi-agent workflow, but **development is single-agent**:
- The Go orchestrator dispatches agents in sequence or parallel.
- The 3 Critics (Formalist, Game Theorist, Skeptic) run in parallel via goroutines.
- The Review Inbox is asynchronous: human picks reviewers in any order.

## Key Constraints

- **Single paper per workflow.** No multi-project workspaces in MVP.
- **No LaTeX dependency.** Markdown is the canonical output.
- **No Simulation Agent.** Papers are purely theoretical. No generated code execution.
- **No web scraping.** arXiv API and Semantic Scholar API only.
- **No GUI in MVP.** CLI only. GUI is a future layer.
