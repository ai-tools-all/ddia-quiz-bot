---
id: farm-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: fault-tolerance
subtopic: coordinator-crash
estimated_time: 8-10 minutes
---

# question_title - Coordinator Crash During Commit

## main_question - Core Question
"Explain what happens if the transaction coordinator crashes at different points during FaRM's two-phase commit protocol: (1) before sending LOCK messages, (2) after receiving all LOCK responses but before COMMIT-PRIMARY, and (3) after sending some but not all COMMIT-PRIMARY messages. How does FaRM ensure correctness in each case?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Case 1 (before LOCK): Transaction hasn't started commit, safe to abort, no state written
- Case 2 (after LOCK, before COMMIT): Locks held at primaries, need timeout/recovery to release locks
- Case 3 (during COMMIT): Partial commit state, need recovery protocol to complete or abort consistently
- Write-ahead logs at primaries enable recovery of commit decisions

### expected_keywords
- coordinator crash, two-phase commit, locks, timeout, recovery, WAL, partial commit, abort

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Coordinator recovery can read logs to determine commit state
- Blocking problem of 2PC when coordinator crashes with locks held
- Per-client message queues in server memory aid recovery
- ZooKeeper/configuration manager role in tracking failures

### bonus_keywords
- blocking, recovery protocol, configuration manager, message queue, distributed commit

## sample_excellent - Example Excellence
"(1) Before LOCK: no commit state exists, transaction implicitly aborts, primaries unaware. (2) After LOCK responses: primaries hold locks and WAL entries; coordinator recovery (or timeout) must either complete commit or release locks. (3) During COMMIT-PRIMARY: some primaries committed, some locked; recovery reads logs to determine commit decision and ensures all participants reach the same outcome, preventing split-brain."

## sample_acceptable - Minimum Acceptable
"Before commit starts, the transaction can safely abort. After locks are acquired, the system needs recovery to either complete or abort consistently. Logs help determine what happened."

## common_mistakes - Watch Out For
- Not distinguishing between different crash points
- Missing the blocking problem when locks are held
- Ignoring the role of write-ahead logs in recovery
- Assuming automatic rollback is always possible

## follow_up_excellent - Depth Probe
**Question**: "How could FaRM modify the protocol to avoid blocking when the coordinator crashes with locks held?"
- **Looking for**: Paxos-based commit, coordinator replication, lock timeouts with retry logic

## follow_up_partial - Guided Probe
**Question**: "What information must be in the write-ahead log to enable recovery?"
- **Hint embedded**: Transaction ID, commit decision, participants, data values

## follow_up_weak - Foundation Check
**Question**: "Why is it problematic for some primaries to commit while others abort the same transaction?"
