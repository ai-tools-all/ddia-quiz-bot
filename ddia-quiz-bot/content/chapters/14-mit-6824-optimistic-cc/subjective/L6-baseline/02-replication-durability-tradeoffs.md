---
id: occ-subjective-L6-002
type: subjective
level: L6
category: baseline
topic: occ
subtopic: replication-durability
estimated_time: 12-15 minutes
---

# question_title - Replication and Durability Trade-offs

## main_question - Core Question
"Discuss durability and availability in the single–data center primary-backup design with non-volatile RAM. What acknowledgments are needed before returning success? Analyze failure windows (primary crash, backup crash, power loss) and how logs enable recovery."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Primary applies write and updates backup(s) for F+1 replication
- Define when to ack client vs when backups confirm
- Non-volatile RAM and per-client logs aid crash recovery
- Availability vs durability trade-off in ack policy

### expected_keywords
- primary, backup, ack policy, recovery log, power failure

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Risk of acknowledging before backup persistence
- Synchronous vs async backup updates and latency impact
- Recovery ordering and idempotent re-apply

### bonus_keywords
- window of vulnerability, fsync analogue, idempotency

## sample_excellent - Example Excellence
"Ack after primaries commit and at least one backup durably stores log to bound data loss. With NVRAM, power failure preserves state, but node crashes rely on logs in server memory. If you ack on primary-only, a primary crash before backup update risks loss; synchronous backup update removes that window at latency cost. Recovery replays per-client logs idempotently."

## sample_acceptable - Minimum Acceptable
"To reduce loss risk, wait for a backup to store the update before ack. Logs let you recover after crashes."

## common_mistakes - Watch Out For
- Assuming replication across DCs
- Ignoring ack policy’s effect on durability window
- Non-idempotent replays

## follow_up_excellent - Depth Probe
**Question**: "How would you change policy under extremely tight latency SLOs?"
- **Looking for**: primary-only ack trade-off, risk acceptance, telemetry

## follow_up_partial - Guided Probe  
**Question**: "What is the ‘window of vulnerability’?"
- **Hint embedded**: Between primary commit and backup persistence

## follow_up_weak - Foundation Check
**Question**: "Why keep per-client logs?"
