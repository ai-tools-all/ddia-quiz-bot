---
id: craq-subjective-L4-006
type: subjective
level: L4
category: baseline
topic: craq
subtopic: consistency-comparison
estimated_time: 6-8 minutes
---

# question_title - CRAQ vs Multi-Leader Replication

## main_question - Core Question
"Contrast CRAQ's read scaling approach with the multi-leader replication model from DDIA. Under which workload patterns would you prefer CRAQ, and when would multi-leader be superior?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Consistency Model**: CRAQ remains single-writer linearizable; multi-leader favors AP with conflict resolution
- **Read Distribution**: CRAQ provides read scaling via clean replicas; multi-leader via independent leaders
- **Conflict Handling**: CRAQ avoids conflicts; multi-leader requires reconciliation strategies
- **Workload Fit**: CRAQ excels with read-heavy workloads needing strong consistency

### expected_keywords
- Primary keywords: linearizability, conflict resolution, workload pattern, latency
- Technical terms: CRDTs, last-write-wins, write contention, synchronous replication

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Geographic Considerations**: CRAQ for strong cross-region reads; multi-leader for low-latency writes
- **Operational Complexity**: Conflict resolution tooling vs clean metadata tracking
- **Client Simplification**: CRAQ clients simpler due to no merge logic
- **Failure Modes**: Partition behavior differences (CP vs AP)

### bonus_keywords
- Implementation: vector clocks, causal ordering, chain length, read quorum
- Scenarios: collaborative editing, financial ledger, social timelines
- Trade-offs: latency, write amplification, operational burden

## sample_excellent - Example Excellence
"CRAQ keeps a single write path (head→tail) and scales reads by letting any clean replica answer, preserving linearizability. Multi-leader replication from DDIA allows each datacenter to accept writes independently, boosting write availability but demanding conflict resolution (LWW timestamps, CRDTs, or custom merges). I'd choose CRAQ for financial ledgers or inventory systems where correctness trumps availability but read throughput must scale. Multi-leader fits collaborative apps or social timelines where latency and availability dominate and conflicts can be resolved programmatically." 

## sample_acceptable - Minimum Acceptable
"CRAQ keeps one writer and scales reads, so it stays strongly consistent. Multi-leader lets multiple sites write at once, but you have to handle conflicts. CRAQ is better when you need linearizable reads; multi-leader helps when you need low-latency writes and can cope with conflicts." 

## common_mistakes - Watch Out For
- Claiming CRAQ allows concurrent writes from multiple heads
- Ignoring conflict resolution cost in multi-leader
- Forgetting partition behavior differences
- Overlooking workload-based recommendation

## follow_up_excellent - Depth Probe
**Question**: "How would you integrate CRAQ with a downstream multi-leader cache layer without reintroducing conflicts?"
- **Looking for**: Read-through caching, invalidation tied to tail commit, avoid independent writes
- **Red flags**: Letting caches perform writes

## follow_up_partial - Guided Probe  
**Question**: "Which metrics would you compare to decide between CRAQ and multi-leader deployment?"
- **Hint embedded**: Write rate, conflict frequency, latency requirements
- **Concept testing**: Quantitative decision making

## follow_up_weak - Foundation Check
**Question**: "If two chefs change a recipe separately, what coordination do you need to keep the final dish consistent?"
- **Simplification**: Conflict resolution analogy
- **Building block**: Single vs multi writer trade-offs

## bar_raiser_question - L4→L5 Challenge
"Propose a hybrid architecture that uses CRAQ for critical tables and multi-leader for social features in the same product, ensuring cross-component consistency."

### bar_raiser_concepts
- Data domain segmentation
- Consistency boundaries
- Event propagation between systems
- Operational governance

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Hybrid architectures, conflict mitigation, SLA alignment
