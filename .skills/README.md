# Skills

These are model-agnostic system prompts. Paste the contents of any skill file as the system prompt before starting a coding session. They work with Claude, Kimi, GPT, Codex, or any other LLM.

---

## Which skill to use

| You want to... | Use |
|---|---|
| Write code file by file | `coder.md` |
| Review changes before committing | `reviewer.md` |
| Make a design decision or evaluate a tradeoff | `architect.md` |

---

## How to use with Claude Code

1. Start a new session.
2. Type `/agent` or open the agent selector.
3. Paste the contents of the skill file as the system prompt.

Alternatively, copy the file contents and open your session with:
```
<paste skill content here>

Now let's work on: <your task>
```

---

## How to use with Kimi / Moonshot

1. Open a new chat on platform.moonshot.cn or via the API.
2. Set the system prompt to the full contents of the skill file.
3. Start your task in the first user message.

---

## How to use with GPT / Codex (OpenAI)

1. Open a new chat in ChatGPT or the Playground.
2. In the **System** field, paste the full contents of the skill file.
3. Start your task in the first user message.

---

## Tips

- Always start a **fresh session** when switching skills. Mixed context leads to mixed behavior.
- The `coder` skill works best when you also share the relevant files from the project as context.
- The `reviewer` skill works best when you paste the git diff or the changed files directly.
- The `architect` skill expects you to have read `docs/PRD.md` — share it if the model has not seen it.
