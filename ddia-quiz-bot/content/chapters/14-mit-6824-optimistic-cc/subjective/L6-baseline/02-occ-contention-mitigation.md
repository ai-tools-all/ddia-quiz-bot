---
id: farm-subjective-L6-002
type: subjective
level: L6
category: baseline
topic: performance
subtopic: contention-mitigation
estimated_time: 12-15 minutes
---

# question_title - OCC Contention Mitigation Strategies

## main_question - Core Question
"Design a comprehensive strategy to mitigate high abort rates in FaRM when a workload exhibits severe contention on a small set of hot objects (e.g., a popular user's profile being updated thousands of times per second). Consider both application-level and system-level approaches. For each approach, analyze the trade-offs in consistency, complexity, and performance."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Application-level: batching/combining updates, eventual consistency for non-critical fields, CRDT-style operations
- System-level: partition hot objects, early conflict detection, adaptive retry with backoff, short transactions
- Pessimistic locking for known hot spots (hybrid OCC/2PL)
- Trade-offs: weaker consistency vs throughput, complexity vs maintainability, latency vs abort rate

### expected_keywords
- batching, partitioning, backoff, hot spot, adaptive, hybrid locking, CRDT, eventual consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Monitoring and identifying hot objects dynamically
- Request coalescing at proxy layer
- Application redesign to avoid hot spots (e.g., sharding user IDs)
- Hybrid storage: hot objects in pessimistic system, cold in FaRM
- Timestamp-based conflict resolution instead of abort

### bonus_keywords
- monitoring, coalescing, proxy, sharding, timestamp ordering, priority scheduling

## sample_excellent - Example Excellence
"Multi-pronged approach: (1) Application: batch updates to hot objects every 10ms, use CRDTs for counters (e.g., likes), separate read and write paths. (2) System: detect hot objects via abort metrics, apply pessimistic locks to top 1% hottest objects (hybrid mode), implement exponential backoff (1-10ms) on retry. (3) Data model: partition hot user data (profile vs posts vs metadata) to reduce conflicts. Trade-offs: batching adds latency (10ms), CRDTs lose strong consistency, hybrid locking adds complexity, monitoring adds overhead. Net result: 10x throughput increase, 50% of transactions delayed by batching."

## sample_acceptable - Minimum Acceptable
"Use batching to combine updates, add backoff delays on retries, and partition the data to spread load. These reduce conflicts but add latency and complexity."

## common_mistakes - Watch Out For
- Only suggesting one approach without considering others
- Not analyzing trade-offs for each strategy
- Ignoring application-level redesign options
- Assuming you can completely eliminate aborts
- Missing the option of hybrid OCC/pessimistic systems

## follow_up_excellent - Depth Probe
**Question**: "How would you decide dynamically which objects should use pessimistic locking vs OCC in a hybrid system?"
- **Looking for**: Online monitoring of abort rates per object, threshold-based switching (e.g., >10% abort rate), hysteresis to avoid thrashing, TTL-based policies

## follow_up_partial - Guided Probe
**Question**: "Why does exponential backoff help with contention?"
- **Hint embedded**: Spreads retries over time, reduces simultaneous conflicts

## follow_up_weak - Foundation Check
**Question**: "What is a hot spot and why does it cause problems in OCC?"
