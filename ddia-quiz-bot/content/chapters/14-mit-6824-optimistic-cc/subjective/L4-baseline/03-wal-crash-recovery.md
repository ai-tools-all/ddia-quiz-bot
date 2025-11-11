---
id: farm-subjective-L4-003
type: subjective
level: L4
category: baseline
topic: fault-tolerance
subtopic: wal-recovery
estimated_time: 8-10 minutes
---

# question_title - Write-Ahead Logs and Crash Recovery

## main_question - Core Question
"Describe how FaRM uses write-ahead logs (WAL) stored in per-client message queues for crash recovery. Explain what information must be logged, where logs are physically stored, and walk through the recovery process when a primary crashes mid-commit. How does the per-client structure affect recovery?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- WAL entries stored in per-client queues in server non-volatile RAM
- Logs contain transaction ID, commit decision, participant list, and data values
- LOCK messages append to logs before responding
- On crash, surviving servers read logs to determine transaction state
- Per-client structure means each client's logs are distributed across servers it wrote to

### expected_keywords
- write-ahead log, per-client queue, NVRAM, crash recovery, transaction state, log replay

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- N² communication channels (each client has queue on each server)
- Log entries act as prepare records in 2PC
- Recovery coordinator can be any surviving server
- Configuration manager detects failures and triggers recovery
- Idempotent recovery operations prevent duplicate applies

### bonus_keywords
- N² channels, prepare record, recovery coordinator, configuration manager, idempotency

## sample_excellent - Example Excellence
"LOCK messages append entries to per-client WAL queues in primary's NVRAM, recording transaction ID, objects locked, and expected versions. If primary crashes after some LOCK responses but before COMMIT-PRIMARY, recovery reads logs: if all participants logged LOCK entries, complete commit; else abort. Per-client logs mean recovery must scan queues from the crashed client across all servers to reconstruct transaction state. NVRAM ensures logs survive crashes, enabling deterministic recovery."

## sample_acceptable - Minimum Acceptable
"WAL logs are in server memory in per-client queues. They record transaction information during LOCK phase. After crash, logs are read to determine if transaction should commit or abort."

## common_mistakes - Watch Out For
- Thinking logs are on disk (they're in NVRAM)
- Missing the per-client queue structure and its implications
- Not explaining what specific information must be logged
- Forgetting that logs must be written before LOCK responses

## follow_up_excellent - Depth Probe
**Question**: "Design a garbage collection strategy for WAL entries—when is it safe to delete log entries?"
- **Looking for**: After transaction completes and all replicas acknowledge, periodic checkpoints, log compaction

## follow_up_partial - Guided Probe
**Question**: "Why store logs in per-client queues rather than a single centralized log?"
- **Hint embedded**: Avoids bottleneck, parallelism, matches transaction coordinator model

## follow_up_weak - Foundation Check
**Question**: "What does 'write-ahead' mean in write-ahead logging?"
