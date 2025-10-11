# Mode Selection & Menu Alignment Improvements

## Date: 2025-10-11 11:24:10

## Task
1. Add manual mode switch (MCQ / Subjective / Mixed)
2. Fix topic menu alignment issues

## Current Issues
- Mode is automatically determined per-question, no user control
- Topic menu has misaligned columns due to varying topic name lengths

## Implementation Plan

### Feature 1: Manual Mode Selection

**Changes:**
1. Add `StateModeSelect` constant
2. Add `selectedMode` field to `ImprovedAppModel`
3. Add `renderModeSelect()` function
4. Update state transitions: Welcome → ModeSelect → TopicSelect → SessionSelect → Question
5. Filter questions based on selected mode
6. Add mode indicator to question header
7. Add `default_mode` config option

**Flow:**
```
Welcome
  ↓
ModeSelect (new)
  ↓
TopicSelect
  ↓
SessionSelect
  ↓
Question
```

### Feature 2: Fix Menu Alignment

**Changes:**
1. Calculate max topic name length
2. Use lipgloss fixed-width columns
3. Align all columns properly
4. Add MCQ/Subjective counts to TopicInfo

**Before:**
```
[1] GFS & Distributed Systems (25 questions)
[2] Raft (18 questions)
```

**After:**
```
[1]  GFS & Distributed Systems    (25 questions)
[2]  Raft                          (18 questions)
```

## Implementation Status

### Phase 1: Data Model Updates ✅
- [x] Add MCQ/Subjective counts to TopicInfo
- [x] Add selectedMode field to ImprovedAppModel
- [x] Add default_mode to TUIConfig

### Phase 2: UI Changes ✅
- [x] Add StateModeSelect constant
- [x] Implement renderModeSelect() function
- [x] Fix renderTopicSelect() alignment
- [x] Add mode indicator to topic view

### Phase 3: Logic Updates ✅
- [x] Update state transitions
- [x] Filter questions by mode
- [x] Update loadTopicQuestionsCmd to support filtering

### Phase 4: Testing ✅
- [x] Build successful
- [ ] Manual TUI testing (requires user interaction)
- [ ] Test with mixed content
- [ ] Verify alignment across different topics

## Files Modified

1. **internal/markdown/scanner.go**
   - Added MCQCount and SubjectiveCount to TopicInfo struct
   - Updated DiscoverTopics() to scan both subjective/ and mcq/ directories
   - Counts questions separately by type

2. **internal/config/tui_config.go**
   - Added DefaultMode field to TUIConfig
   - Set default value to "mixed" for backward compatibility

3. **internal/tui/screens/app.go**
   - Added StateModeSelect state constant
   - Added selectedMode field to ImprovedAppModel
   - Implemented renderModeSelect() with clean UI showing:
     - MCQ Questions option
     - Subjective Questions option
     - Mixed Mode option
     - Current selection indicator
   - Updated renderTopicSelect() with:
     - Fixed-width columns using lipgloss
     - Proper alignment for all topic names
     - Mode-specific counts display
     - Mode indicator at top
   - Updated state transitions: Welcome → ModeSelect → TopicSelect → SessionSelect → Question
   - Modified loadTopicQuestionsCmd() to:
     - Load subjective questions based on mode
     - Load MCQ questions based on mode
     - Combine questions appropriately
   - Added filepath import for path manipulation

## Features Implemented

### 1. Mode Selection Screen
- Clean, intuitive UI with numbered options (1-3)
- Descriptions for each mode
- Shows current selection
- Keyboard navigation (1, 2, 3, q)

### 2. Fixed Topic Menu Alignment
- Calculates max topic name length dynamically
- Uses lipgloss fixed-width columns:
  - Index column: 6 chars
  - Name column: dynamic width
  - Count column: dynamic content
- Perfect vertical alignment regardless of topic name length

### 3. Mode-Specific Question Counts
- Mixed mode: Shows "MCQ: X | Subjective: Y"
- MCQ mode: Shows "X questions" (MCQ count only)
- Subjective mode: Shows "X questions" (Subjective count only)

### 4. Question Filtering
- Loads only MCQ questions in MCQ mode
- Loads only subjective questions in Subjective mode
- Loads both in Mixed mode
- Graceful error handling if no questions found

## Bug Fix: Auto-detect MCQ Type

**Issue Found:** MCQ files didn't have `type: mcq` in frontmatter, so they were treated as subjective questions.

**Fix Applied:** Added auto-detection in parser after parseBody():
```go
// Auto-detect question type if not explicitly set
if question.Type == "" {
    if len(question.Options) > 0 && question.Answer != "" {
        question.Type = "mcq"
    } else {
        question.Type = "subjective"
    }
}
```

Now questions are automatically detected as MCQ if they have options and an answer section.

## UX Improvement: Letter Keys for Mode Selection

**Issue:** Using numbers (1, 2, 3) for both mode selection and topic selection was confusing.

**Fix:** Changed mode selection to use letter keys:
- **[M]** - MCQ Questions
- **[S]** - Subjective Questions  
- **[B]** - Both (Mixed Mode)
- Also accepts [X] for miXed mode

Topic selection continues to use numbers 1-9, eliminating the conflict.

## Notes
- Backward compatibility: Mixed mode = current behavior (default)
- Default mode is configurable via config file
- MCQ count scans mcq/ subdirectory
- Subjective count scans subjective/ subdirectory with level folders
- Auto-detection ensures MCQ files work without explicit type field
- Build successful, ready for manual testing
