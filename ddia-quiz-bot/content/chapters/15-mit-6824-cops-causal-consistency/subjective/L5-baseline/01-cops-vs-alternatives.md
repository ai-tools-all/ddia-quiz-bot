---
id: cops-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: cops
subtopic: design-comparison
estimated_time: 10-12 minutes
---

# question_title - COPS vs. Alternative Approaches

## main_question - Core Question
"Compare COPS with (1) pure eventual consistency and (2) Strawman 2's sync barriers. For each comparison, explain what trade-offs COPS makes in terms of consistency guarantees, latency, availability, and implementation complexity."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Pure eventual: lower latency, higher availability, but causal anomalies
- COPS: prevents causal anomalies, maintains local write latency, but complexity in tracking
- Sync barrier: strong ordering, but high write latency (cross-datacenter waits)
- COPS sits in middle: stronger than eventual, weaker than linearizability
- Trade-off: COPS sacrifices availability during partitions for causal consistency

### expected_keywords
- causal anomalies, local writes, cross-datacenter latency, dependency tracking, consistency spectrum

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- COPS has no coordination on critical path (like eventual consistency)
- But COPS has deferred visibility overhead (metadata, checking)
- Sync barriers similar to Spanner/Memcache write patterns
- CAP theorem considerations: COPS chooses consistency over availability when dependencies blocked

### bonus_keywords
- coordination-free writes, metadata overhead, consistency-availability trade-off, CAP

## sample_excellent - Example Excellence
"Compared to pure eventual consistency, COPS provides stronger guarantees by preventing causal anomalies (like seeing a list entry before its photo), while maintaining the same local write latency since clients don't wait for remote acknowledgments. However, COPS adds complexity through context tracking and dependency metadata. Compared to Strawman 2's sync barriers, COPS achieves much lower write latency by not waiting for cross-datacenter propagation, but provides weaker consistency—only causal, not linearizable. COPS cannot handle all concurrent writes perfectly (uses LWW). During network partitions blocking dependencies, COPS sacrifices availability to maintain causal consistency, unlike pure eventual consistency which remains available but inconsistent."

## sample_acceptable - Minimum Acceptable
"COPS is stronger than eventual consistency (prevents causal anomalies) but weaker than sync barriers (not linearizable). COPS has low write latency like eventual consistency, but adds dependency tracking complexity. During partitions, COPS may stall writes unlike eventual consistency."

## common_mistakes - Watch Out For
- Not addressing all three dimensions: consistency, latency, availability
- Missing that COPS is not linearizable
- Not mentioning the complexity cost of dependency tracking

## follow_up_excellent - Depth Probe
**Question**: "When would you choose pure eventual consistency over COPS, despite the anomaly risks?"
- **Looking for**: Use cases where causal ordering doesn't matter, extreme scale requirements

## follow_up_partial - Guided Probe
**Question**: "Why doesn't COPS provide linearizability like Spanner?"
- **Hint embedded**: What would be required for linearizability?

## follow_up_weak - Foundation Check
**Question**: "Does COPS require cross-datacenter coordination for writes?"
