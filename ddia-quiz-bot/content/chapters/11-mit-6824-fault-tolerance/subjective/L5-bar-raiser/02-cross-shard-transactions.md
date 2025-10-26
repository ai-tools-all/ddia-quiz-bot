---
id: fault-tolerance-subjective-L5-005
type: subjective
level: L5
category: bar-raiser
topic: fault-tolerance
subtopic: cross-shard-transactions
estimated_time: 9-12 minutes
---

# question_title - Cross‑Shard Transactions over Multiple Raft Groups

## main_question - Core Question
"Design an atomic, linearizable cross‑shard transaction protocol across multiple Raft groups (shards). Explain coordinator/participant roles, how prepare/commit integrate with Raft logs, and how you recover from coordinator or participant failures without violating correctness."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **2PC over Raft**: Use a transaction coordinator with two‑phase commit; each shard is a Raft group
- **Durable Prepare**: Participants log a prepare record (via Raft) before voting YES
- **Atomic Commit**: Coordinator logs and replicates COMMIT/ABORT; participants apply accordingly
- **Linearizability**: Per‑shard linearizability + atomic decision → globally linearizable outcome
- **Failure Handling**: Recovery via durable logs; participants block after YES until decision

### expected_keywords
- Primary: 2PC, coordinator, participants, atomicity, linearizability
- Technical: prepare, commit index, transaction ID (XID), write‑ahead log, idempotence

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Deadlock Avoidance**: Deterministic lock ordering or timeouts + retry
- **Exactly‑Once**: Client IDs and sequence numbers across shards
- **Coordinator Replication**: Coordinator state replicated (Raft/Paxos) to avoid SPOF
- **Sagas**: Compensation for non‑atomic workflows
- **Commit‑Wait**: Reads wait until decision visible (fencing)

### bonus_keywords
- Implementation: per‑shard lock table, XID→state table, idempotent apply
- Related: Spanner 2PC over Paxos, Percolator, Calvin

## sample_excellent - Example Excellence
"Each shard is a Raft group. The client drives a coordinator (itself replicated) that assigns XID and readies write sets. Phase 1: the coordinator sends PREPARE(XID) to each shard; each shard proposes a Raft log entry that validates/locks the keys, fsyncs via Raft commit, and replies YES with its prepare LSN, or NO. If any NO/timeouts, coordinator logs ABORT(XID) (replicated) and broadcasts ABORT. If all YES, coordinator logs COMMIT(XID) in its Raft and then sends COMMIT to participants, which append COMMIT entries and apply under idempotence. Participants that voted YES block reads/writes on locked keys until decision (commit‑wait) to preserve linearizability. Recovery: participants consult their XID table; if PREPARED with no decision and coordinator unreachable, they wait and/or run termination protocol by querying a majority of coordinator replicas. Exactly‑once uses clientID+seq to dedupe per shard. Deterministic lock ordering prevents distributed deadlocks; otherwise use timeouts and backoff."

## sample_acceptable - Minimum Acceptable
"Use two‑phase commit with a coordinator and shard participants. Participants must persist PREPARE before voting YES. The coordinator replicates COMMIT/ABORT and participants apply the final decision. Recovery checks logs to decide."

## common_mistakes - Watch Out For
- Voting YES before durable PREPARE
- Coordinator as SPOF (not replicated)
- Unlocking keys before commit visible (breaking linearizability)
- Missing idempotence/deduplication on retries

## follow_up_excellent - Depth Probe
**Question**: "Coordinator crashes after all participants are PREPARED but before broadcasting COMMIT. Walk through safe recovery without violating atomicity."
- **Looking for**: Coordinator’s replicated log decides; participants in PREPARED wait; new coordinator reads XID→COMMIT and completes; or ABORT if no durable decision

## follow_up_partial - Guided Probe  
**Question**: "How do you avoid or handle deadlocks across shards?"
- **Hint**: Deterministic key ordering or wait‑die + timeouts; abort and retry

## follow_up_weak - Foundation Check
**Question**: "Why must a participant write PREPARE to its Raft log before replying YES?"
- **Simplification**: Ensures the vote survives crashes and can be recovered

## bar_raiser_question - L5→L6 Challenge
"Design commit‑wait semantics that guarantee a post‑commit read on any shard reflects the transaction outcome. Specify the fences or version checks required and analyze the latency impact."

### bar_raiser_concepts
- Require reads to observe commit index ≥ decision LSN
- Fencing tokens or version vectors to prevent stale leaders
- Latency trade‑offs vs freshness

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: TrueTime/commit‑wait, Spanner vs Percolator, transactional KV design

## assistant_answer
Use 2PC with a replicated coordinator and Raft‑backed participants: PREPARE must be durably committed per shard before YES; coordinator replicates COMMIT/ABORT and participants apply idempotently. Commit‑wait ensures linearizable reads by waiting for a decision fence; recovery consults replicated logs to terminate in‑doubt transactions safely.

## improvement_suggestions
- Require an XID→state table schema and exact recovery algorithm (termination protocol) description.
- Ask for lock ordering strategy and fairness/backoff policy to control contention.

## improvement_exercises
### exercise_1 - XID Table and Recovery
**Question**: "Define the XID table fields and write a termination protocol when a participant is PREPARED but cannot reach the coordinator."

**Sample answer**: "XID→{state, participants[], decisionLSN, startTS}. If PREPARED and coordinator unreachable, query a quorum of coordinator replicas: if COMMIT/ABORT found, apply; else wait/backoff or elect new coordinator to decide ABORT after timeout."

### exercise_2 - Commit‑Wait Fence
**Question**: "Specify the precise fence needed for linearizable post‑commit reads on participants."

**Sample answer**: "Before serving a read on keys touched by XID, ensure local commitIndex ≥ decisionLSN for COMMIT/ABORT entry and that leader lease/fence token is valid; otherwise issue ReadIndex/leader confirmation or wait."
