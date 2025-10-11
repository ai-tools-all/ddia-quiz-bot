## 2025-10-11: TUI Confidence Tests

### Goals
- Capture checklist for high-confidence TUI tests
- Implement coverage for topic discovery, question ordering, autosave, resume, and quit flows

### Checklist
- [x] Draft detailed plan for five critical tests
- [x] Implement topic discovery navigation test
- [x] Implement progressive question ordering test
- [x] Implement auto-save cadence test using short intervals
- [x] Implement resume-most-recent-session test with prepared sessions
- [x] Implement graceful quit persistence test
- [x] Run `go test ./...`
- [x] Prepare concise summary & commit

### Notes
- Use temporary directories for session manager interactions to avoid polluting real sessions
- Mock keyboard events with `tea.KeyMsg` versions matching shortcuts
- Reduce `AutoSaveInterval` within tests to milliseconds to keep runtime short
