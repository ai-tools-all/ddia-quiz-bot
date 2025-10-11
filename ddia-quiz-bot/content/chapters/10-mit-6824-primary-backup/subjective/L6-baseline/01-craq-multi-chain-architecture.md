---
id: craq-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: craq
subtopic: multi-chain-architecture
estimated_time: 9-12 minutes
---

# question_title - Architecting Multi-Chain CRAQ Deployments

## main_question - Core Question
"Design an architecture that scales CRAQ to thousands of shards while guaranteeing cross-shard transactions with bounded latency. Synthesize concepts from DDIA's partitioning, consensus, and transaction chapters to justify your design choices." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Shard-to-Chain Mapping**: Deterministic mapping with dynamic rebalancing (e.g., consistent hashing with virtual nodes)
- **Cross-Shard Coordination**: Use two-phase commit or saga orchestrator layered over CRAQ chains
- **Global Ordering/Metadata**: Configuration manager tracks epochs per chain and transaction coordinator enforces ordering
- **Latency Budgeting**: Bound tail ack and cross-chain coordination time using locality-aware placement

### expected_keywords
- Primary keywords: sharding, coordinator, epoch, bounded latency
- Technical terms: consistent hashing, dynamic rebalancing, 2PC, saga, metadata service

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Hotspot Redistribution**: Automated shard splits with chain cloning
- **Observability**: Cross-chain tracing, per-transaction telemetry
- **Failure Domains**: Separate fault domains for chains and coordinators
- **Backpressure**: Coordinated throttling to protect tail latency

### bonus_keywords
- Implementation: placement driver, metadata cache, commit index watermark, tail latency SLO
- Scenarios: large ecommerce catalog, global financial ledger, flash sale surge
- Trade-offs: coordination overhead, operational complexity, isolation levels

## sample_excellent - Example Excellence
"I'd partition data using consistent hashing into thousands of CRAQ chains of three replicas each, with placement tuned so each chain sits in a latency-friendly region trio. The configuration manager maintains per-chain epochs and publishes them to a transaction coordinator. For cross-shard operations, we run two-phase commit: prepare logs writes on each chain (dirty state) and only commit after all tails respond, keeping latency bounded with a maximum of two round trips per chain plus one coordinator decision. Hotspots trigger automated shard cloning and key repartitioning, while observability tracks per-transaction dirty duration. This synthesizes DDIA's partitioning, consensus, and transaction guidance." 

## sample_acceptable - Minimum Acceptable
"Partition shards across many CRAQ chains using consistent hashing and coordinate cross-shard work with 2PC: prepare logs on each chain, then commit once all tails acknowledge. Keep placement regional to bound latency and let the configuration manager track epochs." 

## common_mistakes - Watch Out For
- Ignoring coordination cost for cross-shard transactions
- No plan for hotspot rebalancing
- Missing link to DDIA concepts (partitioning, transactions)
- Overlooking latency budgets or placement design

## follow_up_excellent - Depth Probe
**Question**: "How would you evolve this design to support monotonic snapshot reads across chains without sacrificing write latency?"
- **Looking for**: Read timestamp service, hybrid logical clocks, asynchronous snapshot dissemination
- **Red flags**: Forcing global serialization on every read

## follow_up_partial - Guided Probe  
**Question**: "What metrics reveal that a rebalancing event is necessary?"
- **Hint embedded**: Tail lag variance, coordinator wait time, shard-level dirty duration
- **Concept testing**: Data-driven decision making

## follow_up_weak - Foundation Check
**Question**: "Why is it dangerous to coordinate a multi-team project without agreeing on deadlines for each task?"
- **Simplification**: Cross-shard transaction scheduling
- **Building block**: Coordinated timing ensures bounded latency

## bar_raiser_question - L6→L7 Challenge
"Extend your design to support multi-region regulatory requirements that mandate local sovereignty while keeping a global settled ledger. How do you reconcile conflicting jurisdictions?"

### bar_raiser_concepts
- Legal/regulatory partitioning
- Federated consensus, multi-level coordination
- Deferred settlement strategies
- Data sovereignty vs global consistency

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: Snapshot reads, hierarchical consensus, regulatory compliance
