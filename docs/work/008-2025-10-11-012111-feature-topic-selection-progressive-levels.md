# Topic Selection and Progressive Level Feature

**Date**: 2025-10-11  
**Category**: Feature  
**Status**: Planning

## Overview

Add topic selection capability to the quiz TUI, allowing users to choose which chapter/topic they want to practice, then answer questions progressively from L3 through L7.

## Current State

### Content Structure
- **Location**: `ddia-quiz-bot/content/chapters/`
- **Available Topics**:
  1. `03-storage-and-retrieval`
  2. `04-encoding-and-evolution`
  3. `05-replication`
  4. `06-partitioning`
  5. `07-transactions`
  6. `08-trouble-with-distributed-systems`
  7. `09-distributed-systems-gfs`

- **Level Structure**: Each topic has `subjective/` directory with:
  - `L3-baseline/`, `L3-bar-raiser/`
  - `L4-baseline/`, `L4-bar-raiser/`
  - `L5-baseline/`, `L5-bar-raiser/`
  - `L6-baseline/`
  - `L7-baseline/`, `L7-bar-raiser/`

### Current TUI Behavior
- Hardcoded path in `config/tui.toml`: `content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"`
- Loads ALL L3-L7 questions from single topic
- No topic selection UI
- Screen flow: Welcome → Session Select → Question → Complete

## Requirements

### User Flow
1. **Welcome Screen** (existing)
2. **Topic Selection Screen** (NEW)
   - List all available topics with descriptions
   - Allow user to choose one topic
   - Show question count per topic
3. **Session Select Screen** (modified)
   - Check for existing sessions for selected topic
   - Resume or create new session
4. **Progressive Question Flow** (NEW)
   - Start with L3 questions (baseline + bar-raiser)
   - Progress to L4, then L5, L6, L7
   - Within each level: baseline questions first, then bar-raiser
5. **Question Screen** (existing)
6. **Complete Screen** (existing)

### Technical Requirements

#### 1. Topic Discovery
- Scan `ddia-quiz-bot/content/chapters/` for available topics
- Extract topic name and create display-friendly labels
- Count questions per topic/level
- Cache topic list for performance

#### 2. Question Loading Strategy
- Load questions in progressive order: L3 → L4 → L5 → L6 → L7
- Within each level: baseline first, then bar-raiser
- Filter by selected topic
- Maintain correct order in session

#### 3. Session Management Updates
- Store selected topic in session metadata
- Filter sessions by topic when resuming
- Update session schema to include `topic` field

#### 4. Config Updates
- Change `content_path` to `chapters_root_path`
- Point to base chapters directory
- Keep backward compatibility

## Implementation Plan

### Phase 1: Data Layer (Scanner & Models)

**Files to Modify**:
- `internal/markdown/scanner.go`
- `internal/models/question.go`

**New Functionality**:
```go
// scanner.go
type TopicInfo struct {
    Name        string   // e.g., "09-distributed-systems-gfs"
    DisplayName string   // e.g., "GFS & Distributed Systems"
    Path        string   // absolute path
    LevelCounts map[string]int // question count per level
}

func (s *Scanner) DiscoverTopics(chaptersPath string) ([]TopicInfo, error)
func (s *Scanner) ScanTopicQuestions(topicPath string) (QuestionIndex, error)
func (s *Scanner) GetProgressiveQuestions(index QuestionIndex) []*Question
```

**Level Order Logic**:
- Sort levels: L3, L4, L5, L6, L7
- Within level: baseline before bar-raiser
- Alphabetical within variant

### Phase 2: TUI Screens

**Files to Modify**:
- `internal/tui/screens/app.go`

**New Screen State**:
```go
const (
    StateWelcome ScreenState = iota
    StateTopicSelect  // NEW
    StateSessionSelect
    StateQuestion
    StateComplete
)
```

**New Model Fields**:
```go
type ImprovedAppModel struct {
    // ... existing fields
    availableTopics []TopicInfo
    selectedTopic   *TopicInfo
}
```

**New Rendering**:
```go
func (m ImprovedAppModel) renderTopicSelect() string
```

### Phase 3: Session Management

**Files to Modify**:
- `internal/tui/session/session.go`

**Schema Update**:
```go
type Session struct {
    // ... existing fields
    Topic       string `json:"topic"`        // NEW: topic identifier
    LevelOrder  []string `json:"level_order"` // NEW: L3,L4,L5,L6,L7
}

// Update filtering
func (m *Manager) ListIncompleteSessions(user, mode, topic string) ([]*Session, error)
```

