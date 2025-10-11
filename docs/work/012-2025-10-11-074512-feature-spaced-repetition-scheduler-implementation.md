# Spaced Repetition Scheduler Implementation

**Date**: 2025-10-11  
**Type**: Feature  
**Status**: Planning

## 1. Overview

Implement a spaced repetition system (SRS) to help users retain concepts by resurface questions at optimal intervals. The system tracks when users answer questions, evaluates their performance, and schedules reviews based on proven spaced repetition algorithms to maximize long-term retention.

## 2. Goals

### Primary Goals
- Track user's learning history for each question
- Schedule question reviews at scientifically-backed intervals
- Focus review sessions on weak areas
- Improve long-term concept retention

### Non-Goals (Future Enhancements)
- Cross-device synchronization (can be added later)
- Social/competitive features
- AI-powered difficulty adjustment
- Custom interval configuration by users

## 3. Background: Spaced Repetition Science

### What is Spaced Repetition?
- Learning technique that schedules reviews at increasing intervals
- Based on "spacing effect" - information reviewed at intervals is retained longer
- Optimal review timing: just before you're about to forget

### Proven Algorithms
1. **SM-2 (SuperMemo 2)** - Classic, widely used in Anki
2. **SM-2+** - Enhanced version with decay factors
3. **FSRS (Free Spaced Repetition Scheduler)** - Modern, ML-optimized

