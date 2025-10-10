---
id: raft-subjective-L5-002
type: subjective
level: L5
category: bar-raiser
topic: raft-consensus
subtopic: log-replication
estimated_time: 12-15 minutes
---

# question_title - Raft Log Replication and Consistency Guarantees

## main_question - Core Question
"Explain how Raft ensures that committed log entries are never lost and how it handles inconsistencies in follower logs after leader changes. Include the role of the AppendEntries RPC in maintaining log consistency."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Log Matching Property**: If entries match at an index, all preceding entries match
- **AppendEntries Consistency Check**: prevLogIndex and prevLogTerm validation
- **Commitment Rule**: Entry committed only after replicated to majority
- **Leader Completeness**: Committed entries present in all future leaders
- **Log Reconciliation**: Leader overwrites conflicting follower entries

### expected_keywords
- Primary keywords: log entry, commit, replication, consistency, AppendEntries
- Technical terms: prevLogIndex, prevLogTerm, commit index, log matching

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **No-holes Guarantee**: Logs are contiguous sequences
- **Fast Rollback Optimization**: Batch conflict resolution
- **Persistent vs Volatile State**: What must survive crashes
- **Safety vs Liveness**: Consistency over availability choice
- **Snapshot Integration**: Log compaction interaction
- **Client Interaction**: Exactly-once semantics

### bonus_keywords
- Implementation: nextIndex, matchIndex arrays, heartbeat piggybacking
- Optimizations: Pipeline, batch operations, parallel RPCs
- Edge cases: Duplicate detection, retry handling, snapshot boundaries

## sample_excellent - Example Excellence
"Raft maintains log consistency through several interlocking mechanisms:

**Log Matching Property**: The foundation - if two logs have identical entries at index i with same term, then all entries from start through i are identical. This is maintained by AppendEntries consistency check.

**AppendEntries Protocol**: Each AppendEntries RPC includes prevLogIndex and prevLogTerm - the entry immediately before new entries. Followers reject if their log doesn't match at this position. This detects inconsistencies.

**Reconciliation Process**: When AppendEntries fails, leader decrements nextIndex for that follower and retries. Eventually finds matching point, then overwrites any conflicting entries. This ensures followers' logs become identical to leader's.

**Commitment Safety**: Leader only commits an entry after replicating to majority. The commit index advances monotonically and is communicated in AppendEntries. Once committed, the entry is guaranteed to appear in all future leaders' logs because:
1. Election restriction: Candidates with incomplete logs can't get majority votes
2. Voters compare log completeness using last log index and term

**Crash Recovery**: Leaders never modify their own logs - only append. After crash, new leader's log becomes truth. Any uncommitted entries from old terms might be discarded, but committed entries are preserved.

Example: Leader has [1,1,2,3], follower has [1,1,2,4]. AppendEntries for term 3 entry will fail at position 4 (term mismatch). Leader retries at position 3, succeeds, then overwrites position 4 with correct entry."

## sample_acceptable - Minimum Acceptable
"Raft uses AppendEntries RPC to maintain log consistency. The leader includes the previous entry's index and term with each replication request. If the follower's log doesn't match at that position, it rejects the request. The leader then backs up and retries until finding a match point, then overwrites the follower's conflicting entries. Entries are only committed after being replicated to a majority of servers."

## common_mistakes - Watch Out For
- Confusing committed vs uncommitted entries
- Not understanding log reconciliation direction (leader→follower)
- Missing the election restriction role
- Ignoring persistent state requirements

## follow_up_excellent - Depth Probe
**Question**: "How does Raft handle the scenario where a leader commits entries from previous terms? Why is there a restriction on this?"
- **Looking for**: Can only commit from current term, indirect commitment of older entries
- **Red flags**: Not understanding the figure 8 problem in the Raft paper

## follow_up_partial - Guided Probe
**Question**: "You mentioned the leader backs up nextIndex on failure. Could this be optimized? How?"
- **Hint embedded**: Follower could send conflict information
- **Concept testing**: Understanding optimization opportunities

## follow_up_weak - Foundation Check
**Question**: "If a server crashes and restarts, which parts of its log must it have saved to disk?"
- **Simplification**: Persistent vs volatile state
- **Building block**: Durability requirements

## bar_raiser_question - L5→L6 Challenge
"Design a mechanism to allow Raft to support geo-distributed clusters with high-latency links between regions. How would you optimize log replication while maintaining strong consistency?"

### bar_raiser_concepts
- Witness/learner nodes
- Regional batching
- Quorum intersection strategies
- WAN optimization techniques
- Flexible quorums

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 7-9 min discussion
- **Common next topics**: Chain replication, Multi-Paxos, Viewstamped Replication, EPaxos
