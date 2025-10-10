# Feature: Subjective Quiz TUI with Bubbletea

## Date: 2025-10-10
## Category: Feature
## Status: Planning (updated after stakeholder clarifications)

## Overview
Design and implement a new terminal-native experience for subjective quizzes. The TUI will read subjective questions from our markdown corpus, present them interactively, collect free-form answers, and persist sessions strictly in JSON (no CSV exports). We will ship this as a dedicated binary for subjective quizzes while designing shared infrastructure so an objective/MCQ TUI can be added later without rework.

## Updated Scope Decisions

1. **Quiz Type**: Deliver only subjective-question flows in the first release. MCQ support is out of scope for now, but the architecture must make adding an objective mode straightforward.
2. **Session Storage**: Persist drafts and completed runs solely as JSON using a consistent schema. CSV generation stays outside the TUI.
3. **Session Naming**: Use filenames that embed the quiz taker and timestamp to support multiple resumable sessions: `sessions/<user>/<mode>/<timestamp>-<user>-<mode>.json` (e.g., `sessions/abhishek/subjective/20251011-101530-abhishek-subjective.json`).
4. **Binary Strategy**: Introduce a new binary (`quiz-tui`) with a `--mode` flag. The initial implementation will fully support `subjective`; an `objective` mode stub plus shared plumbing will keep the future MCQ flow consistent.
5. **Markdown Rendering**: Embed `github.com/charmbracelet/glow`/`glamour` to render markdown inline; no external command invocation.
6. **Evaluation Integration**: When the user opts to evaluate, convert in-memory responses to the evaluator’s expected format temporarily, but do not emit CSV artifacts from the TUI.

## Tentative Library Stack
- **github.com/charmbracelet/bubbletea** – core TUI state machine and event loop.
- **github.com/charmbracelet/bubbles** – progress bars, viewports, and other reusable UI widgets.
- **github.com/charmbracelet/lipgloss** – terminal styling, layouts, and color palettes.
- **github.com/charmbracelet/huh** – form and textarea components for subjective answer entry.
- **github.com/charmbracelet/glow** and **github.com/charmbracelet/glamour** – inline markdown rendering with syntax highlighting.
- **github.com/spf13/cobra** – command-line interface for the new binary.
- **github.com/spf13/viper** – configuration loading and environment binding.
- **github.com/google/uuid** (or similar) – optional helper for unique session IDs if timestamps alone are insufficient.

## Objectives

### Primary Goals
1. **Interactive Subjective Questioning**: Render questions (with rich markdown) and collect multiline answers smoothly.
2. **JSON Session Lifecycle**: Auto-save progress, resume incomplete runs, and mark completion in a single JSON artifact.
3. **Progress Awareness**: Visual cues for position, completion percentage, and unanswered items.
4. **Review & Submission**: Let users review, edit, and confirm answers before finalizing.
5. **Optional Evaluation Step**: Chain to the existing evaluator using the freshly collected responses without writing CSV to disk.

### User Experience Goals
- Keyboard-first navigation and shortcuts.
- Styled terminal layout using Lipgloss/Bubbles.
- Automatic draft saves (interval + state-change triggers).
- Resume flow that detects available sessions per user.
- Future-ready layout hooks to plug in MCQ widgets when objective mode lands.

## Reusable Foundations

