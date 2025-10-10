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
4. **Binary Strategy**: Introduce a new binary (`quiz-tui`) - standalone binary, not extending the existing `quiz-evaluator`.
5. **Markdown Rendering**: Embed `github.com/charmbracelet/glow`/`glamour` to render markdown inline; no external command invocation.
6. **Evaluation Integration**: When the user opts to evaluate, convert in-memory responses to the evaluator’s expected format temporarily, but do not emit CSV artifacts from the TUI.

## Implementation Details (User-Specified)

### Content Location
- **Subjective Questions Path**: `ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective/`
- **Structure**: Organized by difficulty levels (L3-L7) with baseline and bar-raiser subdirectories
- **Question Format**: YAML frontmatter with fields: id, type, level, category, topic, subtopic, estimated_time
- **Content Sections**: main_question, core_concepts, peripheral_concepts, sample answers, common_mistakes, follow_up questions
- **Available Topics**: GFS (replication, consistency, chunk design) and Raft consensus (intro, basics, election, log replication, performance, evolution)

### User Management
- **User Identity**: Via CLI flag `--user <username>` (e.g., `quiz-tui --user abhishek`)
- **Required**: User must be specified at startup

### Question Selection
- **Interactive Selection**: TUI presents a list of available topics/questions from the markdown folder
- **User Choice**: Users select questions from within the TUI interface (not via CLI flags)
- **Simple Approach**: Show all available questions grouped by level and category for selection

### Session Management
- **Resume on Startup**: If user has existing incomplete sessions, show selection menu
- **Session Selection**: Allow user to pick which session to resume or start new
- **Auto-save**: Every 30 seconds after any change
- **Answer Format**: Plain text only for initial version

### Configuration
- **Config File**: `config/tui.toml` (project-local)
- **Default Settings**:
  - `auto_save_interval = 30` (seconds)
  - `sessions_dir = "sessions"`
  - `content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"`

### Binary Details
- **Name**: `quiz-tui` (standalone, not extending quiz-evaluator)
- **Location**: Will be built to `build/quiz-tui`
- **Build Script**: New `scripts/build_tui.sh` (inspired by existing build.sh)

## Final Implementation Approach (Clarified)

### Question Selection Flow
- **Simple Linear Flow**: Start from L3 and progress sequentially to L7
- **No Individual Selection**: Questions are presented in order by difficulty level
- **Level Progression**: L3-baseline → L3-bar-raiser → L4-baseline → L4-bar-raiser → ... → L7
- **All Questions**: User goes through all available subjective questions in sequence

### Session Behavior
- **No Review Screen**: Users cannot review all answers at once
- **No Explicit Submit**: Session auto-completes when last question is answered
- **Auto-save Triggers**:
  - When user moves to next question (immediate save)
  - Every 30 seconds while user is typing/idle
- **Immediate Persistence**: Each answer is saved as soon as user navigates away

### Navigation Rules
- **Linear Only**: No jumping between questions
- **Sequential Enforcement**: Must proceed through questions in order
- **No Back Navigation**: Once moved to next question, cannot go back (for initial version)
- **Progress Indicator**: Show current position (e.g., "Question 3 of 15")

### Code Reuse Strategy
- **Leverage quiz-evaluator**: Reuse existing code where applicable
- **Shared Modules**: Extract common functionality into shared packages:
  - `internal/models/` - Reuse Question, UserResponse structs
  - `internal/markdown/` - Reuse Scanner, Parser for reading questions
  - `internal/config/` - Extend existing config management
- **New Modules**: 
  - `internal/tui/` - New TUI components and screens
  - `internal/session/` - New session management (JSON-based)

### Display Formatting
- **Minimal Display**: Show only the main_question text
- **Hidden Metadata**: Do not show:
  - Level, category, topic information
  - Estimated time
  - Core concepts or rubrics
  - Sample answers
- **Clean Interface**: Focus on question and answer area only

### Project Structure
```
quiz-evaluator/           # Existing evaluator (keep as-is)
internal/                 # Shared modules (refactored from quiz-evaluator)
  ├── models/            # Question, UserResponse structs
  ├── markdown/          # Scanner, Parser
  ├── config/            # Configuration management
  └── common/            # Other shared utilities
cmd/quiz-tui/            # New TUI binary entry point
internal/tui/            # TUI-specific components
  ├── screens/           # Welcome, Question, Complete screens
  ├── components/        # Reusable UI widgets
  └── session/           # Session management (JSON)
scripts/
  ├── build.sh          # Existing build script
  └── build_tui.sh      # New TUI build script
config/
  └── tui.toml          # TUI configuration file
```

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
5. **Optional Evaluation Step**: Future optional feature - evaluation integration not required for initial implementation.

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

### Data Flow
1. User starts with `quiz-tui --user <name>`
2. System checks for existing sessions
3. User selects resume or new session
4. Questions loaded from markdown files (L3 → L7 order)
5. Linear progression through all questions
6. Auto-save on navigation and timer
7. Session marked complete after last question

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
- Simple CLI with required `--user` flag
- Main command: `quiz-tui --user <username>` 
- On startup:
  - Check for existing incomplete sessions for the user
  - If sessions exist, show selection menu to resume or start new
  - If no sessions, start new quiz from L3
- Uses Cobra for CLI parsing and Viper for config loading

### 7. Configuration Updates
- New `config/tui.toml` with settings:
  - `sessions_dir = "sessions"`
  - `auto_save_interval = 30`
  - `content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"`
- Configuration loaded via Viper with sensible defaults

## Implementation Phases (Simplified)

### Phase 1 – Foundations & Code Reuse (Day 1)
1. Extract shared modules from `quiz-evaluator` to `internal/` directory
2. Create `SessionManager` for JSON persistence with auto-save
3. Set up project structure and dependencies (Bubbletea, Glamour, etc.)
4. Create `scripts/build_tui.sh` based on existing build script

**Deliverable**: Shared modules extracted, session management ready, build script working.

### Phase 2 – TUI Core (Days 2-3)
1. Build Bubbletea skeleton with linear flow: welcome → questions → completion
2. Implement question display (main_question only) with markdown rendering
3. Add textarea component for plain text answers
4. Wire up linear navigation (next only, no back)

**Deliverable**: Basic TUI flow working with question display and answer input.

### Phase 3 – Session Management (Days 4-5)
1. Implement auto-save on navigation and 30-second timer
2. Add resume functionality with session selection on startup
3. Handle session completion when last question answered
4. Add progress indicator showing current position

**Deliverable**: Complete session lifecycle with auto-save and resume working.

### Phase 4 – Polish & Testing (Day 6)
1. Add error handling and edge cases
2. Test with all subjective questions (L3-L7)
3. Ensure config file (`config/tui.toml`) works properly
4. Add keyboard shortcuts and help text

**Deliverable**: Production-ready TUI binary for subjective quiz sessions.

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
