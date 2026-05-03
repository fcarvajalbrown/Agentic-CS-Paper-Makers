# AGENTS.md — Agentic CS Paper Makers

## Project Purpose

Agentic CS Paper Makers is a CLI-first academic paper generation workflow targeting computer science and game theory research. It orchestrates multiple specialized LLM agents through a Go-based CLI with an embedded Python LLM bridge. The output targets preprint servers (Zenodo) and professional networks (LinkedIn), with academic rigor as a quality aspiration rather than a journal-submission guarantee.

## Technology Stack

### Core
- **Go 1.22+** — CLI, state machine, checkpointing, schema validation, cross-OS binary distribution
- **Python 3.10+** — LLM API bridge, web search, tool execution (embedded via `//go:embed`)
- **JSON** — All inter-agent artifacts use versioned JSON with strict schema validation

### Key Dependencies (Planned)
- Go: `cobra` (CLI), standard library for HTTP, JSON, filesystem
- Python: `requests` or `httpx` (HTTP), minimal external deps

### Output
- Canonical: **Markdown** (`paper.md`) with YAML frontmatter, LaTeX math blocks
- Optional PDF: delegated to external `pandoc` + `wkhtmltopdf` (not bundled)
- Bibliography: auto-generated `references.bib` from placeholder citations

## Architecture Principles

1. **Language separation:** Go handles everything except LLM interaction. Python handles LLM calls, search, and tool use.
2. **Go-Python IPC:** Strict JSON over stdin/stdout. No shared memory, no sockets, no file-based IPC.
3. **Immutable artifacts:** Every stage produces a versioned JSON file in `.paperflow/artifacts/`. Never mutate in place.
4. **Checkpoint survival:** `Ctrl+C` at any point must gracefully save state. `paperflow resume` must always work.
5. **Schema-first:** No artifact moves to the next stage without passing JSON Schema validation.

## Directory Conventions

```
.
├── cmd/paperflow/          # Go entry point
├── internal/               # Go private packages
│   ├── cli/                # Commands, flags, help
│   ├── config/             # Config hierarchy
│   ├── state/              # Checkpoint, resume, rollback
│   ├── agents/             # Orchestrator, runner, registry
│   ├── artifacts/          # Schema definitions, validator, store
│   ├── llm/                # Python bridge interface, cache, cost tracking
│   └── export/             # Finalize and export logic
├── pkg/embed/              # go:embed directives
├── python/                 # Python bridge (embedded into binary)
│   ├── bridge.py           # Main router
│   ├── agents/             # Agent implementations
│   ├── tools/              # Web search tools
│   ├── models/             # LLM API client
│   └── utils/              # IPC handler
├── agents/                 # Agent profile .md files (embedded)
├── schemas/                # JSON Schema files (embedded)
├── go.mod, go.sum
├── Makefile
├── PRD.md                  # Full product spec
├── README.md               # User-facing docs
└── CLAUDE.md               # Claude Code specific instructions
```

## Coding Conventions

### Go
- Use `cobra` for subcommands. One command per file under `internal/cli/`.
- Return errors, do not panic. Every function that can fail returns `(T, error)`.
- All file paths use `filepath.Join`, never string concatenation.
- Configuration loading: `flags > env > local JSON > global JSON > defaults`.
- State writes are atomic: write to temp file, then rename.

### Python
- Modules are standalone. Each agent can be imported and tested independently.
- The bridge receives JSON on stdin and prints JSON on stdout. Nothing else to stdout.
- Use type hints on all public functions.
- Minimal external dependencies. Prefer stdlib over third-party.

### JSON Schemas
- Schemas live in `schemas/` and are embedded in the Go binary.
- Every artifact has a `version` field for forward compatibility.
- Schemas validate structure, required fields, and basic types. Not business logic.

## File Creation Policy

**This is enforced for all agents working on this project:**

- Create or modify **one file at a time**.
- Ask for explicit user confirmation before each file.
- Wait for the user to approve before proceeding to the next file.
- Exception: trivial single-line fixes in a known file may be done directly, but announce them.

## Workflow When Working on This Project

1. **Read the PRD** before implementing any feature. The PRD is the source of truth.
2. **Read existing code** before modifying. Use exploration agents if you need to scan >3 files.
3. **Propose, then execute.** Say what file you're about to create/modify, wait for confirmation.
4. **Test after writing.** Run `go test` or `pytest` as appropriate. Report results.
5. **Update docs** if you change behavior described in PRD, README, or this file.

## Testing Standards

- Go: `go test ./...` must pass before any commit.
- Python: `pytest` in a dedicated venv.
- Every schema change needs a corresponding test fixture (valid and invalid examples).
- Cache and checkpoint tests must verify crash recovery.

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
