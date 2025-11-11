---
id: farm-subjective-L3-003
type: subjective
level: L3
category: baseline
topic: replication
subtopic: commit-backup
estimated_time: 6-8 minutes
---

# question_title - COMMIT-BACKUP and Durability

## main_question - Core Question
"Explain the purpose of the COMMIT-BACKUP phase in FaRM's commit protocol. What information is sent to backups, when is it sent relative to COMMIT-PRIMARY, and how does this provide fault tolerance? Give an example of a failure scenario where backups are essential for recovery."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- COMMIT-BACKUP propagates committed writes to backup replicas
- Sent after COMMIT-PRIMARY completes at primaries
- Backups store the same data and version numbers as primaries
- F+1 replicas tolerate F failures; only one surviving replica needed
- Example: primary fails after commit, backup can serve reads and become new primary

### expected_keywords
- backup replica, fault tolerance, replication, primary failure, recovery, F+1, durability

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Non-volatile RAM ensures both primary and backup data survives power failures
- Configuration manager (ZooKeeper) tracks which servers are primary/backup
- Write-ahead logs help coordinate recovery
- Unlike Paxos, doesn't need majority—any surviving replica works

### bonus_keywords
- NVRAM, ZooKeeper, configuration manager, WAL, Paxos comparison, quorum

## sample_excellent - Example Excellence
"COMMIT-BACKUP sends committed writes and updated versions to backup replicas after primaries commit via COMMIT-PRIMARY. Backups replicate the data for fault tolerance. With F+1 replicas, the system tolerates F failures—only one replica needed to recover. Example: T1 writes X@v5→X@v6 on primary P1 and backup B1. If P1 crashes after commit, B1 has X@v6 and can serve reads or become the new primary, ensuring no data loss."

## sample_acceptable - Minimum Acceptable
"COMMIT-BACKUP replicates data to backups after the primary commits. This provides fault tolerance so if the primary fails, backups have the data."

## common_mistakes - Watch Out For
- Thinking backups must acknowledge before client sees commit
- Confusing with Paxos-style quorum requirements
- Not explaining the F+1 replica model
- Missing the role of configuration manager in tracking replicas

## follow_up_excellent - Depth Probe
**Question**: "How would you modify the protocol to tolerate both server failures and data center-wide network partitions?"
- **Looking for**: Multi-datacenter replication, consensus protocols, trade-off with latency

## follow_up_partial - Guided Probe
**Question**: "Why doesn't FaRM require backup acknowledgments before responding to the client?"
- **Hint embedded**: NVRAM persistence and primary authority

## follow_up_weak - Foundation Check
**Question**: "What does F+1 replicas mean in terms of fault tolerance?"
