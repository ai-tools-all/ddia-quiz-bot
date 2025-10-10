---
id: gfs-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: gfs
subtopic: consistency-architecture
estimated_time: 7-10 minutes
---

# question_title - GFS Consistency Trade-off Analysis

## main_question - Core Question
"GFS chose relaxed consistency over strong consistency. Walk me through the technical and business reasoning behind this decision. What would be different if they had chosen strong consistency?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Performance Impact**: Coordination overhead for strong consistency
- **Availability Trade-off**: CAP theorem implications
- **Target Workload**: MapReduce and batch processing tolerance
- **Scale Requirements**: Thousands of nodes make coordination expensive

### expected_keywords
- Primary keywords: CAP theorem, coordination, throughput, latency
- Technical terms: consensus protocols, atomic operations, linearizability

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Economic Factors**: Hardware costs vs software complexity
- **Network Partitions**: Handling split-brain scenarios
- **Specific Use Cases**: Web indexing, log processing
- **Alternative Approaches**: Eventual consistency, causal consistency
- **Evolution Path**: How Colossus/Spanner addressed this

### bonus_keywords
- Implementation: Paxos, Raft, two-phase commit
- Metrics: Throughput vs latency, P99 performance
- Business: Time to market, operational complexity

## sample_excellent - Example Excellence
"GFS's relaxed consistency was a deliberate choice optimizing for Google's specific needs. Technical reasoning: Strong consistency would require protocols like Paxos/2PC for every write, adding multiple network round-trips and reducing throughput by 10-100x. With thousands of nodes, coordinator failures would frequently stall the system. Business reasoning: Google's primary workload (MapReduce for web indexing) could tolerate inconsistencies through checksums and retries. The engineering cost of building strongly consistent storage would have delayed deployment by years. If they'd chosen strong consistency, GFS would have lower throughput, higher latency, require more complex failure handling, and probably wouldn't scale to thousands of nodes. This is why Google later built Spanner for workloads needing strong consistency - different requirements, different system."

## sample_acceptable - Minimum Acceptable
"GFS chose relaxed consistency because strong consistency would require coordination between replicas on every write, significantly reducing performance. Their MapReduce workloads could handle inconsistencies, so they optimized for throughput and availability instead of consistency."

## common_mistakes - Watch Out For
- Not explaining the specific trade-offs
- Missing the workload-specific reasoning
- Oversimplifying to just "performance"
- Not understanding coordination costs

## follow_up_excellent - Depth Probe
**Question**: "Given GFS's choice, how would you modify it to support a financial transaction system that requires strong consistency for account balances?"
- **Looking for**: Hybrid approach, separate systems, consistency zones
- **Red flags**: Trying to force GFS to do something it wasn't designed for

## follow_up_partial - Guided Probe
**Question**: "You mentioned coordination overhead. Can you estimate the network round-trips needed for a strongly consistent write across 3 replicas in different datacenters?"
- **Hint embedded**: Leader election, prepare, commit phases
- **Concept testing**: Consensus protocol understanding

## follow_up_weak - Foundation Check
**Question**: "Let's consider a simpler scenario: three friends maintaining a shared shopping list. What happens if two update it simultaneously without coordination?"
- **Simplification**: Human-scale consistency problem
- **Building block**: Conflict scenarios

## bar_raiser_question - L4→L5 Challenge
"Design a storage system that provides both GFS-style relaxed consistency for bulk data AND strong consistency for metadata/critical paths. How would you architect this hybrid system?"

### bar_raiser_concepts
- Separate storage tiers with different guarantees
- Consistency boundaries and transaction domains
- Performance isolation between systems
- Operational complexity of hybrid approach
- Migration paths between consistency levels

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-5 min discussion
- **Common next topics**: Spanner, NewSQL systems, consistency models spectrum