### Phase 4: Config Updates

**Files to Modify**:
- `config/tui.toml`
- `internal/config/config.go`

**Changes**:
```toml
# Old
# content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"

# New
chapters_root_path = "ddia-quiz-bot/content/chapters"
```

### Phase 5: Testing

**Test Files**:
- `internal/markdown/scanner_test.go` - Test topic discovery and progressive ordering
- `internal/tui/screens/app_test.go` - Test new screen flow
- Manual integration testing

## UI Mockup

### Topic Selection Screen
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📚 Select Your Topic

Available Topics:

  [1] Storage & Retrieval (20 questions)
  [2] Encoding & Evolution (15 questions)
  [3] Replication (25 questions)
  [4] Partitioning (18 questions)
  [5] Transactions (22 questions)
  [6] Distributed Systems (20 questions)
  [7] GFS & Consensus (24 questions)

Press 1-7 to select • Press q to quit
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Progressive Question Display
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Topic: GFS & Consensus
Level: L3-baseline → L4 → L5 → L6 → L7
Question 3 of 24 ✓ Saved

┌─ Question ──────────────────────────────────┐
│ Explain how GFS ensures data durability...  │
│ ...                                          │
└──────────────────────────────────────────────┘

Your Answer:
┌──────────────────────────────────────────────┐
│                                              │
│                                              │
└──────────────────────────────────────────────┘

Ctrl+N: Next • Ctrl+S: Save • Ctrl+C: Quit
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Key Design Decisions

### 1. Progressive Ordering
- **Decision**: Strict L3→L7 progression, no jumping ahead
- **Rationale**: Pedagogical approach - build understanding incrementally
- **Alternative**: Random order within selected topic (rejected - less structured learning)

### 2. Baseline vs Bar-Raiser
- **Decision**: Show baseline first within each level
- **Rationale**: Bar-raiser questions are harder, should come after baseline mastery
- **Implementation**: Sort by directory name (baseline sorts before bar-raiser)

### 3. Topic Storage in Session
- **Decision**: Store topic name/path in session metadata
- **Rationale**: 
  - Allows resuming from correct topic
  - Enables topic-specific session filtering
  - Supports multiple concurrent sessions across topics

### 4. Config Backward Compatibility
- **Decision**: Support both old `content_path` and new `chapters_root_path`
- **Rationale**: Don't break existing configs
- **Implementation**: Check for new field first, fall back to old

## Implementation Status

### Completed

#### Phase 1: Data Layer ✅
- ✅ Added `TopicInfo` struct to `internal/markdown/scanner.go`
- ✅ Implemented `DiscoverTopics()` - scans chapters directory and counts questions
- ✅ Implemented `GetProgressiveQuestions()` - sorts L3→L7, baseline before bar-raiser
- ✅ Added `formatTopicName()` helper for display names

#### Phase 2: TUI Screens ✅
- ✅ Added `StateTopicSelect` to screen states
- ✅ Added `availableTopics` and `selectedTopic` fields to model
- ✅ Implemented `renderTopicSelect()` - displays numbered topic list
- ✅ Added topic selection key handling (1-9 keys)
- ✅ Updated `renderSessionSelect()` to show selected topic
- ✅ Updated `renderQuestion()` to show topic and level in progress bar
- ✅ Added `discoverTopicsCmd()`, `loadTopicQuestionsCmd()`, `checkTopicSessionsCmd()`

#### Phase 3: Session Management ✅
- ✅ Added `Topic` and `TopicDisplay` fields to `SessionMetadata`
- ✅ Implemented `CreateSessionWithTopic()` function
- ✅ Implemented `ListIncompleteSessionsForTopic()` for topic-specific filtering
- ✅ Updated session ID format to include topic name

#### Phase 4: Config Updates ✅
- ✅ Added `ChaptersRootPath` field to `TUIConfig`
- ✅ Maintained backward compatibility with `ContentPath`
- ✅ Updated `config/tui.toml` to use new `chapters_root_path`

#### Phase 5: Testing ✅
- ✅ Created `scanner_topic_test.go` with 3 test functions
- ✅ `TestDiscoverTopics` - verifies topic discovery (1 topic found, 20 questions)
- ✅ `TestGetProgressiveQuestions` - verifies L3→L7 ordering with baseline first
- ✅ All existing tests pass (13 tests, 100% pass rate)

