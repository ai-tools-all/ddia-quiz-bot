---
id: occ-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: occ
subtopic: basics
estimated_time: 6-8 minutes
---

# question_title - OCC Basics and Aborts

## main_question - Core Question
"Explain the optimistic concurrency control workflow in the RDMA-based system: execution, validation, locking, and commit. Give a concrete example of two transactions that conflict and show which one aborts and why."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Execute reads without locks; record versions
- At commit: lock write-set on primaries and validate versions
- Abort on version change or observed lock; else apply writes and bump versions
- Serializability ensured by validation and lock bits

### expected_keywords
- read set, write set, version, lock bit, validate, abort, commit

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Read-only path uses header recheck only (no locks)
- Impact of contention and retry strategies
- RDMA one-sided reads during execution

### bonus_keywords
- CAS, one-sided RDMA, backoff, starvation

## sample_excellent - Example Excellence
"T1 reads X@v1, Y@v2; T2 reads X@v1 then plans to write X. T2 commits first: locks X, applies X@v2, clears lock. T1’s commit validates X’s version and detects change (v1→v2), so T1 aborts. Reads are optimistic and serial order is enforced at commit via validation/locks."

## sample_acceptable - Minimum Acceptable
"Read without locks, then lock/validate during commit. If versions changed, abort; otherwise apply writes."

## common_mistakes - Watch Out For
- Assuming reads hold locks
- Ignoring version recording at read time
- Forgetting to clear lock after commit

## follow_up_excellent - Depth Probe
**Question**: "How would you reduce aborts on a hot key?"
- **Looking for**: batching, combining, short transactions, partitioning, limited locking

## follow_up_partial - Guided Probe  
**Question**: "What exactly must the client record during initial reads to validate later?"
- **Hint embedded**: Version and lock state

## follow_up_weak - Foundation Check
**Question**: "What does ‘optimistic’ mean here?"
