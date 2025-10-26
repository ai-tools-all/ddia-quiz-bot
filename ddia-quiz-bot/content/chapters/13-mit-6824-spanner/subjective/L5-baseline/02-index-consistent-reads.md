---
id: spanner-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: spanner
subtopic: secondary-index-consistency
estimated_time: 10-12 minutes
---

# question_title - Consistent Secondary-Index Reads

## main_question - Core Question
"Design a read-only path that joins a secondary index with base rows in Spanner while preserving a single consistent snapshot across shards. What timestamps and freshness checks are required? What anomalies occur if components read at different times?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Single read timestamp T used for both index and base tables
- Each replica must be fresh through T (applied-through ≥ T) before answering
- No locks for read-only; versions chosen with ts ≤ T
- Same T prevents mismatched base/index snapshots

### expected_keywords
- snapshot timestamp, applied-through watermark, local reads, consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Handling backfilled/partially built indexes (read-at-T gating)
- Fallback to leaders when followers lag
- Impact on latency and routing strategy

### bonus_keywords
- replica lag, staleness SLO, index backfill, follower read

## sample_excellent - Example Excellence
"Pick a global T (per start rule). Resolve index ranges at replicas whose applied-through ≥ T; fetch base rows at the same T. If any shard’s watermark < T, wait or route to a fresher replica/leader. Reading index at T1 and base at T2 (T1≠T2) can yield phantom/ghost entries or missing rows, violating snapshot semantics."

## sample_acceptable - Minimum Acceptable
"Use one timestamp for index and base and wait until replicas are fresh through that time."

## common_mistakes - Watch Out For
- Using different timestamps for index vs base
- Assuming leaders are always required
- Adding read locks unnecessarily

## follow_up_excellent - Depth Probe
**Question**: "How would you expose partial freshness to the query planner to pick routes?"
- **Looking for**: Watermark-based routing, SLO-aware fallback

## follow_up_partial - Guided Probe  
**Question**: "When must you route to a leader rather than wait on a follower?"
- **Hint embedded**: Consider tail-latency budgets vs staleness

## follow_up_weak - Foundation Check
**Question**: "What is a consistent snapshot?"
