---
id: occ-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: occ
subtopic: 2pc-coordinator-crash
estimated_time: 8-10 minutes
---

# question_title - 2PC with OCC: Coordinator Crash

## main_question - Core Question
"Describe the two-phase commit flow with OCC (lock/validate, commit primary, commit backup). If the coordinator crashes after some primaries are locked, how is state recovered and what outcomes are possible?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- LOCK phase on primaries validates versions and sets lock bits
- COMMIT to primaries applies writes; COMMIT to backups propagates
- Recovery reads per-client logs / server memory to determine state
- Idempotent commit/abort; unlock/rollback paths

### expected_keywords
- lock bit, version, write-ahead/log, primary, backup, idempotent

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Timeouts and re-issuing messages
- Ensuring locks aren’t left stuck (lock cleanup)
- Difference from majority quorum systems

### bonus_keywords
- recovery coordinator, in-doubt, cleanup protocol

## sample_excellent - Example Excellence
"Coordinator locks each primary with version checks. If it crashes after locking some shards, recovery inspects logs to see whether commits were sent. If no commit records exist, abort and clear locks; if commit was sent to some primaries, re-send commit to the rest (idempotent). Backups are updated after primaries; recovery replays pending backup updates."

## sample_acceptable - Minimum Acceptable
"Lock/validate, then commit primaries, then backups. On crash, read logs and either abort (clearing locks) or finish commits idempotently."

## common_mistakes - Watch Out For
- Assuming majority quorum is needed
- Leaving locks uncleared on abort paths
- Non-idempotent commit messages

## follow_up_excellent - Depth Probe
**Question**: "How do you guarantee locks are eventually cleared even if the coordinator is lost?"
- **Looking for**: lease/timeout cleanup, helper coordinators

## follow_up_partial - Guided Probe  
**Question**: "Why commit primaries before backups?"
- **Hint embedded**: Define durability window vs availability

## follow_up_weak - Foundation Check
**Question**: "What are the two phases and their purpose?"
