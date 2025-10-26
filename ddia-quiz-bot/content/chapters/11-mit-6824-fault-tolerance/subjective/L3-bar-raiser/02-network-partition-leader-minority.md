---
id: fault-tolerance-subjective-L3-005
type: subjective
level: L3
category: bar-raiser
topic: fault-tolerance
subtopic: network-partitions-leader-minority
estimated_time: 7-10 minutes
---

# question_title - Partition With Leader in Minority

## main_question - Core Question
"A 5-server Raft cluster partitions into {A,B} and {C,D,E}. The original leader was A (now in the minority). Describe precisely what happens in each partition during the partition and after healing. Explain leader behavior, client experience, term advancement, and log reconciliation."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Majority Elects New Leader**: {C,D,E} holds a majority and will elect a new leader
- **Minority Cannot Commit**: {A,B} cannot reach quorum; the old leader cannot commit entries
- **Term Advancement**: Majority side advances term; minority leader steps down upon seeing higher term
- **Log Reconciliation**: On heal, minority catches up by overwriting conflicts from the divergence point

### expected_keywords
- Primary: partition, majority, minority, quorum, commit
- Technical: election timeout, term number, AppendEntries, prevLogIndex/prevLogTerm

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Leader Persistence Rule**: A leader remains leader until it learns of a higher term
- **Client Behavior**: Minority-side requests time out/fail; majority-side succeed
- **Safety Over Availability**: Writes on minority not acknowledged (no majority), so no loss of committed data
- **Healing Mechanics**: Stepping down via higher-term AppendEntries/Responses

### bonus_keywords
- Implementation: nextIndex/matchIndex, backoff, retry
- Related: CAP theorem, availability trade-offs
- Operational: retry semantics, deduplication at state machine layer

## sample_excellent - Example Excellence
"When the partition occurs, {A,B} contains the old leader A but lacks quorum, so A cannot commit new entries (no majority acknowledgments). By Raft’s rules, a leader remains leader until it learns of a higher term; thus A stays leader in its view but makes no progress. On {C,D,E}, followers miss A’s heartbeats, time out, increment term, and elect a new leader (say C) via RequestVote; term advances on the majority side. Clients talking to {C,D,E} succeed; clients talking only to {A,B} see timeouts or failures for writes (no commit). When the network heals, C’s AppendEntries carrying the higher term reach A and B. They update their term, A steps down to follower, and they reconcile logs: the leader finds the last matching index/term, truncates conflicting suffixes on A/B if any, and appends the majority’s entries. Safety holds because no entry from the minority was committed. Eventually, all five servers share the same committed log."

## sample_acceptable - Minimum Acceptable
"The majority side {C,D,E} elects a new leader and continues serving requests. The minority {A,B} cannot commit writes. After healing, the minority steps down when it sees the higher term and catches up to the majority’s log."

## common_mistakes - Watch Out For
- Claiming the minority leader can commit entries or continue serving writes safely
- Saying the leader must step down immediately upon losing quorum (it steps down when it learns of a higher term)
- Ignoring client experience (timeouts on minority; success on majority)
- Forgetting the log reconciliation process and divergence handling

## follow_up_excellent - Depth Probe
**Question**: "A accepts a client write during the partition and appends it locally, but receives only B’s acknowledgment. What happens to that entry after healing?"
- **Sample answer**: "It was never committed (no majority), so if it conflicts with the majority leader’s log, it will be overwritten during reconciliation when the new leader finds the divergence point and sends AppendEntries. No acknowledged committed data is lost."

## follow_up_partial - Guided Probe  
**Question**: "Why do clients pointed at {A,B} see timeouts or failures for writes, while {C,D,E} succeed?"
- **Hint embedded**: Commit requires majority
- **Sample answer**: "Minority cannot form quorum for commits; majority can."

## follow_up_weak - Foundation Check
**Question**: "If two groups can’t talk and only one is big enough to decide, which group should make decisions, and why?"
- **Simplification**: Majority decides; guarantees single source of truth

## bar_raiser_question - L3→L4 Challenge
"Assume {C,D,E} becomes leader in term T+1 and commits entries at indices i..j. On healing, A’s log has a different entry at i. Describe the exact prevLogIndex/prevLogTerm checks and the sequence of deletions/appends needed to bring A in sync."

### bar_raiser_concepts
- Use term/index matching to find last common prefix
- Truncate conflicting suffix on A, then append leader’s entries
- After replication to majority, commitIndex can advance and apply to state machines
- No committed entries from the majority can be lost

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: Leader step-down conditions, leases and read safety, reconciliation optimizations

## assistant_answer
Majority {C,D,E} elects a new leader and advances the term; minority {A,B} cannot commit and A remains leader only until it learns of the higher term. After healing, A/B update to the higher term, A steps down, and both catch up via AppendEntries by truncating conflicts and appending the majority log, preserving safety.

## improvement_suggestions
- Add a timed sequence (timestamps/RTTs) to reason about detection and election latency.
- Include a concrete log example (indices/terms) and require the candidate to show the exact prevLogIndex/prevLogTerm probes.
- Ask for client retry/deduplication handling and how exactly-once semantics are preserved across the partition.
