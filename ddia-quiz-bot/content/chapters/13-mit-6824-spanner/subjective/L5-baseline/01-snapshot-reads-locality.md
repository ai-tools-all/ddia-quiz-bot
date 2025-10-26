---
id: spanner-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: spanner
subtopic: snapshot-reads
estimated_time: 10-12 minutes
---

# question_title - Low-Latency Consistent Reads

## main_question - Core Question
"Design Spanner’s read-only transaction path that provides low-latency local reads without violating consistency across shards. What metadata does each replica track, and when must a replica wait?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Assign a read timestamp and read versions with ts ≤ read ts
- Replica freshness tracking up to timestamp T (applied-through)
- Wait (if needed) until replica has applied all updates ≤ T
- No locks required for read-only transactions

### expected_keywords
- Primary: snapshot, timestamp, freshness, local read
- Technical: applied-through watermark, leader vs follower, staleness bounds

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Trade-off between lower latency and waiting for freshness
- Impact of cross-region latency on catching up
- Interaction with TrueTime-based start rule for reads
- Handling range/secondary index reads consistently

### bonus_keywords
- replica lag, majority replication, consistency point-in-time

## sample_excellent - Example Excellence
"The client picks a read timestamp T (e.g., per start rule). Each replica tracks a watermark indicating it has applied all Paxos entries ≤ W. If W ≥ T, it can serve reads immediately from local state by returning the latest versions with ts ≤ T. Otherwise, it delays until W reaches T. Because the snapshot is as-of T, no read locks are needed and different shards remain consistent at the same logical time."

## sample_acceptable - Minimum Acceptable
"Pick a timestamp and read a consistent snapshot. Replicas wait until they’re up-to-date for that timestamp before answering."

## common_mistakes - Watch Out For
- Requiring leaders for all reads
- Ignoring cross-shard consistency of the snapshot
- Using locks for read-only traffic

## follow_up_excellent - Depth Probe
**Question**: "How would you expose and bound read latency for users given variable replica lag?"
- **Looking for**: SLOs based on freshness metrics, fallback to leaders, dynamic routing

## follow_up_partial - Guided Probe  
**Question**: "When would a follower be unable to serve a read at timestamp T immediately?"
- **Hint embedded**: Compare T to applied-through watermark

## follow_up_weak - Foundation Check
**Question**: "What does a consistent snapshot mean in a distributed database?"