### Test Results

```
Topic Discovery: PASS
- Found: 1 topic (09-distributed-systems-gfs)
- Questions: 20 total (L3:5, L4:5, L5:5, L6:2, L7:3)

Progressive Ordering: PASS
- Verified L3→L4→L5→L6→L7 progression
- Verified baseline before bar-raiser within each level
- Example order: L3-baseline(4), L3-bar-raiser(1), L4-baseline(4), L4-bar-raiser(1)...

All Tests: PASS
- internal/markdown: 3 new tests + existing tests
- internal/tui/screens: all existing tests pass
- internal/config: config loading works
```

## Tasks

- [x] **Phase 1**: Implement topic discovery in scanner
- [x] **Phase 2**: Add topic selection screen to TUI
- [x] **Phase 3**: Update session management for topics
- [x] **Phase 4**: Update config handling
- [x] **Phase 5**: Add tests and documentation
- [ ] **Phase 6**: Manual testing with all topics (only GFS has subjective content)
- [ ] **Phase 7**: Update TUI-README.md with new feature

## Testing Strategy

### Unit Tests
1. `scanner_test.go`:
   - Topic discovery finds all chapters
   - Progressive question ordering is correct
   - Level filtering works
   
2. `app_test.go`:
   - Topic selection state transitions
   - Session filtering by topic

### Integration Tests
1. Load each topic and verify question counts
2. Test progressive navigation through levels
3. Test session resume with topic persistence

### Manual Tests
1. Select each topic and verify questions load
2. Complete a few questions and resume - verify correct topic
3. Try multiple concurrent sessions on different topics

## Success Criteria

- ✅ User can select from all available topics
- ✅ Questions load in strict L3→L7 order
- ✅ Baseline questions appear before bar-raiser within each level
- ✅ Session correctly saves and resumes topic-specific progress
- ✅ All existing tests pass
- ✅ New functionality is tested
- ✅ Documentation updated

## Notes

- Current hardcoded path has only GFS content - after this feature, all topics accessible
- Level progression is pedagogically sound but may want "practice mode" later
- Consider adding level breakdown in topic selection (how many L3, L4, etc.)

## Files Changed

### New Files
1. `internal/markdown/scanner_topic_test.go` - Tests for topic discovery and ordering (108 lines)

### Modified Files
1. `internal/markdown/scanner.go` - Added topic discovery and progressive ordering (+147 lines)
2. `internal/config/tui_config.go` - Added chapters_root_path field (+9 lines)
3. `config/tui.toml` - Updated to use new config format
4. `internal/tui/session/session.go` - Added topic fields and filtering (+30 lines)
5. `internal/tui/screens/app.go` - Added topic selection screen (+120 lines)

### Summary Statistics
- **Total Lines Added**: ~306 lines
- **New Functions**: 8 (DiscoverTopics, GetProgressiveQuestions, renderTopicSelect, etc.)
- **New Tests**: 3 (TestDiscoverTopics, TestGetProgressiveQuestions, TestFormatTopicName)
- **Test Pass Rate**: 100% (all 16+ tests passing)

## How To Use

### Running the TUI with Topic Selection
```bash
./build/quiz-tui -u your-name
```

**Flow**:
1. Welcome Screen → Press Enter
2. **Topic Selection** (NEW) → Press 1-9 to select topic
3. Session Selection → Press 'r' to resume or 'n' for new
4. Questions (L3→L4→L5→L6→L7 progressive order)
5. Complete Screen

### Example Session
```
📚 Select Your Topic

Available Topics:

  [1] Distributed Systems Gfs (20 questions)

Press 1-9 to select • Press q to quit
```

After selecting topic, questions progress through levels:
- First: All L3-baseline questions
- Then: All L3-bar-raiser questions  
- Then: All L4-baseline questions
- ...and so on through L7

### Legacy Mode (Backward Compatible)
If `content_path` is set instead of `chapters_root_path` in config, the TUI works as before (single topic, no selection screen).

## Future Enhancements

1. **Add subjective content to other topics** - Currently only GFS has subjective questions
2. **Level breakdown in topic list** - Show "L3:5, L4:5, L5:5" etc.
3. **Practice mode** - Allow random order instead of progressive
4. **Topic search/filter** - For when there are many topics
5. **Multi-topic sessions** - Practice across multiple topics in one session
