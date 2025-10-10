# Quiz TUI - Interactive Subjective Quiz Terminal Interface

An interactive terminal-based quiz application for answering subjective questions from the DDIA quiz corpus.

## Features

- 📝 **Interactive Question Flow**: Linear progression through subjective questions (L3-L7)
- 💾 **Auto-Save**: Automatic saves every 30 seconds + save on navigation
- 🔄 **Resume Sessions**: Pick up where you left off with session management
- 🎨 **Beautiful UI**: Markdown rendering, styled borders, and intuitive controls
- ⌨️ **Keyboard-First**: Efficient keyboard shortcuts for all actions
- 📊 **Progress Tracking**: Visual progress indicator and session statistics

## Quick Start

### Build

```bash
./scripts/build_tui.sh
```

This creates the binary at `build/quiz-tui`.

### Run

```bash
./build/quiz-tui --user <your-username>
```

Example:
```bash
./build/quiz-tui --user abhishek
```

## Usage

### Starting a New Quiz

1. Launch the application with your username
2. Press `Enter` at the welcome screen
3. If you have incomplete sessions, choose:
   - `r` to resume the most recent session
   - `n` to start a new session
4. Begin answering questions!

### Answering Questions

- Type your answer in the text area
- Use standard text editing keys:
  - Arrow keys, Home, End for navigation
  - Backspace, Delete for editing
  - Enter for new lines within your answer
- Your answer auto-saves every 30 seconds
- Progress is shown at the top: "Question X of Y"

### Navigation & Controls

| Key | Action |
|-----|--------|
| `Ctrl+N` or `Ctrl+Enter` | Save current answer and move to next question |
| `Ctrl+S` | Manually save current answer |
| `Ctrl+C` | Quit (auto-saves current answer) |
| `q` | Quit (on welcome/complete screens) |
| `r` | Resume most recent incomplete session |
| `n` | Start new session |

### Completing the Quiz

- After answering the last question, the session automatically completes
- Your answers are saved in `sessions/<username>/subjective/`
- Press `q` to exit

## Configuration

Edit `config/tui.toml` to customize:

```toml
# Auto-save interval in seconds
auto_save_interval = 30

# Directory where session files are stored
sessions_dir = "sessions"

# Path to the content directory for subjective questions
content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"
```

## Session Files

Sessions are stored as JSON in `sessions/<user>/subjective/`:

```
sessions/
└── abhishek/
    └── subjective/
        └── 20251011-101530-abhishek-subjective.json
```

Each session file contains:
- Session metadata (user, status, timestamps, progress)
- List of questions in the quiz
- Your responses with timestamps and time spent

## Question Content

Questions are organized by difficulty level:
- **L3**: Baseline and bar-raiser questions
- **L4**: Baseline and bar-raiser questions
- **L5**: Baseline and bar-raiser questions
- **L6**: Baseline questions
- **L7**: Baseline and bar-raiser questions

Topics include:
- GFS (Google File System)
  - Replication
  - Consistency models
  - Chunk design
- Raft Consensus
  - Introduction and basics
  - Leader election
  - Log replication
  - Performance and evolution

## Tips

1. **Take Your Time**: There's no time limit on questions
2. **Auto-Save**: Your work is saved regularly, but you can manually save with `Ctrl+S`
3. **Resume Anytime**: If you need to take a break, just quit with `Ctrl+C` and resume later
4. **Linear Flow**: Questions must be answered in order (no skipping or going back)
5. **Markdown Support**: Questions may contain formatted text, code blocks, and lists

## Troubleshooting

### Binary Not Found
```bash
# Make sure you've built the project
./scripts/build_tui.sh

# Or build manually
go build -o build/quiz-tui ./cmd/quiz-tui
```

### Config File Not Found
The application uses default settings if `config/tui.toml` is not found. You can create it by copying the example in this README.

### No Questions Found
Ensure the content path in your config points to the correct location:
```toml
content_path = "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective"
```

### Session Won't Resume
Check that your username matches and the session file exists in `sessions/<username>/subjective/`.

## Architecture

The TUI is built with:
- **Bubbletea**: TUI framework with Elm-inspired architecture
- **Bubbles**: Pre-built UI components (textarea)
- **Lipgloss**: Terminal styling and layouts
- **Glamour**: Markdown rendering
- **Cobra**: CLI framework
- **Viper**: Configuration management

## Development

### Project Structure
```
cmd/quiz-tui/           # Main entry point
internal/
  ├── models/          # Question and response data structures
  ├── markdown/        # Question file parsing
  ├── config/          # Configuration management
  └── tui/
      ├── screens/     # Main application logic
      ├── components/  # Reusable UI widgets
      └── session/     # Session management
```

### Adding Features

The codebase is structured to support future enhancements:
- Objective/MCQ mode (architecture ready)
- Back navigation
- Review screen before submission
- In-app evaluation integration

## License

Part of the DDIA Clicker project.

## Support

For issues or questions, refer to the main project documentation or open an issue in the repository.
