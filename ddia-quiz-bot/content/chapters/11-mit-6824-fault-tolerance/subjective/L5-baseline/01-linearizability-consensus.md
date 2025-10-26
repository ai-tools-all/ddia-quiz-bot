---
id: fault-tolerance-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: fault-tolerance
subtopic: linearizability-consensus
estimated_time: 7-9 minutes
---

# question_title - Linearizability Through Consensus

## main_question - Core Question
"How does Raft provide linearizable semantics for client operations? Explain the relationship between consensus, state machine replication, and linearizability, including how Raft handles duplicate client requests."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Consensus on Order**: All servers agree on operation sequence
- **State Machine Replication**: Identical initial state + same ops = same final state
- **Client Interaction**: Read/write through leader, exactly-once semantics
- **Duplicate Detection**: Client ID and sequence numbers

### expected_keywords
- Primary keywords: linearizability, consensus, state machine, idempotence
- Technical terms: client session, sequence number, applied index

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Read Optimizations**: ReadIndex, lease reads for performance
- **Session Management**: Client registration, timeout handling
- **Exactly-once Semantics**: Deduplication at state machine layer
- **Time Bounds**: Real-time ordering guarantees
- **Split-brain Prevention**: How consensus prevents inconsistent reads

### bonus_keywords
- Implementation: client table, request cache
- Related: vector clocks, causal consistency
- Optimizations: read-only queries, follower reads

## sample_excellent - Example Excellence
"Raft achieves linearizability by using consensus to establish a total order of all operations, then applying them to replicated state machines. When a client submits an operation, it goes to the leader, who assigns it a log index. Through consensus (majority replication), all servers agree on what operation occupies each log position. Once committed, servers apply operations in log order to their state machines. Since all servers start from the same state and apply the same operations in the same order, they maintain identical state - this is the state machine replication principle. For linearizability's real-time ordering requirement, operations appear to take effect atomically at the moment they're committed. To handle failures and retries, clients include unique IDs and sequence numbers with requests. The state machine layer tracks completed operations per client, ignoring duplicates but returning cached responses. This prevents double-execution while preserving exactly-once semantics even across leader failures. For reads, the simplest approach routes them through the leader's commit process, though optimizations like ReadIndex (checking commit status without new entry) or lease-based reads can improve performance while maintaining linearizability."

## sample_acceptable - Minimum Acceptable
"Raft provides linearizability by making all operations go through the leader, which orders them in the log. Once operations are committed through consensus, all servers apply them in the same order to their state machines, ensuring identical state. Clients include IDs with requests so the system can detect and ignore duplicates from retries, providing exactly-once execution. This gives the appearance that operations happen atomically at a single point in time."

## common_mistakes - Watch Out For
- Confusing consensus with linearizability (related but different)
- Not explaining duplicate detection mechanism
- Thinking followers can serve reads by default
- Missing the state machine replication principle

## follow_up_excellent - Depth Probe
**Question**: "A client submits a write, the leader replicates it to a majority and responds, but then crashes before the entry is applied to any state machine. What happens when the client retries with the new leader?"
- **Looking for**: Entry still committed, new leader will apply it, dedup prevents double execution
- **Red flags**: Thinking the operation is lost

## follow_up_partial - Guided Probe  
**Question**: "You mentioned reads go through the leader. Why can't followers serve reads directly from their state machines?"
- **Hint embedded**: Followers might be behind
- **Concept testing**: Understanding staleness issues

## follow_up_weak - Foundation Check
**Question**: "Imagine a bank with multiple tellers (servers) processing transactions. How would you ensure everyone sees the same account balance?"
- **Simplification**: Single source of truth
- **Building block**: Ordering importance

## bar_raiser_question - L5→L6 Challenge
"Design a optimization where followers can serve linearizable reads without contacting the leader. What invariants must you maintain? What are the trade-offs?"

### bar_raiser_concepts
- Lease-based reads: Leader heartbeats as time-bounded lease
- Clock synchronization requirements and risks
- ReadIndex: Followers confirm leader's commit index
- Trade-offs: Availability vs latency vs complexity
- Safety conditions: Bounded clock drift, lease timeout > max clock skew

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 3-4 min discussion
- **Common next topics**: Jepsen testing, session guarantees, multi-datacenter consensus

## assistant_answer
Raft linearizes by establishing a single committed log order via leader-based consensus and applying operations identically on all replicas (state machine replication). Exactly-once semantics use client IDs and monotonically increasing sequence numbers to dedupe retries and return cached results. Reads can go through commit, or use ReadIndex/lease reads while preserving linearizability.

## improvement_suggestions
- Ask for the dedup table schema (clientID→lastSeq, result) and eviction/GC strategy.
- Require comparison of commit-path reads vs ReadIndex vs leader leases, including safety preconditions and trade-offs.

## improvement_exercises
### exercise_1 - Deduplication Table Design
**Question**: "Design the dedup/request cache schema to ensure exactly-once semantics across leader failover. Include eviction and recovery."

**Sample answer**: "Map clientID → {lastSeq, lastResult}. On apply, if seq==lastSeq return cached result; if seq>lastSeq, execute and update; if seq<lastSeq, ignore and return cached. Persist in the state machine or WAL-snapshots; evict using LRU with TTL per client session; ensure snapshot carries table contents for recovery."

### exercise_2 - Read Path Comparison
**Question**: "Compare leader-commit reads, ReadIndex, and lease reads: latency, safety assumptions, and failure modes."

**Sample answer**: "Commit reads: highest latency, strongest (no extra assumptions). ReadIndex: single RTT to confirm leader’s commitIndex ≥ read index; safe without clocks, requires leader availability. Lease reads: lowest latency but need bounded clock skew and lease validity; unsafe if leader’s lease time isn’t truly exclusive due to skew or pause-the-world events."
