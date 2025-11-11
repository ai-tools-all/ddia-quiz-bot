---
id: memcached-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: memcached
subtopic: race-conditions
estimated_time: 12-15 minutes
---

# question_title - Cache Invalidation Race Conditions

## main_question - Core Question
"Design a robust cache invalidation protocol for a look-aside cache system. Consider race conditions between concurrent reads, writes, and deletes. How would you handle: (1) read occurring between database write and cache delete, (2) concurrent writes to the same key, (3) lost delete messages? What mechanisms would you add to Facebook's design to detect and mitigate these issues?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Race between write commit and delete propagation
- Concurrent read might populate cache with stale data
- Concurrent writes can cause out-of-order deletes/sets
- Lost deletes leave stale data indefinitely
- Mechanisms: version numbers, lease/token system, timeouts, idempotent operations

### expected_keywords
- race condition, concurrent, version, lease, timeout, staleness window, idempotent, ordering

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Thundering herd problem on cache miss
- Lease mechanism prevents duplicate computation
- Short invalidation window acceptable for most use cases
- Monitoring and detection of stale cache state
- Compare-and-set for conditional updates
- Replication log as source of truth

### bonus_keywords
- thundering herd, lease holder, monitoring, CAS, replication stream, eventual convergence

## sample_excellent - Example Excellence
"Several race conditions threaten consistency: (1) A read occurring between database write commit and cache delete might populate the cache with the old value, briefly serving stale data until the delete arrives. (2) Concurrent writes to the same key could leave the wrong value cached if their deletes/sets arrive out of order. (3) Lost delete messages (network failure) leave stale data indefinitely. Facebook's design mitigates these with dual invalidation (front-end + MySQL async deletes) and preferring deletes over sets to avoid ordering issues. To strengthen this, I'd add: version numbers stored with cached values, checked against database on cache population; lease tokens issued on cache miss, preventing thundering herd and ensuring only the lease holder populates the cache; aggressive timeouts (TTLs) as a safety net against lost deletes; monitoring to detect divergence between cache and database; idempotent delete operations to safely retry. The lease system particularly helps: on miss, client gets exclusive right to populate, preventing duplicate database loads and racing cache sets."

## sample_acceptable - Minimum Acceptable
"Race between write and delete can leave stale data temporarily. Concurrent writes need ordering. Lost deletes need timeouts. Add version checks, leases for cache misses, and monitoring."

## common_mistakes - Watch Out For
- Not identifying specific race condition scenarios
- Missing the lease/thundering herd problem
- Proposing only timeouts without version checks
- Not considering idempotency of operations
- Assuming synchronous operations eliminate all races

## follow_up_excellent - Depth Probe
**Question**: "How would you quantify the acceptable staleness window for different types of data (e.g., profile photos vs friend counts)? Design a tiered consistency system."
- **Looking for**: SLO-based approach, data classification, different TTLs, critical path vs background updates

## follow_up_partial - Guided Probe
**Question**: "What is the thundering herd problem and how do leases solve it?"
- **Hint embedded**: Multiple clients miss cache simultaneously, causing duplicate database loads

## follow_up_weak - Foundation Check
**Question**: "Draw a timeline showing the race between a write, delete, and concurrent read."
