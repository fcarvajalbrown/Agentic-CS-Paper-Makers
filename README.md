# Agentic CS Paper Makers

> A CLI-first, GUI-ready multi-agent workflow for producing structured computer science and game theory papers. Go orchestrator + embedded Python LLM bridge.

## What It Does

1. **Research** — Takes your seed idea, searches arXiv and Semantic Scholar, and grounds it in real literature.
2. **Socratic Architect** — Asks clarifying questions to formalize your idea into a rigorous blueprint.
3. **Lead Writer** — Generates academic Markdown prose from the approved blueprint.
4. **Simulated Peer Review** — Three specialist critics review the draft in parallel.
5. **Inbox Revisions** — You pick which critiques to address, in any order, until you're satisfied.
6. **Finalize** — Produces `paper.md`, `references.bib`, and `metadata.json`, with optional PDF export.

## Quick Start

```bash
# Build
go build -o paperflow ./cmd/paperflow

# Create a new paper project
paperflow init my_stackelberg_paper
cd my_stackelberg_paper

# Run the full workflow
paperflow start "Byzantine Fault Tolerance using Stackelberg competition"

# Check status
paperflow status

# Resume after crash or Ctrl+C
paperflow resume

# Export to PDF (requires pandoc + wkhtmltopdf)
paperflow finalize
paperflow export --pdf
```

## Requirements

- **Go 1.22+** (for building)
- **Python 3.10+** (for CLI mode; embedded in future GUI version)
- **pandoc** + **wkhtmltopdf** (optional, for PDF export)
- A **Kimi / Moonshot API key**

## Architecture

- **Go** — CLI, state machine, checkpointing, schema validation, cross-OS binary
- **Python (embedded)** — LLM API calls, web search, tool execution
- **Agent profiles** — Markdown files defining each agent's role and instructions
- **JSON artifacts** — Immutable, versioned outputs passed between stages

## Configuration

```bash
export PAPERFLOW_API_KEY="your-key"
# or
paperflow --agent-model=architect:kimi-k2.5 --budget=10.00
```

## License

Copyleft — see `LICENSE`.

## Status

Work in progress. See `PRD.md` for full specification.
