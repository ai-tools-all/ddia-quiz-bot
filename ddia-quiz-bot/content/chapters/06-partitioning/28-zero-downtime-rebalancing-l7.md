---
id: ch06-zero-downtime-rebalancing-l7
day: 28
level: L7
tags: [partitioning, rebalancing, operations, principal-engineer, migration]
related_stories: []
---

# Zero-Downtime Dynamic Rebalancing

## question
Your team manages a critical financial ledger system (500TB, 100K TPS) using range-based partitioning. The business is acquiring another company, which will double the data volume and change access patterns significantly. As Principal Engineer, design a rebalancing strategy that: (1) Migrates from range to hash partitioning, (2) Maintains strict consistency and zero data loss, (3) Achieves zero downtime with <10ms p99 latency impact, (4) Can be rolled back at any point. Consider that some transactions span multiple partitions and the system processes $10M/minute.

## expected_concepts
- Online schema/partition migration patterns
- Dual writes and incremental cutover
- Distributed transaction handling during migration
- State machine for migration orchestration
- Rollback strategies and checkpointing
- Performance isolation techniques
- Split-brain prevention during transition

## answer
The solution requires a sophisticated state-machine-driven migration with multiple safety checks:

Migration architecture: (1) **Parallel universe approach** - Build complete hash-partitioned system alongside existing range-partitioned system. (2) **Logical replication layer** - Abstract physical partitioning from application using logical partition mapping that can route to either system. (3) **Phased migration state machine** - States: STABLE → BUILDING → SHADOWING → VALIDATING → SWITCHING → FINALIZING → COMPLETE, with rollback paths from each state.

Consistency preservation: (1) **Dual write protocol** - All writes go to both systems with 2PC coordination ensuring atomicity. (2) **Read router** with configurable strategy: start with 100% range, gradually shift to hash after validation. (3) **Transaction bridging** - Distributed transaction coordinator handles transactions spanning old/new partitions using saga pattern with compensating transactions.

Zero-downtime techniques: (1) **Incremental partition migration** - Migrate 1% of keyspace at a time with continuous validation. (2) **Request hedging** - Send reads to both systems, use faster response, detect inconsistencies. (3) **Smart client routing** - Clients maintain dual routing tables, handle partition ownership changes without server roundtrips. (4) **Backpressure coordination** - Both systems share congestion signals to prevent overload during migration.

Validation and safety: (1) **Continuous consistency checker** - Sample reads from both systems, alert on mismatches. (2) **Shadow traffic analysis** - Compare performance metrics between systems before switching. (3) **Automated rollback triggers** - Latency spike, error rate, or consistency violations trigger automatic rollback. (4) **Financial reconciliation** - Periodic sum checks ensure no money lost during migration.

Rollback strategy: (1) **Checkpoint-based recovery** - Save state after each successful partition migration. (2) **Reverse replication** - Maintain change stream from new to old system for instant rollback. (3) **Traffic drain protocol** - Gracefully drain in-flight requests before rollback. 

Critical insight: The real risk isn't in the technology but in the coordination. A migration of this scale requires treating the migration system itself as a distributed system with its own failure modes, consistency requirements, and operational playbooks.

## hook
How do you handle a network partition between old and new systems during the migration?

## follow_up
Three weeks into your successful migration (now at 60% completion), a critical bug is discovered: the hash function has a subtle bias causing 15% of keys to hash to the same partition, creating severe hot spots. The bug is in production with 300TB already migrated. You cannot simply fix the hash function as it would require re-migrating all data. How do you complete the migration while fixing this issue, knowing that rolling back means losing three weeks of complex migration progress?

## follow_up_answer
This requires a surgical intervention that fixes the problem while preserving progress:

Immediate mitigation: (1) **Virtual partition subdivision** - Split hot partitions into virtual sub-partitions, distribute these across nodes to spread load. This is a routing-layer fix, not data movement. (2) **Admission control** - Rate-limit writes to hot partitions, queue in distributed buffer with fair scheduling. (3) **Read replicas** - Add dedicated read replicas for hot partitions to handle query load.

Permanent solution: (1) **Three-phase migration evolution** - Current phase continues with biased hash for already-migrated data. New phase uses corrected hash function for unmigrated data. Final phase gradually re-redistributes biased data. (2) **Hybrid routing table** - Maintain mapping of key ranges to hash functions used, route appropriately. (3) **Background rebalancer** - Slowly move data from biased partitions to correct locations, 1% per day to avoid impact.

Implementation strategy: (1) **Hash function versioning** - Tag each partition with hash version used for its data. (2) **Lazy migration** - When reading biased keys, opportunistically migrate to correct partition if system load permits. (3) **Write-time correction** - New writes to biased keys go to correct partition, maintain pointer from old location. (4) **Compaction-triggered migration** - During regular compaction, move biased data to correct partitions.

Risk management: (1) **Partial rollback option** - Can still rollback unmigrated portion (40%) if new approach fails. (2) **A/B testing** - Test corrected hash on small subset before full deployment. (3) **Financial impact analysis** - Calculate cost of completing biased migration vs restart, present to business. (4) **Hot spot prediction model** - ML model predicts future hot spots based on access patterns, proactively mitigates.

Long-term learning: (1) **Pre-production validation** - Should have run statistical analysis on hash distribution with production data sample. (2) **Canary migrations** - Start with 0.1% of data to detect issues early. (3) **Reversible decisions** - Design migrations where early decisions can be revised without full restart.

Key insight: In production systems, perfect solutions are often impossible. The art lies in finding creative ways to fix problems while the plane is flying, accepting temporary complexity for long-term correctness.
