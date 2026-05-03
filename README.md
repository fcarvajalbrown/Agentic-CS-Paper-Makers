# Agentic CS Paper Makers

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows%20%7C%20arm-lightgrey)
![License](https://img.shields.io/badge/License-GPLv3-blue.svg)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

> A CLI-first multi-agent workflow for producing structured computer science and game theory papers. Pure Go — no external runtime required.

## What It Does

1. **Research** — Takes your seed idea, searches arXiv and Semantic Scholar, and grounds it in real literature.
2. **Socratic Architect** — Asks clarifying questions to formalize your idea into a rigorous blueprint.
3. **Lead Writer** — Generates academic Markdown prose from the approved blueprint.
4. **Simulated Peer Review** — Three specialist critics review the draft in parallel.
5. **Inbox Revisions** — You pick which critiques to address, in any order, until you are satisfied.
6. **Finalize** — Produces `paper.md`, `references.bib`, and `metadata.json`, with optional PDF export.

## Quick Start

```bash
# Build
go build -o paperflow ./cmd/paperflow

# Create a new paper project
paperflow init my_stackelberg_paper
cd my_stackelberg_paper

# Run the full workflow step by step
paperflow research "Byzantine Fault Tolerance using Stackelberg competition"
paperflow architect
paperflow write
paperflow review
paperflow inbox
paperflow finalize

# Check status or resume after a crash
paperflow status
paperflow resume

# Export to PDF (requires pandoc + wkhtmltopdf)
paperflow export --pdf
```

## Requirements

- **Go 1.26+**
- **A Kimi / Moonshot API key** (`PAPERFLOW_API_KEY`)
- **pandoc** + **wkhtmltopdf** — optional, for PDF export only

## Architecture

- **Go** — CLI, state machine, checkpointing, schema validation, LLM calls, web search, cross-OS binary
- **Agent profiles** — Markdown files defining each agent's role and system prompt, embedded in the binary
- **JSON artifacts** — Immutable, versioned outputs passed between stages, validated against JSON Schema

## Configuration

```bash
export PAPERFLOW_API_KEY="your-key"

# Override model or set a spend cap
paperflow research --agent-model=research:kimi-k2.5 --budget=10.00

# Cheap mode (moonshot-v1-32k) or production mode (kimi-k2.6)
paperflow write --cheap
paperflow write --production
```

## AI Skills

The `.skills/` directory contains model-agnostic system prompts for working on this codebase. They work with any LLM — Claude, Kimi, GPT, or Codex. See [`.skills/README.md`](.skills/README.md) for usage instructions.

| Skill | Purpose |
|---|---|
| `coder.md` | Writing code file by file |
| `reviewer.md` | Reviewing changes before committing |
| `architect.md` | Design decisions and tradeoff evaluation |

## License

GPL v3 — see `LICENSE`.

## Status

In development. See [`docs/PRD.md`](docs/PRD.md) for the full specification and [`codingplan.ini`](codingplan.ini) for the build order.
