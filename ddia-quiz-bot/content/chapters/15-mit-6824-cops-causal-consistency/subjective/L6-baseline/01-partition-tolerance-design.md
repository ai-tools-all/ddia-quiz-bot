---
id: cops-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: cops
subtopic: partition-tolerance
estimated_time: 12-15 minutes
---

# question_title - Network Partitions and System Evolution

## main_question - Core Question
"You're designing a geo-replicated key-value store and considering COPS as the consistency model. Analyze how COPS behaves during a prolonged network partition between two data centers. What operational challenges arise, and how would you extend COPS to handle partitions more gracefully while preserving causal semantics?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- During partition, writes with dependencies on unreachable datacenter stall indefinitely
- Availability sacrifice: affected keys cannot become visible even if data physically present
- Creates split-brain scenario where each partition can accept writes but not see other's
- Trade-off: COPS chooses consistency over availability (CP behavior for dependent data)
- Challenge: detecting partition vs. slow network, deciding when to degrade guarantees

### expected_keywords
- partition stall, availability sacrifice, split-brain, CP behavior, degraded mode, dependency checking

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Possible mitigations: timeout-based visibility, degraded causal mode, conflict-free subsets
- Version vectors to detect partition vs. lag
- Operational monitoring: dependency wait times, visibility delays
- Comparison with AP systems (Dynamo) and CP systems (Spanner)
- Human operational burden: deciding when to intervene

### bonus_keywords
- timeout heuristics, causal cuts, version vectors, operational visibility, system observability

## sample_excellent - Example Excellence
"During a prolonged partition between DC-East and DC-West, COPS exhibits CP behavior for causally-dependent keys. If DC-East has a write W2 that depends on W1 from DC-West, W2 cannot become visible at other data centers until the partition heals and W1 propagates. This creates operational challenges: (1) detecting partition vs. slow network—need version vectors or heartbeat monitoring; (2) deciding acceptable wait times—balance freshness vs. availability; (3) managing split-brain where both partitions accept local writes but cannot exchange them.

Possible extensions: (a) Timeout-based degradation—after T seconds, make writes visible with warning flags, accepting potential causal violations; (b) Causal cuts—identify independent key subsets that can proceed without cross-partition dependencies; (c) Conflict-free dependency relaxation—allow reads to proceed with stale dependencies if writes use CRDTs. Implementation would require: metadata to track dependency sources, monitoring dashboards for visibility delays, automated or manual intervention policies, and post-partition reconciliation (like Bayou anti-entropy)."

## sample_acceptable - Minimum Acceptable
"COPS stalls visibility when dependencies are unreachable during partition, sacrificing availability. Operationally hard to detect partition vs. lag. Could extend with timeouts to degrade to eventual consistency after threshold, or use version vectors to detect partition. Need monitoring for dependency wait times."

## common_mistakes - Watch Out For
- Not addressing operational detection and intervention challenges
- Missing the CP vs. AP trade-off in COPS
- Not proposing concrete extensions with trade-off analysis
- Ignoring split-brain implications

## follow_up_excellent - Depth Probe
**Question**: "How would you design a monitoring dashboard for COPS operators to detect and respond to partition-induced stalls?"
- **Looking for**: Metrics (dependency wait time, visibility lag), alerting thresholds, manual override capabilities

## follow_up_partial - Guided Probe
**Question**: "What are the risks of using a timeout to forcibly make writes visible during partitions?"
- **Hint embedded**: Causal consistency violation consequences

## follow_up_weak - Foundation Check
**Question**: "Why does COPS stall writes during a partition, rather than applying them immediately?"
