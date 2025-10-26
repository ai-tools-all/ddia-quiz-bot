---
id: spanner-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: spanner
subtopic: truetime-commit-wait
estimated_time: 6-8 minutes
---

# question_title - TrueTime and Commit-Wait

## main_question - Core Question
"Explain Spanner’s TrueTime API and how commit-wait uses it to achieve external consistency. Give a simple example that shows why the wait is necessary."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- TrueTime returns an interval [earliest, latest] capturing clock uncertainty
- Commit timestamp chosen using latest bound; must wait until earliest > ts
- External consistency: real-time order respected by serialization order
- Bounded uncertainty (e.g., ~7–10ms) drives wait duration

### expected_keywords
- Primary: TrueTime, uncertainty, commit-wait, timestamp
- Technical: TT.now().earliest/latest, external consistency, monotonicity

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- GPS/atomic clocks as time sources
- Effects on latency of write transactions
- Why read-only transactions don’t need locks
- Failure/scenario where skipping wait breaks ordering

### bonus_keywords
- time drift, bound ε, snapshot reads, serialization order

## sample_excellent - Example Excellence
"TrueTime exposes bounded uncertainty via [earliest, latest]. A transaction chooses a commit timestamp ts≈TT.latest; to ensure wall-clock order matches serialization, Spanner waits until TT.earliest > ts before releasing locks. If T1 commits then T2 starts later, T2’s ts > T1’s ts by construction, so T2 observes T1. Without the wait, a commit could publish with a timestamp in the future, letting a later-starting transaction read a snapshot before T1’s ts and miss T1’s write, violating external consistency."

## sample_acceptable - Minimum Acceptable
"TrueTime gives a time interval. Spanner waits a small amount so the commit timestamp is definitely in the past. That preserves real-time ordering."

## common_mistakes - Watch Out For
- Treating TrueTime as a single precise timestamp
- Saying commit-wait waits for all replicas
- Confusing external consistency with linearizability on a single shard only

## follow_up_excellent - Depth Probe
**Question**: "How would you tune commit-wait if uncertainty bounds vary by region?"
- **Looking for**: Using dynamic ε per data center, impact on tail latency

## follow_up_partial - Guided Probe  
**Question**: "If TT.latest is ahead of actual time, what goes wrong without the wait?"
- **Hint embedded**: Think about a commit timestamp in the future

## follow_up_weak - Foundation Check
**Question**: "What does it mean that time is uncertain in distributed systems?"
