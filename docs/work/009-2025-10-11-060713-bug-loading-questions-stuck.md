# Bug: TUI Gets Stuck on "Loading Questions" Screen

**Date**: 2025-10-11
**Type**: bug
**Status**: Test Written & Issue Identified

## Problem Description
The TUI application gets stuck displaying "Loading Questions..." on the welcome screen and doesn't proceed to topic selection.

## Root Cause Analysis
After investigation and writing tests, the issue was identified:

1. **State Transition Bug**: The app remains in `StateWelcome` (state 0) even after topics are successfully discovered
2. **Welcome Screen Logic Issue**: The `renderWelcome()` function checks for `m.questions == nil` to show "Loading questions...", but in topic selection mode, questions aren't loaded until after a topic is selected
3. **Missing State Update**: After receiving `topicsDiscoveredMsg`, the code doesn't transition to `StateTopicSelect` when topics are available

## Test Coverage Added
Created comprehensive test file: `/internal/tui/screens/app_test.go`

### Tests Written:
1. **TestQuestionLoadingFromSubjectiveFolder**: ✅ PASS
   - Verifies questions load correctly from GFS subjective folder
   - Tests handling of non-existent folders

2. **TestTopicDiscoveryWithChaptersRoot**: ✅ PASS
   - Confirms topic discovery finds GFS chapter with 20 questions
   - Verifies only GFS has subjective questions currently

3. **TestAppModelQuestionLoading**: ❌ FAIL (Catches the bug!)
   - Demonstrates app stays in StateWelcome after topic discovery
   - This test failure confirms the user's reported issue

4. **TestWelcomeScreenLoadingState**: ✅ PASS
   - Tests welcome screen rendering with/without questions

5. **TestSingleTopicModeFallback**: ✅ PASS
   - Tests fallback to single topic mode when chapters_root_path is empty

## Additional Findings
- Only `09-distributed-systems-gfs` has a `subjective/` folder
- Other chapters (03-08) lack subjective question directories
- The scanner correctly handles missing directories but the UI state machine doesn't update properly

## Workaround
Users can switch to single topic mode by editing `config/tui.toml`:
```toml
# Comment out chapters_root_path to disable topic selection
# chapters_root_path = "ddia-quiz-bot/content/chapters"

# Use single topic mode instead
content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"
```

## Fix Applied
The bug has been fixed with the following changes:

1. **State Transition Fix** (`app.go`):
   - Added state transition to `StateTopicSelect` after topics are discovered
   - Now properly transitions from Welcome → Topic Select when topics are loaded

2. **Welcome Screen Rendering** (`app.go`):
   - Updated to show "Discovering topics..." in topic mode
   - Shows "Loading questions..." only in single topic mode
   - Displays appropriate info based on mode (topics vs questions)

3. **Test Updates** (`app_test.go`):
   - Updated `TestWelcomeScreenLoadingState` to handle both topic and single topic modes
   - All tests now pass, including the previously failing `TestAppModelQuestionLoading`

## Verification
✅ All tests pass
✅ TUI builds successfully
✅ State transition from Welcome to Topic Select works correctly
✅ Welcome screen shows appropriate loading messages based on mode
