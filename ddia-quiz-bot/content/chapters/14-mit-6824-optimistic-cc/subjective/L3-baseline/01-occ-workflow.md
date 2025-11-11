---
id: farm-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: occ
subtopic: workflow
estimated_time: 6-8 minutes
---

# question_title - OCC Workflow in FaRM

## main_question - Core Question
"Explain the optimistic concurrency control workflow in FaRM: execution phase, validation phase, and commit phase. Give a concrete example of two transactions that conflict on the same object and explain which one commits and which one aborts, and why."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Execute phase: read objects without locks using one-sided RDMA; record versions
- Validation phase: send LOCK messages to primaries; verify versions match expected values
- Commit phase: if all locks acquired, apply writes, increment versions, clear locks; else abort
- Conflicts detected through version mismatches or lock bits already set

### expected_keywords
- read set, write set, version number, lock bit, LOCK message, abort, commit, serializability

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- COMMIT-PRIMARY and COMMIT-BACKUP messages
- One-sided RDMA reads during execution
- Retry mechanisms after aborts
- Transaction coordinator on client side

### bonus_keywords
- two-phase commit, primary-backup, RDMA, coordinator

## sample_excellent - Example Excellence
"T1 reads X@v5, Y@v3 via RDMA, plans to write X. T2 reads X@v5, plans to write X. T2 commits first: sends LOCK to X's primary, verifies v5, sets lock, sends COMMIT-PRIMARY applying X@v6, clears lock. When T1 commits, its LOCK finds X@v6 (not expected v5), so T1 aborts. OCC allows optimistic reads but enforces serial order at commit via version validation."

## sample_acceptable - Minimum Acceptable
"Read objects and record versions. At commit, check versions haven't changed. If changed, abort; otherwise apply writes and update versions."

## common_mistakes - Watch Out For
- Assuming reads acquire locks during execution
- Not explaining version recording at read time
- Forgetting lock clearing after commit
- Confusing execution phase with commit phase

## follow_up_excellent - Depth Probe
**Question**: "How would you modify the protocol to reduce aborts on a heavily contended object?"
- **Looking for**: batching writes, combining operations, partitioning data, early abort detection, backoff strategies

## follow_up_partial - Guided Probe
**Question**: "What specific information must the client record during each read to enable later validation?"
- **Hint embedded**: Version number and object location

## follow_up_weak - Foundation Check
**Question**: "What does 'optimistic' mean in optimistic concurrency control?"
