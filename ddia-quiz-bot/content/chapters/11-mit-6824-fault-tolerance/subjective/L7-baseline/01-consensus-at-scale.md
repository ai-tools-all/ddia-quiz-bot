---
id: fault-tolerance-subjective-L7-001
type: subjective
level: L7
category: baseline
topic: fault-tolerance
subtopic: consensus-at-scale
estimated_time: 10-12 minutes
---

# question_title - Consensus at Global Scale

## main_question - Core Question
"Design a consensus system for a globally distributed database spanning 20+ regions with 100ms+ latencies between continents. Requirements: 100k writes/second globally, 99.99% availability, strict serializability for financial transactions, and 5-second recovery time objective (RTO) for regional failures. Explain your architecture, trade-offs, and how you'd evolve it over 5 years."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Hierarchical Consensus**: Regional Raft/Paxos clusters with global coordination
- **Geographic Sharding**: Data locality and sovereignty considerations
- **Multi-Version Concurrency**: MVCC for read scaling and consistency
- **Adaptive Protocols**: Dynamic leader placement, flexible consistency

### expected_keywords
- Primary keywords: geo-replication, sharding, consistency, latency
- Technical terms: Spanner, Calvin, MDCC, clock synchronization

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Atomic Clocks/TrueTime**: Bounded uncertainty for global ordering
- **Deterministic Execution**: Calvin-style scheduling for cross-region transactions
- **Witness Replicas**: Voting-only replicas in expensive regions
- **Edge Caching**: CDN integration for read-heavy workloads
- **Chaos Engineering**: Systematic failure injection
- **Economic Optimization**: Cross-region bandwidth costs
- **Regulatory Compliance**: Data residency, GDPR, sovereignty

### bonus_keywords
- Systems: Spanner, CockroachDB, YugabyteDB, FaunaDB
- Techniques: loosely synchronized clocks, hybrid logical clocks
- Evolution: quantum networks, edge computing trends

## sample_excellent - Example Excellence
"I'd architect a hierarchical consensus system inspired by Google Spanner but adapted for cloud deployment:

**Architecture Layers**:

1. **Regional Layer**: Each region runs 5-node Raft clusters for local shards. Leader placement optimized for write locality. 20 regions × 5 nodes = 100 nodes minimum.

2. **Global Coordination**: Lightweight global Paxos group (7 nodes across stable regions) assigns timestamp ranges and coordinates cross-region transactions. Not on critical path for regional writes.

3. **Time Synchronization**: Without atomic clocks, use hybrid logical clocks (HLC) with NTP discipline. Accept 10-50ms uncertainty vs Spanner's <10ms. Design assumes bounded clock drift.

**Data Model**:
- **Sharding**: Range-based sharding with locality-aware placement. Financial accounts sharded by region of origin.
- **Replication**: 3-way replication within region, asynchronous replication to 2 remote regions for disaster recovery
- **MVCC**: Multi-version storage with garbage collection after global snapshot time advances

**Transaction Processing**:
- **Single-Region**: Fast path through regional Raft (5-10ms commit)
- **Cross-Region**: Two approaches:
  - Spanner-style: 2PC over Paxos with synchronized clocks (100-200ms)
  - Calvin-style: Deterministic scheduling with pre-declared read/write sets (150ms but more predictable)

**Achieving 100k writes/second**:
- 1000 shards globally, each handling 100 writes/second
- Regional batching and pipelining 
- Connection multiplexing and kernel bypass networking

**99.99% Availability** (52 minutes downtime/year):
- No single points of failure
- Automatic leader migration within 2 seconds
- Pre-positioned witness replicas for quorum during failures
- Graceful degradation: accept regional consistency during partition

**5-Second RTO**:
- Hot standby replicas with real-time log shipping
- Pre-computed failure detection with BFD
- DNS and client-side failover in parallel
- Automated runbook execution

**Trade-offs**:
- Complexity vs single-region systems
- Write latency (100ms+) for global consistency  
- Cost: 3-5x infrastructure for redundancy
- Operational burden: requires SRE team

**5-Year Evolution**:

Years 1-2: Stabilize core, improve observability
- Comprehensive distributed tracing
- Automated capacity planning
- Cost optimization through better placement

Years 2-3: Performance optimizations
- RDMA in same-region replication
- Adaptive consistency levels per transaction class
- Edge read replicas with causal consistency

Years 4-5: Next-gen infrastructure
- Integration with edge computing (5G MEC)
- Quantum-safe cryptography migration
- ML-driven failure prediction and mitigation
- Serverless deployment model for elastic scaling

This design balances theoretical soundness with practical engineering, providing a clear path from MVP to mature system while maintaining correctness guarantees throughout evolution."

## sample_acceptable - Minimum Acceptable
"I'd use a hierarchical architecture with regional Raft clusters for local writes and a global coordination layer for cross-region transactions. Each region would have 5-node clusters with 3-way replication. For global consistency, implement Spanner-like TrueTime with bounded clock uncertainty, using 2PC over Paxos for distributed transactions. Achieve 100k writes through sharding (1000 shards handling 100 ops each). Ensure availability through automatic failover, witness replicas, and graceful degradation. Meet 5-second RTO with hot standbys and automated runbooks. Key trade-offs include complexity, higher write latency for global operations, and increased infrastructure costs."

## common_mistakes - Watch Out For
- Over-simplifying clock synchronization challenges
- Ignoring economic factors (bandwidth costs)
- Not addressing data sovereignty requirements
- Unrealistic latency expectations

## follow_up_excellent - Depth Probe
**Question**: "Your design has 100-200ms latency for cross-region transactions. A competitor claims they achieve 50ms. What might they be sacrificing and would you make that trade-off?"
- **Looking for**: Consistency relaxation, limited geographic scope, or unrealistic claims
- **Red flags**: Not recognizing physics limitations

## follow_up_partial - Guided Probe  
**Question**: "How would your system handle a scenario where the entire Europe region loses connectivity to other continents but remains internally connected?"
- **Hint embedded**: Regional autonomy vs global consistency
- **Concept testing**: Partition tolerance strategies

## follow_up_weak - Foundation Check
**Question**: "If you had to choose between everyone seeing the same data (consistency) or the system always responding (availability), which would you pick for a bank?"
- **Simplification**: CAP theorem basics
- **Building block**: Understanding trade-offs

## bar_raiser_question - L7→Industry Leader Challenge
"Your system is successful, processing $1 trillion daily. Regulators now require: 1) Cryptographic proof of transaction ordering, 2) Multi-jurisdiction data sovereignty with computation sovereignty, 3) Quantum-resistant security. How do you evolve your architecture while maintaining service?"

### bar_raiser_concepts
- Blockchain integration for verifiable logs
- Secure multiparty computation for cross-border processing
- Post-quantum cryptography migration
- Homomorphic encryption for computation on encrypted data
- Zero-knowledge proofs for compliance without data exposure

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 7-8 min answer + 3-4 min discussion
- **Common next topics**: Blockchain consensus, global state machines, planet-scale systems

## assistant_answer
Use hierarchical consensus: regional Raft/Paxos shards for locality and throughput, plus a lightweight global coordinator for cross-region ordering. Employ MVCC, locality-aware sharding, DR replicas for RTO≤5s, and batching/pipelining to scale to 100k writes/s; optimize leader placement and failover to meet 99.99% availability.

## improvement_suggestions
- Require explicit SLA math (quorum RTTs, shard counts, bandwidth budgets) tied to 100k writes/s and 99.99% availability.
- Ask for a 5-year evolution plan (time sync choices, PQ crypto migration, data sovereignty strategy).