### Why It Works
- Fights the forgetting curve
- Optimizes study time (don't review what you already know well)
- Builds strong long-term memory through active recall

## 4. Requirements

### 4.1 Functional Requirements

**FR1: Learning History Tracking**
- Track when user answers each question
- Store answer quality/confidence level
- Record time spent on question
- Track hint usage

**FR2: Review Scheduling**
- Calculate next review date based on performance
- Support multiple difficulty levels for same question
- Reschedule based on answer quality
- Handle both new and review questions

**FR3: Review Session Management**
- Generate daily review list
- Prioritize overdue reviews
- Mix new content with reviews
- Show progress indicators

**FR4: Performance Analytics**
- Track retention rate per concept
- Identify weak areas needing focus
- Show learning streaks
- Display upcoming review counts

**FR5: Data Persistence**
- Store learning history locally
- Backup/restore capability
- Export learning data
- Maintain history across app restarts

### 4.2 Non-Functional Requirements

**NFR1: Performance**
- Review calculation < 100ms for 1000+ cards
- Session startup < 200ms
- Smooth TUI interactions

**NFR2: Data Integrity**
- No data loss on crashes
- Atomic writes for learning history
- Validation of stored data

**NFR3: Privacy**
- All data stored locally
- No external tracking
- User controls their data

## 5. Architecture Design

### 5.1 Component Structure

```
internal/
├── srs/
│   ├── scheduler.go          # Core scheduling algorithm (SM-2+)
│   ├── card.go               # SRS card representation
│   ├── session.go            # Review session management
│   ├── history.go            # Learning history tracking
│   └── algorithm/
│       ├── sm2.go            # SM-2 algorithm implementation
│       ├── fsrs.go           # FSRS (future)
│       └── algorithm.go      # Algorithm interface
├── storage/
│   └── srs_store.go          # Persistence layer for SRS data
└── tui/
    └── screens/
        ├── review_session.go # Review session UI
        ├── stats.go          # Progress/statistics UI
        └── calendar.go       # Review calendar view

data/
└── srs/
    ├── cards.json            # Card state (intervals, ease, due dates)
    ├── history.json          # Answer history log
    └── stats.json            # Aggregated statistics
```

### 5.2 Data Models

```go
// Card represents an SRS flashcard for a question
type Card struct {
    QuestionID    string    `json:"question_id"`
    Topic         string    `json:"topic"`          // e.g., "zookeeper", "gfs"
    Level         string    `json:"level"`          // L3, L4, L5, etc.
    
    // SRS State
    Interval      int       `json:"interval"`       // Days until next review
    Repetitions   int       `json:"repetitions"`    // Number of successful reviews
    EaseFactor    float64   `json:"ease_factor"`    // Difficulty multiplier (default 2.5)
    DueDate       time.Time `json:"due_date"`       // Next review date
    LastReviewed  time.Time `json:"last_reviewed"`  // Last review timestamp
    
    // Performance Tracking
    TotalReviews  int       `json:"total_reviews"`  // All review attempts
    SuccessCount  int       `json:"success_count"`  // Good/excellent reviews
    AverageTime   int       `json:"average_time"`   // Average answer time (seconds)
    HintsUsed     int       `json:"hints_used"`     // Total hints across reviews
    
    // Lifecycle
    State         CardState `json:"state"`          // New, Learning, Review, Mature
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type CardState string
const (
    CardStateNew      CardState = "new"       // Never reviewed
    CardStateLearning CardState = "learning"  // In initial learning phase
    CardStateReview   CardState = "review"    // Regular reviews
    CardStateMature   CardState = "mature"    // Long intervals (21+ days)
    CardStateSuspended CardState = "suspended" // User paused
)

// ReviewResult captures user's performance on a question
type ReviewResult struct {
    QuestionID   string        `json:"question_id"`
    Timestamp    time.Time     `json:"timestamp"`
    TimeSpent    int           `json:"time_spent"` // seconds
    Quality      ReviewQuality `json:"quality"`    // How well answered
    HintsUsed    int           `json:"hints_used"`
    PreviousInterval int       `json:"previous_interval"`
    NewInterval     int        `json:"new_interval"`
    WasOverdue   bool          `json:"was_overdue"`
}

type ReviewQuality int
const (
    QualityBlackout   ReviewQuality = 0 // Complete failure
    QualityWrong      ReviewQuality = 1 // Incorrect answer
    QualityHard       ReviewQuality = 2 // Correct with difficulty
    QualityGood       ReviewQuality = 3 // Correct with some thought
    QualityEasy       ReviewQuality = 4 // Correct, felt easy
    QualityPerfect    ReviewQuality = 5 // Instant correct recall
)

// ReviewSession manages a study session
type ReviewSession struct {
    ID            string       `json:"id"`
    StartTime     time.Time    `json:"start_time"`
    EndTime       time.Time    `json:"end_time"`
    
    // Session Config
    MaxNewCards   int          `json:"max_new_cards"`    // New questions per session
    MaxReviews    int          `json:"max_reviews"`      // Review limit
    Topics        []string     `json:"topics"`           // Filter by topics
    Levels        []string     `json:"levels"`           // Filter by levels
    
    // Session State
    CardsReviewed []string     `json:"cards_reviewed"`
    NewCardsSeen  int          `json:"new_cards_seen"`
    ReviewsDone   int          `json:"reviews_done"`
    
    // Statistics
    TotalTime     int          `json:"total_time"`       // seconds
    AvgQuality    float64      `json:"avg_quality"`
    CardsCorrect  int          `json:"cards_correct"`
    CardsFailed   int          `json:"cards_failed"`
}

// Statistics for progress tracking
type Statistics struct {
    // Overall Stats
    TotalCards        int       `json:"total_cards"`
    NewCards          int       `json:"new_cards"`
    LearningCards     int       `json:"learning_cards"`
    ReviewCards       int       `json:"review_cards"`
    MatureCards       int       `json:"mature_cards"`
    
    // Today
    DueToday          int       `json:"due_today"`
    CompletedToday    int       `json:"completed_today"`
    NewSeenToday      int       `json:"new_seen_today"`
    
    // Upcoming
    DueTomorrow       int       `json:"due_tomorrow"`
    DueThisWeek       int       `json:"due_this_week"`
    
    // Performance
    RetentionRate     float64   `json:"retention_rate"`    // % of reviews answered well
    CurrentStreak     int       `json:"current_streak"`     // Days with reviews
    LongestStreak     int       `json:"longest_streak"`
    
    // Per-Topic Stats
    TopicStats        map[string]*TopicStatistics `json:"topic_stats"`
}

type TopicStatistics struct {
    Topic         string  `json:"topic"`
    TotalCards    int     `json:"total_cards"`
    MasteredCards int     `json:"mastered_cards"`
    RetentionRate float64 `json:"retention_rate"`
    AvgInterval   int     `json:"avg_interval"`
}
```

## 6. Algorithm Selection: SM-2+ (Enhanced SuperMemo 2)

### Why SM-2+?
- **Proven**: Used by Anki (millions of users)
- **Simple**: Easy to understand and implement
- **Effective**: Good balance between simplicity and results
- **Customizable**: Can enhance with our own tweaks

### Core SM-2 Algorithm

```
After each review, update card based on quality (0-5):

IF quality < 3 (failed):
    repetitions = 0
    interval = 1 day
    
IF quality >= 3 (passed):
    repetitions += 1
    
    IF repetitions == 1:
        interval = 1 day
    ELSE IF repetitions == 2:
        interval = 6 days
    ELSE:
        interval = previous_interval × ease_factor
    
    ease_factor = ease_factor + (0.1 - (5 - quality) × (0.08 + (5 - quality) × 0.02))
    ease_factor = max(1.3, ease_factor)  // minimum ease
```

### Our Enhancements (SM-2+)

1. **Decay Factor for Long Intervals**
   - After 30 days, reduce ease factor slightly
   - Accounts for natural forgetting on mature cards

2. **Hint Penalty**
   - Using hints reduces review quality by 1 level
   - Encourages unaided recall

3. **Time-Based Adjustments**
   - Very fast answers (< 10s) → treat as "easy"
   - Very slow answers (> 5min) → reduce quality by 1
   - Accounts for confidence in addition to correctness

4. **Overdue Adjustment**
   - Cards reviewed late get interval reduction
   - Reduction proportional to how overdue (up to 50%)

5. **New Card Graduation**
   - New cards go through: 1d → 3d → 6d before "graduating"
   - Ensures solid foundation before long intervals

## 7. Implementation Plan

### Phase 1: Core SRS Engine (Week 1)

**Day 1-2: Data Models & Storage**
- [ ] Define Card, ReviewResult, ReviewSession structs
- [ ] Implement JSON storage for cards.json, history.json
- [ ] Add CRUD operations for cards
- [ ] Write unit tests for storage layer

**Day 3-4: SM-2+ Algorithm**
- [ ] Implement base SM-2 algorithm
- [ ] Add enhancements (decay, hints, time adjustments)
- [ ] Calculate intervals and ease factors
- [ ] Unit tests with known scenarios

**Day 5: Session Management**
- [ ] Implement review session creation
- [ ] Build due card selection logic (overdue + new)
- [ ] Add session state management
- [ ] Test session workflows

### Phase 2: TUI Integration (Week 2)

**Day 6-7: Review Session UI**
- [ ] Create ReviewSessionScreen
- [ ] Show question + quality buttons (0-5)
- [ ] Display progress bar (X of Y cards)
- [ ] Show streak and daily goal
- [ ] Handle card transitions

**Day 8-9: Statistics Dashboard**
- [ ] Create StatsScreen
- [ ] Display overall statistics
- [ ] Show per-topic breakdown
- [ ] Add retention rate charts (text-based)
- [ ] Display upcoming reviews calendar

**Day 10: Home Screen Integration**
- [ ] Add "Daily Review" option to home screen
- [ ] Show review count badge
- [ ] Display streak counter
- [ ] Add "Quick Stats" widget

### Phase 3: Polish & Testing (Week 3)

**Day 11-12: User Experience**
- [ ] Add keyboard shortcuts for review quality (1-5 keys)
- [ ] Implement undo last review
- [ ] Add session pause/resume
- [ ] Show concept reminders before review

**Day 13-14: Data Management**
- [ ] Implement backup/restore
- [ ] Add data export (CSV/JSON)
- [ ] Create data validation/repair tool
- [ ] Handle migration from old format

**Day 15: Testing & Documentation**
- [ ] Integration tests for full review cycle
- [ ] Load testing with 1000+ cards
- [ ] Write user documentation
- [ ] Create tutorial/onboarding

## 8. User Workflows

### 8.1 First-Time User Flow

```
1. User launches TUI
2. Sees "Start Daily Review" option with badge: "10 new questions available"
3. Selects review mode
4. Sees onboarding: "Let's learn with spaced repetition!"
5. System shows 5 new questions (default new card limit)
6. After each answer, user rates quality (1-5)
7. System schedules next review
8. End of session: "Great start! See you tomorrow for 5 reviews."
```

### 8.2 Regular User Flow

```
1. User launches TUI
2. Home screen shows:
   - "Daily Review: 12 cards due (3 new, 9 reviews)"
   - "Current streak: 7 days 🔥"
   - "Retention rate: 87%"
3. User starts review session
4. System prioritizes: overdue > due today > new
5. For each card:
   - Show question
   - User answers
   - User rates quality (Easy/Good/Hard/Again)
   - System shows next review date
6. End of session:
   - "Session complete! ✓"
   - "Reviewed: 12 cards in 8 minutes"
   - "See you tomorrow: 8 cards due"
```

### 8.3 Review Session UI

```
┌────────────────────────────────────────────────────────────────┐
│ Daily Review Session                      [6/12] 🔥 Streak: 7d │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│ Topic: ZooKeeper • Level: L3 • Category: Baseline              │
│                                                                 │
│ Question:                                                       │
│ Explain what linearizability means in the context of           │
│ ZooKeeper. Why is it important for a coordination service?     │
│                                                                 │
│ Core Concepts:                                                  │
│ • Linearizability Definition                                    │
│ • Strong Consistency                                            │
│ • Coordination Service Need                                     │
│                                                                 │
│ [Your answer area - press Enter to submit]                     │
│                                                                 │
├────────────────────────────────────────────────────────────────┤
│ Rate your answer:                                               │
│                                                                 │
│  [1] Again      - Need to relearn                              │
│  [2] Hard       - Correct but struggled                        │
│  [3] Good       - Correct with some thought                    │
│  [4] Easy       - Correct, felt confident                      │
│  [5] Perfect    - Instant recall                               │
│                                                                 │
│ [H] Hint • [S] Sample Answer • [ESC] Pause                     │
└────────────────────────────────────────────────────────────────┘
```

### 8.4 Statistics Dashboard

```
┌────────────────────────────────────────────────────────────────┐
│ Learning Statistics                                             │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│ Overall Progress                                                │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 156/200     │
│                                                                 │
│ Card States:                                                    │
│ • New:        44 cards                                          │
│ • Learning:   23 cards (intervals < 21 days)                   │
│ • Review:     67 cards                                          │
│ • Mature:     22 cards (intervals > 21 days)                   │
│                                                                 │
│ Today:                                                          │
│ • Due:         12 cards (3 new, 9 reviews)                     │
│ • Completed:    8 cards                                         │
│ • Time spent:   12 minutes                                      │
│                                                                 │
│ Performance:                                                    │
│ • Retention:    87% (good: 132/152 reviews)                    │
│ • Streak:       7 days 🔥                                       │
│ • Longest:      12 days                                         │
│                                                                 │
│ Upcoming Reviews:                                               │
│ • Tomorrow:     8 cards                                         │
│ • This week:    34 cards                                        │
│ • Next week:    28 cards                                        │
│                                                                 │
├────────────────────────────────────────────────────────────────┤
│ By Topic:                                                       │
│                                                                 │
│ ZooKeeper       ████████████████░░░░ 81% (13/16) Avg: 8d      │
│ GFS             ██████████████░░░░░░ 72% (18/25) Avg: 6d      │
│ Raft            ████████████████████ 93% (14/15) Avg: 12d     │
│ Transactions    ██████████░░░░░░░░░░ 64% (32/50) Avg: 4d      │
│                                                                 │
│ [R] Start Review • [ESC] Back                                  │
└────────────────────────────────────────────────────────────────┘
```

## 9. File Locations & Persistence

### 9.1 Data Directory Structure

```
~/.ddia-clicker/
└── srs/
    ├── cards.json          # Current state of all cards
    ├── history.json        # Complete review history (append-only log)
    ├── stats.json          # Pre-calculated statistics (cache)
    ├── sessions/           # Individual session logs
    │   ├── 2025-10-11.json
    │   ├── 2025-10-12.json
    │   └── ...
    └── backups/            # Automatic backups
        ├── cards-2025-10-11.json.bak
        └── ...
```

### 9.2 Data Format Examples

**cards.json:**
```json
{
  "version": "1.0",
  "updated_at": "2025-10-11T14:30:00Z",
  "cards": [
    {
      "question_id": "zookeeper-subjective-L3-001",
      "topic": "zookeeper",
      "level": "L3",
      "interval": 7,
      "repetitions": 3,
      "ease_factor": 2.6,
      "due_date": "2025-10-18T00:00:00Z",
      "last_reviewed": "2025-10-11T14:25:00Z",
      "total_reviews": 4,
      "success_count": 3,
      "average_time": 180,
      "hints_used": 1,
      "state": "review",
      "created_at": "2025-10-05T10:00:00Z",
      "updated_at": "2025-10-11T14:25:00Z"
    }
  ]
}
```

**history.json (append-only log):**
```json
{
  "version": "1.0",
  "reviews": [
    {
      "question_id": "zookeeper-subjective-L3-001",
      "timestamp": "2025-10-11T14:25:00Z",
      "time_spent": 180,
      "quality": 4,
      "hints_used": 0,
      "previous_interval": 3,
      "new_interval": 7,
      "was_overdue": false
    }
  ]
}
```

### 9.3 Backup Strategy

- **Automatic backups**: Daily before first review
- **Retention**: Keep last 7 days
- **On-demand**: Before major operations (bulk operations, imports)
- **Format**: Same as main files with `.bak` extension + timestamp

## 10. Algorithm Implementation Details

### 10.1 Review Quality Mapping

For subjective questions without auto-grading, we use self-assessment:

| Quality | Score | Description | Interval Effect |
|---------|-------|-------------|-----------------|
| Perfect | 5 | Instant recall, comprehensive answer | +100% interval |
| Easy | 4 | Correct, felt confident | +80% interval |
| Good | 3 | Correct with some thought | +60% interval (base) |
| Hard | 2 | Correct but struggled significantly | Repeat with 1d interval |
| Wrong | 1 | Incorrect or incomplete | Reset to learning |
| Blackout | 0 | No recall at all | Reset to new |

### 10.2 Interval Calculation Pseudocode

```go
func CalculateNextInterval(card *Card, quality ReviewQuality, timeSpent int, hintsUsed int) int {
    // Adjust quality based on context
    adjustedQuality := quality
    
    // Hint penalty
    if hintsUsed > 0 {
        adjustedQuality = max(0, adjustedQuality - 1)
    }
    
    // Time-based adjustment
    if timeSpent < 10 && quality >= 3 {
        adjustedQuality = min(5, adjustedQuality + 1) // Very fast = easy
    } else if timeSpent > 300 && quality >= 3 {
        adjustedQuality = max(0, adjustedQuality - 1) // Very slow = harder
    }
    
    // Failed review
    if adjustedQuality < 3 {
        card.Repetitions = 0
        card.EaseFactor = max(1.3, card.EaseFactor - 0.2)
        return 1 // Review again tomorrow
    }
    
    // Successful review
    card.Repetitions++
    
    // Calculate new interval based on repetitions
    var interval int
    switch card.Repetitions {
    case 1:
        interval = 1
    case 2:
        interval = 6
    default:
        interval = int(float64(card.Interval) * card.EaseFactor)
        
        // Apply decay for mature cards (30+ days)
        if interval >= 30 {
            decayFactor := 0.95
            interval = int(float64(interval) * decayFactor)
        }
    }
    
    // Update ease factor using SM-2 formula
    card.EaseFactor = card.EaseFactor + (0.1 - (5 - adjustedQuality) * (0.08 + (5 - adjustedQuality) * 0.02))
    card.EaseFactor = max(1.3, card.EaseFactor)
    
    // Handle overdue reviews
    if card.IsOverdue() {
        overdueDays := card.DaysSinceLastReview() - card.Interval
        reductionFactor := min(0.5, float64(overdueDays) / float64(card.Interval) * 0.5)
        interval = int(float64(interval) * (1.0 - reductionFactor))
    }
    
    return max(1, interval)
}
```

### 10.3 Due Card Selection Logic

```go
func SelectDueCards(maxNew, maxReviews int, topics, levels []string) []Card {
    // Priority order:
    // 1. Overdue reviews (past due date)
    // 2. Due today reviews
    // 3. New cards (up to maxNew)
    
    overdue := GetOverdueCards(topics, levels)
    dueToday := GetDueTodayCards(topics, levels)
    newCards := GetNewCards(maxNew, topics, levels)
    
    // Combine with limits
    selected := []Card{}
    
    // All overdue (critical)
    selected = append(selected, overdue...)
    
    // Due today up to review limit
    remaining := maxReviews - len(overdue)
    if remaining > 0 {
        selected = append(selected, dueToday[:min(remaining, len(dueToday))]...)
    }
    
    // New cards if under total limit
    if len(selected) < maxReviews {
        remaining = min(maxNew, maxReviews - len(selected))
        selected = append(selected, newCards[:min(remaining, len(newCards))]...)
    }
    
    // Shuffle to mix topics
    shuffle(selected)
    
    return selected
}
```

## 11. Configuration & Defaults

### 11.1 System Defaults

```yaml
# Default SRS configuration
srs:
  # Session limits
  max_new_cards_per_day: 5          # Conservative start
  max_reviews_per_day: 50           # Reasonable daily load
  
  # Timing
  new_card_intervals: [1, 3, 6]     # Days for graduating new cards
  min_interval: 1                    # Minimum 1 day
  max_interval: 365                  # Maximum 1 year
  
  # Algorithm
  starting_ease: 2.5                 # Default ease factor
  min_ease: 1.3                      # Floor for ease factor
  ease_bonus: 0.15                   # Bonus for quality 5
  ease_penalty: 0.20                 # Penalty for failed review
  
  # Decay
  mature_threshold: 21               # Days to be "mature"
  decay_enabled: true
  decay_threshold: 30                # Apply decay after 30 days
  decay_factor: 0.95
  
  # Adjustments
  hint_penalty: 1                    # Quality reduction per hint
  fast_answer_threshold: 10          # Seconds for "easy" bonus
  slow_answer_threshold: 300         # Seconds for difficulty penalty
  overdue_reduction_max: 0.5         # Max 50% interval reduction
```

### 11.2 User Preferences (Configurable)

```yaml
# User can customize
preferences:
  daily_goal: 20                     # Cards per day target
  session_style: "mixed"             # "mixed", "new_first", "reviews_first"
  auto_play_audio: false             # For future audio features
  show_hints_by_default: false
  
  # Topics to focus on
  enabled_topics: 
    - "zookeeper"
    - "gfs"
    - "raft"
  
  # Notification preferences
  daily_reminder: true
  reminder_time: "09:00"
```

## 12. Testing Strategy

### 12.1 Unit Tests

1. **Algorithm Tests**
   - Test SM-2 calculations with known inputs/outputs
   - Verify ease factor adjustments
   - Test interval calculations for all quality levels
   - Test edge cases (quality 0, quality 5, hints, overdue)

2. **Storage Tests**
   - Test CRUD operations
   - Test concurrent access
   - Test data corruption recovery
   - Test backup/restore

3. **Session Tests**
   - Test card selection logic
   - Test session state management
   - Test statistics calculation

### 12.2 Integration Tests

1. **Full Review Cycle**
   - New card → learning → review → mature
   - Failed review → reset → relearn
   - Overdue handling

2. **Multi-Session Flow**
   - Review across multiple days
   - Verify interval progression
   - Check statistics accuracy

3. **Data Persistence**
   - Save/load state
   - Handle app restart
   - Recover from crashes

### 12.3 User Acceptance Testing

1. **First-Time Experience**
   - Can user understand SRS concept?
   - Is onboarding clear?
   - Are quality ratings intuitive?

2. **Regular Use**
   - Is daily workflow smooth?
   - Are statistics motivating?
   - Is progress visible?

3. **Edge Cases**
   - Missed days (overdue cards)
   - Bulk learning (many new cards)
   - Topic switching

## 13. Metrics & Success Criteria

### 13.1 System Metrics

- **Algorithm Accuracy**: Retention rate > 80% for mature cards
- **Performance**: Review calculation < 100ms
- **Reliability**: Zero data loss incidents

### 13.2 User Metrics

- **Engagement**: Daily active users
- **Retention**: % of users returning day 2, 7, 30
- **Learning**: Average cards mastered per week
- **Satisfaction**: Self-reported confidence improvement

### 13.3 Success Criteria

Phase 1 (MVP):
- [ ] Can track 100+ cards without performance issues
- [ ] Correct interval calculations (verified against Anki)
- [ ] Data persists across restarts
- [ ] Basic statistics display

Phase 2 (Enhanced):
- [ ] Smooth TUI integration
- [ ] Users understand quality ratings
- [ ] Statistics are motivating
- [ ] Daily workflow < 15 minutes

Phase 3 (Production):
- [ ] 80%+ retention rate for mature cards
- [ ] Users maintain 7+ day streaks
- [ ] Positive user feedback on learning outcomes

## 14. Future Enhancements

### Phase 4+ (Post-Launch)

1. **Advanced Algorithm (FSRS)**
   - ML-based optimization
   - Personalized difficulty curves
   - Cross-topic learning transfer

2. **Analytics Dashboard**
   - Learning velocity charts
   - Concept mastery heat maps
   - Forecast future review load

3. **Social Features**
   - Leaderboards (optional)
   - Study groups
   - Shared decks

4. **Cross-Device Sync**
   - Cloud backup (optional)
   - Multi-device support
   - Web interface

5. **AI Integration**
   - Auto-grade subjective answers
   - Suggest related questions
   - Generate follow-up questions

## 15. Implementation Status

### Phase 1: Core SRS Engine ✅ COMPLETED
- [x] Data models defined (card.go, session.go)
- [x] Storage layer implemented (storage.go)
- [x] SM-2+ algorithm implemented (algorithm/sm2.go)
- [x] Session management implemented (scheduler.go)
- [x] Unit tests written (card_test.go, algorithm/sm2_test.go)
- [x] Scheduler with full card management
- [x] Statistics calculation
- [x] Backup functionality

**Completed**: 2025-10-11

### Phase 2: TUI Integration 🚧 IN PROGRESS
- [ ] Review session UI created
- [ ] Statistics dashboard created
- [ ] Home screen integrated
- [ ] Keyboard shortcuts added

### Phase 3: Polish & Testing
- [ ] User experience refined
- [ ] Data management features added
- [ ] Integration tests passing
- [ ] Documentation written

## 16. Phase 1 Summary

### What Was Implemented

#### Core Data Structures
1. **Card** (`card.go`) - Complete SRS card with:
   - State management (New, Learning, Review, Mature, Suspended)
   - Performance tracking (retention rate, average time)
   - Due date calculation
   - Overdue detection
   
2. **ReviewSession** (`session.go`) - Session management with:
   - Configurable limits (max new cards, max reviews)
   - Topic and level filtering
   - Session statistics (quality, time, correct/failed)
   
3. **Statistics** (`session.go`) - Comprehensive stats:
   - Card counts by state
   - Due cards (today, tomorrow, this week)
   - Retention rates
   - Streak tracking
   - Per-topic statistics

#### Storage Layer (`storage.go`)
- JSON-based persistence in ~/.ddia-clicker/srs/
- Atomic writes with temp file + rename
- Separate files for cards, history, and sessions
- Automatic daily backups (keeps last 7 days)
- Thread-safe with mutex protection

#### SM-2+ Algorithm (`algorithm/sm2.go`)
- Enhanced SuperMemo 2 algorithm
- Graduated intervals for new cards: 1d → 3d → 6d
- Quality-based ease factor adjustment
- Hint penalty (reduces effective quality)
- Timing adjustments (very fast/slow answers)
- Decay factor for mature cards (30+ days)
- Overdue reduction (proportional to delay)
- Configurable parameters via Config struct

#### Scheduler (`scheduler.go`)
- Central orchestrator for all SRS operations
- Add questions individually or in bulk
- Record reviews with full tracking
- Get due cards with priority: overdue > due today > new
- Topic/level filtering
- Card state management (reset, suspend, unsuspend)
- Statistics generation
- Session persistence

### File Structure
```
internal/srs/
├── card.go              # Card data model and state management
├── session.go           # Session and statistics models
├── storage.go           # JSON persistence layer
├── scheduler.go         # Main SRS orchestrator
├── card_test.go         # Card unit tests
├── algorithm/
│   ├── algorithm.go     # Algorithm interface
│   ├── sm2.go          # SM-2+ implementation
│   └── sm2_test.go     # Algorithm tests
```

### Tests Passing
- ✅ Card creation and state management
- ✅ Card review tracking
- ✅ SM-2+ graduation intervals
- ✅ Algorithm ease factor adjustments
- (More tests pending for storage and scheduler)

### Next Steps for Phase 2
1. Create TUI review session screen
2. Implement interactive review flow
3. Add statistics dashboard
4. Integrate with home screen
5. Add keyboard shortcuts for quick quality rating

## 16. Notes

- Start simple: Focus on SM-2+ before considering FSRS
- User education is key: Explain how SRS works and why quality ratings matter
- Daily habit formation: Make it easy to maintain streak
- Data integrity: Never lose user progress - it's demotivating
- Performance: Must feel instant even with 1000+ cards
- Motivation: Make progress visible and celebrate achievements

## 17. References

- **SM-2 Algorithm**: https://www.supermemo.com/en/archives1990-2015/english/ol/sm2
- **Anki Manual**: https://docs.ankiweb.net/studying.html
- **FSRS**: https://github.com/open-spaced-repetition/fsrs4anki
- **Spaced Repetition Research**: https://en.wikipedia.org/wiki/Spaced_repetition
- **Forgetting Curve**: https://en.wikipedia.org/wiki/Forgetting_curve
