---
id: cops-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: cops
subtopic: conflict-resolution
estimated_time: 10-12 minutes
---

# question_title - Concurrent Writes and Conflict Resolution

## main_question - Core Question
"How does COPS handle concurrent writes to the same key at different data centers? Explain the role of Lamport clocks and last-writer-wins, and discuss what application-level problems might remain."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Concurrent writes = updates to same key without causal relationship
- COPS uses Lamport clocks for version numbering
- Last-writer-wins (LWW): higher Lamport timestamp wins
- LWW discards one update, causing potential data loss
- Application may need custom conflict resolution (COPS doesn't provide)

### expected_keywords
- concurrent writes, Lamport clock, last-writer-wins, data loss, conflict resolution

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Causal consistency doesn't order concurrent operations
- LWW is simple but not always semantically correct
- Alternative: CRDTs, application-specific merge functions
- Difference from causally-ordered vs. concurrent writes

### bonus_keywords
- commutativity, CRDT, semantic conflicts, merge strategies, concurrent vs causal

## sample_excellent - Example Excellence
"When two clients concurrently update the same key at different data centers (neither read the other's write, so no causal relationship), COPS uses Lamport clocks to assign version numbers and applies last-writer-wins (LWW) based on highest timestamp. This means one write is discarded, potentially losing data. For example, if two users simultaneously update a shopping cart at different locations, one update is lost. Application-level problems include: lost updates (cart items disappear), incorrect semantics (counter decrements lost), and no merge capability. COPS doesn't solve these conflicts—applications may need CRDTs or custom merge logic for commutative updates."

## sample_acceptable - Minimum Acceptable
"Concurrent writes to the same key use Lamport clocks, and LWW picks the higher timestamp. This discards one write, causing data loss. Applications might need custom conflict resolution for semantic correctness."

## common_mistakes - Watch Out For
- Confusing concurrent writes with causally-ordered writes
- Not explaining what LWW actually discards
- Missing application-level implications of data loss

## follow_up_excellent - Depth Probe
**Question**: "How would you extend COPS to support better conflict resolution for a shopping cart use case?"
- **Looking for**: CRDTs (OR-Set), version vectors, application callbacks

## follow_up_partial - Guided Probe
**Question**: "Why can't causal consistency alone solve the problem of concurrent writes?"
- **Hint embedded**: What does 'concurrent' mean in causality terms?

## follow_up_weak - Foundation Check
**Question**: "What is last-writer-wins and when does COPS use it?"
