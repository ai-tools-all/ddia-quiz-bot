---
id: ch07-isolation-level-migration-l7
day: 26
level: L7
tags: [transactions, isolation-levels, migration, principal-engineer, consistency]
related_stories: []
---

# Isolation Level Selection and Migration

## question
Your financial services platform currently uses serializable isolation for all transactions, causing performance bottlenecks (p99 latency > 5s, throughput capped at 1000 TPS). You need to selectively relax isolation levels while maintaining correctness for: (1) money transfers between accounts, (2) fraud detection scanning transactions, (3) balance inquiries, (4) transaction history generation, (5) interest calculations, and (6) regulatory reporting. Design a differentiated isolation strategy with safe migration path, considering that some operations span multiple isolation levels. Include monitoring to detect and prevent consistency violations.

## expected_concepts
- Isolation level hierarchy and anomalies
- Mixed isolation level interactions
- Consistency violation detection
- Performance vs correctness trade-offs
- Read-write split architectures
- Temporal consistency requirements
- Audit and compliance implications
- Rollback mechanisms for isolation changes

## answer
The solution requires careful analysis of operation semantics and intelligent isolation level assignment:

Isolation Level Mapping: (1) **Serializable**: Money transfers, interest calculations - require strict consistency to prevent lost updates and write skew. (2) **Repeatable Read**: Regulatory reporting - needs consistent snapshot but not real-time. (3) **Read Committed**: Fraud detection - can tolerate reading committed but slightly stale data. (4) **Read Uncommitted**: Balance inquiries for UI display (with separate path for transaction initiation). Transaction history can use snapshot isolation with async generation.

Migration Architecture: (1) **Isolation Router** - tags transactions with required isolation level based on operation type. (2) **Multi-version database setup** - separate connection pools per isolation level. (3) **Consistency barriers** - enforce serializable for operations following read-uncommitted queries that lead to writes.

Safety Mechanisms: (1) **Shadow validation** - run critical operations at both old and new isolation levels, comparing results. (2) **Anomaly detection** - monitor for patterns indicating consistency violations (negative balances, duplicate transfers). (3) **Circuit breaker per isolation level** - revert to serializable if anomalies exceed threshold. (4) **Compensating transaction framework** - automated fixes for detected inconsistencies.

Performance Optimization: (1) **Read replicas with lag monitoring** - route read-uncommitted to replicas, ensuring lag < 100ms. (2) **Optimistic locking for hot accounts** - reduce serialization bottlenecks. (3) **Batching interest calculations** - run during low-traffic periods at serializable level. (4) **Caching with versioning** - serve repeat balance inquiries from cache with version validation.

Monitoring Strategy: (1) **Consistency scorecard** - track violation rate per operation type. (2) **A/B testing framework** - gradually roll out isolation changes with control groups. (3) **Transaction genealogy** - trace how read anomalies propagate to writes. (4) **Business impact metrics** - correlate isolation changes with customer complaints or financial discrepancies.

Critical Insight: The danger isn't in the obvious cases but in transaction chains. A balance inquiry at read-uncommitted followed by a transfer decision creates a hidden dependency. Build a transaction dependency graph and enforce consistency escalation when weak reads influence strong writes.

## hook
How does Google Spanner provide external consistency while allowing different isolation levels?

## follow_up
A critical bug is discovered: your fraud detection system at read-committed isolation has been missing patterns due to reading intermediate states of multi-step fraud operations. However, moving to serializable isolation would increase false positives by 10x and add 2-second latency. Design a custom isolation solution that captures complete fraud patterns without the performance penalty, considering that fraudsters actively exploit timing windows.

## follow_up_answer
This requires a specialized temporal isolation mechanism designed for pattern detection:

Custom Isolation Design: (1) **Temporal Coherence Windows** - fraud detection reads see a consistent snapshot of all transactions within a time window (e.g., last 5 minutes), not just individually committed transactions. (2) **Causal Consistency Tracking** - maintain happens-before relationships between transactions, ensuring fraud detection sees causally related operations together. (3) **Deferred Visibility** - transactions become visible to fraud detection only after related transactions complete or timeout.

Implementation Architecture: (1) **Transaction Correlation Service** - identifies related transactions using ML-based pattern recognition (same IP, device fingerprint, payment method). (2) **Coherent Snapshot Manager** - maintains multiple consistent snapshots at different time boundaries. (3) **Stream Processing Pipeline** - Apache Flink/Beam processing transaction streams with watermarks ensuring complete windows.

Pattern Detection Enhancement: (1) **Speculative Execution Paths** - evaluate fraud rules on both committed and pending transaction combinations. (2) **Temporal Join Windows** - correlate transactions within sliding time windows even if committed at different times. (3) **Graph-based Analysis** - build transaction graphs showing money flow, detecting cycles and suspicious patterns regardless of commit order.

Performance Optimization: (1) **Incremental Materialization** - maintain running aggregates updated with each transaction. (2) **Probabilistic Filters** - Bloom filters for quick elimination of non-fraudulent patterns. (3) **Tiered Analysis** - fast path for obvious cases, detailed analysis for suspicious patterns. (4) **Read-only Replicas** - dedicated replicas with controlled lag for fraud detection.

Anti-Gaming Measures: (1) **Random Delay Injection** - add jitter to transaction visibility to prevent timing attacks. (2) **Shadow Periods** - invisible holding period where transactions are monitored before final commit. (3) **Retroactive Analysis** - re-evaluate past transactions when new patterns emerge.

Key Insight: Traditional isolation levels assume independent transactions, but fraud detection needs to see the "story" across related transactions. Build a domain-specific isolation model that understands transaction relationships and provides consistency at the pattern level, not just individual operation level. This might mean seeing "incomplete" individual transactions but "complete" behavioral patterns.
