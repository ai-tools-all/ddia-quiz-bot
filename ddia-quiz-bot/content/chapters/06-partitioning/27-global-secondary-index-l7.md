---
id: ch06-global-secondary-index-l7
day: 27
level: L7
tags: [partitioning, secondary-indexes, distributed-systems, principal-engineer, consistency]
related_stories: []
---

# Global Secondary Index Architecture

## question
You're architecting a distributed document store (100B documents, 50TB) requiring complex secondary indexes for an e-commerce search system. Requirements: (1) Support 20+ index fields with frequent schema evolution, (2) Consistent index reads within 1 second of writes, (3) Index updates shouldn't block document writes, (4) Support both exact match and range queries. As Principal Engineer, design a global secondary index system that balances consistency, performance, and operational maintainability. Consider index corruption recovery and partial index rebuilds.

## expected_concepts
- Synchronous vs asynchronous index updates
- Index partitioning strategies independent of primary data
- Consistency models for secondary indexes
- Index repair and anti-entropy mechanisms
- Schema evolution with zero-downtime
- Distributed transaction coordination
- Write amplification and its mitigation

## answer
The solution requires a decoupled index architecture with tunable consistency:

Architecture: (1) **Separate index clusters** from primary storage, allowing independent scaling and isolation from primary write path. (2) **Hybrid update model** - critical indexes (inventory, pricing) use synchronous updates via distributed saga pattern; non-critical indexes (search, analytics) use async with bounded staleness.

Index partitioning: (1) **Term-based partitioning** for index entries, but partition by hash of term for even distribution. (2) **Composite indexes** stored redundantly with different partition keys for various access patterns. (3) **Covering indexes** for frequent query patterns, trading storage for query performance.

Consistency mechanism: (1) **Write-ahead log (WAL) shipping** from primary to index clusters with guaranteed delivery. (2) **Logical timestamps** (hybrid logical clocks) for ordering updates across distributed system. (3) **Read repair** on staleness detection - client can force refresh if timestamp indicates staleness. (4) **Periodic anti-entropy** sweeps comparing checksums between primary and index data.

Schema evolution: (1) **Online index building** using snapshot + catch-up pattern. (2) **Progressive rollout** - new indexes start as "building", transition to "active" after validation. (3) **Index versioning** - multiple versions coexist during migration with query router handling version selection.

Operational excellence: (1) **Index health scoring** based on latency, staleness, and corruption rate. (2) **Automatic quarantine** of corrupted index partitions with fallback to scan. (3) **Incremental repair** using merkle trees to identify divergent ranges. (4) **Index advisor** analyzing query patterns to suggest new indexes or removal of unused ones.

Key innovation: **Adaptive consistency** - System automatically adjusts sync/async threshold based on load. During peak, more indexes become async to protect write latency.

## hook
How do you handle index updates when a single document change affects 20+ indexes?

## follow_up
Your global secondary index system has been running successfully for a year. Now, the business wants to implement GDPR-compliant "right to be forgotten" where user data must be completely purged within 72 hours, including from all indexes, backups, and replicas. However, your indexes use immutable LSM trees with compaction, and some indexes are derived (materialized views of joins). How do you architect a deletion pipeline that guarantees complete removal while maintaining system performance?

## follow_up_answer
This requires a comprehensive deletion orchestration system that tracks data lineage:

Deletion architecture: (1) **Tombstone propagation system** - Special tombstone records that propagate through all data flows with higher priority than normal updates. Include deletion deadline timestamp for SLA tracking. (2) **Lineage tracking** - Maintain directed graph of data dependencies: primary → indexes → derived indexes → caches → backups. Each node tracks TTL for deletion completion.

Index-specific handling: (1) **Forced compaction triggers** - Upon deletion request, schedule priority compaction for affected SSTable ranges to physically remove data. (2) **Deletion generations** - Track which generation of LSM tree contains user data, verify deletion after major compaction. (3) **Soft-delete with encryption** - Immediately encrypt deleted records with per-user key, destroy key for instant logical deletion while awaiting physical removal.

Derived data challenges: (1) **Reverse index mapping** - Maintain inverse indexes to track which derived views contain specific user data. (2) **Cascading regeneration** - Mark all derived data as "tainted" when source is deleted, trigger regeneration without deleted records. (3) **Join materialization tracking** - Store source record IDs in materialized views to enable targeted deletion.

Verification pipeline: (1) **Deletion certificates** - Cryptographic proof of deletion from each system component. (2) **Scanning verification** - Background job samples data locations to verify absence of deleted data. (3) **Audit logging** - Immutable audit trail of deletion requests and confirmations for compliance proof.

Performance optimization: (1) **Batched deletions** - Aggregate deletion requests in 1-hour windows for efficient processing. (2) **Deletion-aware caching** - Bloom filters for deleted IDs to avoid unnecessary searches. (3) **Resource reservation** - Dedicated capacity for deletion pipeline to meet SLAs without impacting normal operations.

Critical insight: GDPR compliance in distributed systems isn't just about deletion - it's about proving deletion. Design systems with deletion as a first-class operation, not an afterthought. The cost of retrofitting deletion capabilities far exceeds building them in initially.