✅ **internal/markdown/**: `Scanner` + `Parser` already parse subjective content (Core Concepts, Rubrics, etc.).

✅ **internal/models/**: `Question` and `UserResponse` exist; we will extend them to support mode metadata and future MCQ fields without affecting current consumers.

✅ **internal/config/**: Viper-backed config wiring we can reuse for TUI defaults and directories.

✅ **internal/evaluator/**: The CLI evaluator can be orchestrated post-quiz without modification.

⚠️ **Enhancements Needed**:
1. `Question` struct: add `QuestionType` (default `subjective`) plus optional MCQ fields (`Options`, `CorrectAnswer`) for future use.
2. `Parser`: ensure we retain MCQ parsing hooks but gate UI presentation to subjective questions only.
3. New JSON session module for saving/resuming quiz runs.
4. Shared quiz engine abstractions to support multiple modes.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                       quiz-tui binary                       │
│                                                             │
│  ┌──────────────────┐       ┌──────────────────────────┐    │
│  │  Cobra CLI layer │──────►│  Mode Registry           │    │
│  └──────────────────┘       │  - SubjectiveMode (now)  │    │
│                             │  - ObjectiveMode (future)│    │
│                             └──────────────────────────┘    │
│                                                             │
│           ┌───────────────────────────────┐                 │
│           │  Quiz Engine (shared)         │                 │
│           │  - Loader (markdown)          │                 │
│           │  - Session Manager (JSON)     │                 │
│           │  - Response Collector         │                 │
│           └───────────────────────────────┘                 │
│                         │                                    │
│           ┌───────────────────────────────┐                 │
│           │  TUI Layer (Bubbletea)        │                 │
│           │  - Subjective screens         │                 │
│           │  - Shared components          │                 │
│           │  - Glow renderer              │                 │
│           └───────────────────────────────┘                 │
│                         │                                    │
│           ┌───────────────────────────────┐                 │
│           │  Evaluator Bridge (optional) │                 │
│           │  - JSON adapter               │                 │
│           │  - CLI invocation             │                 │
│           └───────────────────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
```

### Mode Strategy
- Introduce a `Mode` interface with lifecycle hooks (`Prepare`, `BuildModel`, `SupportsQuestion(*Question)` etc.).
- Implement `SubjectiveMode`, which filters questions where `QuestionType == "subjective"` and provisions subjective-specific screens.
- Ship an `ObjectiveMode` skeleton returning “not yet implemented” but exercising the same interfaces to validate extensibility.

## JSON Session Specification

```json
{
  "session": {
    "session_id": "20251011-101530-abhishek-subjective",
    "user": "abhishek",
    "mode": "subjective",
    "status": "in_progress",   // in_progress | completed | aborted
    "created_at": "2025-10-11T10:15:30Z",
    "updated_at": "2025-10-11T10:35:12Z",
    "question_count": 10,
    "answered": 6
  },
  "questions": [
    {
      "id": "ch07-write-skew",
      "title": "Write Skew in Snapshot Isolation",
      "level": "L5",
      "chapter": "7",
      "metadata": {
        "tags": ["transactions", "anomalies"]
      }
    }
  ],
  "responses": [
    {
      "question_id": "ch07-write-skew",
      "answer": "Write skew occurs when...",
      "updated_at": "2025-10-11T10:24:02Z",
      "time_spent_seconds": 142
    }
  ]
}
```

- Draft files live under `sessions/<user>/<mode>/`.
- Resuming scans the directory for the latest `status != completed` file for the active user/mode.
- Completed sessions remain in the same directory for audit/history.

## Detailed Component Plan

### 1. Models & Parsing
**Goals**: Enrich `Question` metadata while keeping current downstream compatibility.
- Add `QuestionType` enum string with helper `IsSubjective()` / `IsObjective()`.
- Keep MCQ fields optional so the parser can populate them even though the TUI ignores them for now.
- Parser should recognise MCQ sections but defer to the mode to decide which questions are loadable; this prevents rework when objective support arrives.

### 2. Quiz Engine Layer (new `internal/quizengine` package)
Responsibilities:
- Accept `Mode` definitions and orchestrate quiz runs.
- Provide `Loader` with chapter/level/tag filters and enforce `Mode.SupportsQuestion`.
- Manage `SessionManager` for JSON read/write, including auto-save intervals.
- Offer `ResponseStore` abstraction for the TUI to interact with answers.

### 3. TUI Layer (`internal/tui`)
- Bubbletea `Model` holds generic state plus a `ModeView` interface supplied by the active mode.
- Subjective mode contributes screens: welcome/setup, question viewer (markdown via Glow), response editor (Huh textarea), review, submit, complete.
- Shared components (progress bar, toast notifications, help overlay) live in `internal/tui/components` for re-use.

### 4. Markdown Rendering
- Embed Glow (`github.com/charmbracelet/glow`) with a custom glamour style tuned to our palette.
- Cache rendered markdown per question to avoid repeated transformations while navigating.

### 5. Evaluation Bridge (JSON-only)
- After submission, optionally convert in-memory responses into the evaluator’s request structs without persisting CSV files.
- Stream JSON-derived payloads directly to the evaluator CLI and capture the output for display inside the completion screen.

### 6. CLI Binary (`cmd/quiz-tui/main.go`)
- Cobra root command `quiz-tui` with `--mode` (default `subjective`).
- Commands:
  - `start` (primary entry) — loads mode, sets up session, launches Bubbletea program.
  - `resume` — accepts optional `--session` path; otherwise auto-picks the latest draft for the user/mode.
  - `list-sessions` — helper to inspect stored drafts (preparing for multi-mode usage).
- For now, invoking with `--mode objective` prints a TODO message referencing future work, proving the seam exists.

### 7. Configuration Updates
- Extend config to include `tui.sessions_dir`, `tui.auto_save_interval`, `tui.default_mode`.
- Provide `tui.modes.subjective` settings (default chapter set, default question count, evaluation defaults).
- Reserve `tui.modes.objective` section for future parameters so existing configs need not change later.

## Implementation Phases

### Phase 1 – Foundations (Days 1-2)
1. Introduce `Mode` abstractions and `SubjectiveMode` implementation scaffold.
2. Extend `Question`/`Parser` with `QuestionType` while ensuring existing tooling keeps working.
3. Create `SessionManager` writing JSON using the new naming convention.
4. Add Glow/Glamour dependencies to `go.mod`.

**Deliverable**: CLI can load subjective questions, create a session file, and exit gracefully (no UI yet).

### Phase 2 – TUI Shell (Days 3-5)
1. Build Bubbletea skeleton with welcome → question → completion flow using placeholder content.
2. Implement Glow-based markdown rendering and subjective answer textarea with auto-save hooks.
3. Wire progress bar and navigation shortcuts.

**Deliverable**: Interactive terminal flow for subjective questions with JSON drafts persisted.

### Phase 3 – Review & Submission (Days 6-7)
1. Add review screen summarising answered/unanswered questions.
2. Implement submission confirmation and final session state transition (`completed`).
3. Provide resume logic picking up saved sessions.

**Deliverable**: End-to-end subjective quiz run with resume support.

### Phase 4 – Evaluation Bridge & Polish (Days 8-9)
1. Integrate optional evaluator invocation using in-memory conversion.
2. Add error handling, toast messages, and help overlay.
3. Finalise configuration wiring and CLI UX.

**Deliverable**: Production-ready subjective TUI binary with optional evaluation.

### Phase 5 – Future Mode Hooks (Day 10)
1. Add `ObjectiveMode` stub implementing the Mode interface but returning `mode not implemented` messages.
2. Document extension points within code comments (minimal, only where necessary) to guide future MCQ work.

**Deliverable**: Codebase ready to host an objective TUI with minimal refactor.

## Testing Strategy

### Unit Tests
- `SessionManager` JSON read/write + naming convention.
- `ModeRegistry` selection logic.
- `Loader` filtering for subjective-only mode.
- Resume logic selecting the latest draft per user/mode.

### Integration Tests
- CLI `start` → answer a subset → exit → resume → submit flow.
- Auto-save interval triggers while typing.
- Evaluation bridge invoked after completion (behind a feature flag in tests).

### Manual QA
- Vary terminal sizes/themes.
- Run concurrent sessions for different users.
- Force interruptions (Ctrl+C) to confirm drafts persist.

## Risks & Mitigations
| Risk | Mitigation |
| --- | --- |
| JSON corruption on crash | Write to temp file then atomically rename. |
| Resume collisions when multiple drafts exist | Present list of drafts in resume flow if more than one match. |
| Future objective mode diverges from shared UI | Keep shared components mode-agnostic; isolate subjective specifics in mode package. |

## Future Objective Mode Considerations
- Implement `ObjectiveMode` by reusing Loader but checking `QuestionType == "objective"`.
- Provide MCQ answer widgets (radio/select) inside mode-specific view builders.
- Add converter to aggregate both subjective and objective answers in the same JSON schema (responses array already supports arbitrary answer payloads).
- Same session naming format ensures both modes co-exist cleanly.

## Conclusion
The revised plan focuses on a JSON-only, subjective-first TUI delivered as a dedicated binary while laying groundwork for a future objective experience. By introducing mode abstractions, a shared engine, and a robust session manager, we guarantee that adding MCQ support later requires only mode-specific UI and answer handling, not wholesale rewrites. Glow-based rendering, auto-save, and resume flows provide a polished terminal experience aligned with the updated stakeholder requirements.
