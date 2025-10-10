---
id: ch07-optimistic-pessimistic-hybrid-l7
day: 27
level: L7
tags: [transactions, concurrency-control, architecture, principal-engineer, performance]
related_stories: []
---

# Hybrid Concurrency Control System

## question
You're designing a collaborative document editing system (think Google Docs meets Git) for enterprise with strict requirements: (1) real-time collaborative editing with sub-100ms latency, (2) atomic multi-document transactions for refactoring operations, (3) branching and merging with conflict resolution, (4) compliance requiring immutable audit trails, and (5) offline editing with eventual sync. Design a hybrid optimistic/pessimistic concurrency control system that adapts based on operation type, conflict probability, and network conditions. Consider how different concurrency mechanisms interact when users perform different operation types simultaneously.

## expected_concepts
- Operational Transformation (OT) vs CRDTs
- Multi-granularity locking
- Adaptive concurrency control
- Intention locks and lock escalation
- Optimistic read with pessimistic write
- Vector clocks and version vectors
- Three-way merge strategies
- Consistency levels per operation type

## answer
The solution requires a multi-layered concurrency control system with intelligent adaptation:

Concurrency Control Layers: (1) **Character-level CRDT** for real-time typing - lock-free, automatic conflict resolution. (2) **Block-level optimistic locking** for paragraph edits - compare-and-swap with retry. (3) **Document-level pessimistic locking** for structural changes - exclusive locks for schema modifications. (4) **Transaction-level 2PL** for multi-document operations - ensures atomicity across documents.

Adaptive Strategy: (1) **Conflict prediction model** - ML-based prediction of conflict probability based on edit patterns, user behavior, and document hotspots. (2) **Dynamic escalation** - start optimistic, escalate to pessimistic after N conflicts. (3) **Network-aware modes** - optimistic for good connectivity, pessimistic for flaky networks to avoid repeated retries. (4) **Time-based leases** - pessimistic locks auto-degrade to optimistic after idle period.

Real-time Collaboration: (1) **Operational Transformation pipeline** - transform concurrent ops to maintain consistency. (2) **Presence awareness** - show user cursors and selections to reduce conflicts. (3) **Micro-locking** - temporary exclusive access to text ranges during active typing. (4) **Differential synchronization** - continuous background sync with automatic merge.

Multi-document Transactions: (1) **Hierarchical intention locks** - IS/IX/S/X locks at document collection level. (2) **Deadlock prevention** - consistent lock ordering based on document IDs. (3) **Savepoints** - partial rollback capability within transactions. (4) **Read stability** - ensure consistent reads across documents during transaction.

Offline Support: (1) **Three-way merge** - local changes, remote changes, common ancestor. (2) **Conflict-free replicated operations log** - record operations not state. (3) **Semantic conflict detection** - understand document structure for intelligent merging. (4) **Manual resolution queue** - complex conflicts requiring user intervention.

Audit Compliance: (1) **Event sourcing backbone** - immutable log of all operations. (2) **Merkle tree verification** - cryptographic proof of audit trail integrity. (3) **Point-in-time reconstruction** - rebuild document state at any timestamp.

Critical Insight: The key is recognizing that different parts of the document have different consistency requirements. Text editing can be eventually consistent with CRDTs, while table formulas need stricter consistency. Build a pluggable concurrency control framework where each document element type can specify its consistency requirements and preferred concurrency mechanism.

## hook
How does Google Docs handle the "everyone editing the meeting notes simultaneously" scenario?

## follow_up
During a major product launch, your system experiences a "thundering herd" problem: 10,000 users simultaneously open and edit the launch announcement document. The document contains rich media, embedded spreadsheets with formulas, and real-time comments. Your current architecture is failing with lock timeouts and inconsistent states. Design an emergency scaling strategy that can be deployed without downtime, including sharding the document, managing cross-shard consistency, and ensuring the CEO's edits always succeed.

## follow_up_answer
This requires rapid deployment of document sharding with intelligent coordination:

Emergency Sharding Strategy: (1) **Content-based partitioning** - shard document by sections (paragraphs, tables, media) with independent version control. (2) **Read-write splitting** - multiple read replicas with single write master per shard. (3) **Hierarchical caching** - edge cache for reads, regional cache for writes, global master for conflicts.

Load Distribution: (1) **User affinity routing** - assign users to specific replicas based on hash(user_id + section_id). (2) **Dynamic shard splitting** - automatically split hot sections when edit rate exceeds threshold. (3) **Read-only mode toggle** - gracefully degrade majority to read-only during peak load. (4) **Priority lanes** - dedicated resources for VIP users (CEO, launch team).

Cross-shard Consistency: (1) **Eventual consistency for text** - allow temporary divergence, merge async. (2) **Synchronous coordination for formulas** - distributed locks for spreadsheet calculations. (3) **Causal consistency for comments** - ensure comment threads maintain order. (4) **Global sequence numbers** - total ordering for audit trail despite sharding.

Immediate Deployment: (1) **Proxy layer injection** - add sharding proxy without app changes. (2) **Shadow sharding** - run sharded system in parallel, switch on validation. (3) **Circuit breaker per shard** - isolate failures to prevent cascade. (4) **Gradual migration** - move sections to shards based on heat map.

CEO Priority Lane: (1) **Reserved connection pool** - dedicated database connections for VIP users. (2) **Optimistic lock stealing** - CEO edits can preempt others with notification. (3) **Fast-path processing** - skip non-critical validations for VIP edits. (4) **Dedicated replica** - CEO reads from consistent, low-latency replica.

Formula Recalculation: (1) **Dependency graph caching** - precompute formula dependencies. (2) **Incremental calculation** - only recalc affected cells. (3) **Async calculation queue** - defer complex recalcs, show "calculating..." state. (4) **Calculation sharding** - distribute formula evaluation across workers.

Key Insight: In crisis, perfect consistency is the enemy of availability. Implement "consistency SLAs" - different document parts get different guarantees. The CEO's message must be consistent, while comments can be eventually consistent. Build an adaptive system that can dynamically trade consistency for availability based on load, with clear degradation paths and user communication about current system state.
