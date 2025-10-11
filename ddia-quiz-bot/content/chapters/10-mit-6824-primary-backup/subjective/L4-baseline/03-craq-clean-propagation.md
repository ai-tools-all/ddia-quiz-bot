---
id: craq-subjective-L4-003
type: subjective
level: L4
category: baseline
topic: craq
subtopic: clean-propagation
estimated_time: 6-8 minutes
---

# question_title - Reasoning About Clean Propagation Latency

## main_question - Core Question
"Quantify how CRAQ's clean-flag propagation latency impacts read throughput in a geographically distributed deployment. Compare it against the follower read staleness budget discussed in DDIA's replication chapter." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Propagation Delay Components**: Head→…→tail, tail ack→…→head round trip
- **Read Availability Window**: Replica unusable for reads while dirty
- **Geographic Latency Impact**: Cross-region RTT increases dirty duration
- **DDIA Link**: Similar to staleness bound calculations for follower reads

### expected_keywords
- Primary keywords: propagation delay, read throughput, staleness bound, RTT
- Technical terms: pipeline depth, batching, follower lag, service-level objective

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Latency Amortization**: Batching writes reduces per-update cost
- **Client Routing**: Weighted selection favoring replicas with low dirty windows
- **Monitoring**: Export staleness histogram akin to follower lag metrics
- **Trade-offs**: Additional metadata to track clean timestamps

### bonus_keywords
- Implementation: moving average, watermarks, percentile SLOs
- Scenarios: multi-region chain, CDN-style read distribution, hot partitions
- Trade-offs: long chains vs local read pools, throughput vs freshness balancing

## sample_excellent - Example Excellence
"In CRAQ the time from head write to tail acknowledgment defines how long each intermediate replica stays dirty. If cross-region RTT is 40 ms and the chain holds five replicas, the write requires roughly 4×40 ms to reach the tail plus another 4×40 ms for the ack to flow back, so a middle replica might remain dirty for ~320 ms. During that window, it can't serve reads, effectively reducing available throughput. DDIA's replication chapter suggests budgeting for follower staleness using max lag; CRAQ provides a deterministic bound equal to the clean propagation delay. We can batch updates to amortize the RTT or route reads preferentially to replicas with shorter dirty windows." 

## sample_acceptable - Minimum Acceptable
"Clean-flag propagation takes a head-to-tail round trip, so with high RTT the replicas stay dirty longer and can't serve reads. That's like the follower lag budget from DDIA—it's a predictable staleness window you need to monitor and address with batching or smarter routing."

## common_mistakes - Watch Out For
- Ignoring tail-to-head acknowledgment in latency budget
- Assuming replicas can serve reads while dirty
- Forgetting geographic RTT compounding through the chain
- Not referencing DDIA's follower lag discussion

## follow_up_excellent - Depth Probe
**Question**: "Design a dashboard widget that highlights when clean propagation latency threatens your read SLOs."
- **Looking for**: Dirty duration percentiles, correlation with RTT, automated rerouting
- **Red flags**: Single metric without context

## follow_up_partial - Guided Probe  
**Question**: "If clean propagation becomes too slow, which component should you optimize first?"
- **Hint embedded**: Tail acknowledgment path, cross-region network
- **Concept testing**: Bottleneck identification

## follow_up_weak - Foundation Check
**Question**: "Imagine washing dishes on an assembly line—what happens to serving speed if rinsing takes longer than drying?"
- **Simplification**: Dirty-to-clean pipeline analogy
- **Building block**: Bottleneck determines throughput

## bar_raiser_question - L4→L5 Challenge
"Given a write-heavy workload, propose a batching or pipelining strategy that reduces dirty-window time without sacrificing the failure semantics described in DDIA's log replication chapter."

### bar_raiser_concepts
- Write coalescing, watermark acknowledgments
- Preserving linearizable order
- Risk of batch loss vs per-entry ack
- Lessons from log replication optimization

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Batching design, read routing policy, monitoring strategy
