---
id: cops-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: cops
subtopic: replica-visibility
estimated_time: 8-10 minutes
---

# question_title - Replica Visibility and Cascading Delays

## main_question - Core Question
"Describe how a cascading dependency wait can occur in COPS. Give an example where a write at one data center is delayed at a remote replica due to waiting on multiple levels of dependencies."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Write W3 depends on W2, which depends on W1
- Remote replica receives W3 but cannot make it visible until W2 is visible
- W2 cannot become visible until W1 arrives and becomes visible
- This creates a cascading chain of waits
- Delays can accumulate across dependency chain depth

### expected_keywords
- cascading wait, dependency chain, delayed visibility, transitive dependencies, blocking

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Worse if dependencies span multiple keys/shards
- Network delays or partitions amplify the problem
- Can cause staleness at remote replicas even without failures
- Depth of causal chains affects propagation latency

### bonus_keywords
- multi-hop delay, propagation stall, causal depth, cross-shard dependencies

## sample_excellent - Example Excellence
"Consider three writes: W1 (photo upload), W2 (add photo to album list), W3 (share album notification). W2 depends on W1, and W3 depends on W2. At a remote data center, if W3 arrives first via fast network path, it cannot be made visible yet. When W2 arrives, it also cannot be made visible because it depends on W1, which hasn't arrived. Only when W1 finally arrives can the system make W1 visible, then W2, then W3. This cascading wait means W3's visibility is delayed by the sum of propagation times for all upstream dependencies, even if W3's data itself arrived early."

## sample_acceptable - Minimum Acceptable
"If write C depends on B, which depends on A, a remote replica receiving C must wait for B, which must wait for A. This creates cascading delays where C's visibility depends on the full chain propagating."

## common_mistakes - Watch Out For
- Not providing a concrete example with multiple writes
- Missing that delays accumulate across chain depth
- Not explaining why intermediate writes also block

## follow_up_excellent - Depth Probe
**Question**: "How could you mitigate cascading delays without sacrificing causal consistency?"
- **Looking for**: Dependency prediction, batching, causal cuts, faster propagation

## follow_up_partial - Guided Probe
**Question**: "What happens if A is delayed by a slow network but C arrives quickly?"
- **Hint embedded**: Must preserve causal order

## follow_up_weak - Foundation Check
**Question**: "Can a write become visible if its dependencies are not yet visible?"
