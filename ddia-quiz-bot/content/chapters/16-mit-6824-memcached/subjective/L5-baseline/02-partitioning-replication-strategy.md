---
id: memcached-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: memcached
subtopic: partitioning-replication
estimated_time: 10-12 minutes
---

# question_title - Partitioning vs Replication Strategy

## main_question - Core Question
"Compare Facebook's use of partitioning (sharding) and replication for memcached. When is each strategy appropriate, and how do they address different scalability challenges? What are the trade-offs in terms of capacity, consistency, and hot-key handling?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Partitioning within clusters splits data for capacity and balanced load
- Replication across clusters multiplies serving capacity for hot keys
- Partitioning provides independent parallelism but can't handle hot spots
- Replication helps performance but complicates consistency
- Regional pool caches cold data to avoid wasteful replication

### expected_keywords
- partitioning, sharding, replication, hot keys, capacity, consistency, regional pool

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Databases use sharding for write load and data volume
- Hot keys would bottleneck a single partition
- Memory efficiency vs performance trade-off
- Replication creates multiple copies to synchronize
- Hybrid approach balances concerns

### bonus_keywords
- write scalability, bottleneck, memory overhead, synchronization, hybrid strategy

## sample_excellent - Example Excellence
"Facebook employs both partitioning and replication strategically. Within each cluster, data is partitioned (sharded) across memcached servers to distribute capacity and balance load. Partitioning scales total storage and provides independent parallelism, but cannot alleviate hot spots—if one key receives massive traffic, it bottlenecks on a single server. For popular keys, Facebook replicates them across multiple clusters, multiplying serving capacity by distributing load. However, replication has costs: it consumes RAM proportionally to replica count and complicates consistency since multiple copies must be synchronized. To balance these, cold (infrequently accessed) data goes to a shared regional pool, avoiding wasteful replication and dedicating per-cluster RAM to hot items. Databases similarly use sharding for write load and data volume. This hybrid approach illustrates the tension between performance (which replication aids) and consistency/efficiency (which replication undermines)."

## sample_acceptable - Minimum Acceptable
"Partitioning splits data across servers for capacity. Replication copies hot keys across clusters for throughput. Partitioning alone can't handle hot keys. Replication uses more memory and makes consistency harder. Regional pool handles cold data efficiently."

## common_mistakes - Watch Out For
- Treating partitioning and replication as mutually exclusive
- Not explaining why partitioning alone fails for hot keys
- Missing the memory efficiency vs performance trade-off
- Not mentioning consistency complexity with replication

## follow_up_excellent - Depth Probe
**Question**: "How would you detect hot keys in a live system and dynamically decide whether to replicate them? What metrics would trigger replication?"
- **Looking for**: Request rate monitoring per key, latency percentiles, cache hit ratio, dynamic threshold adjustment

## follow_up_partial - Guided Probe
**Question**: "Why does replication make consistency harder?"
- **Hint embedded**: Multiple copies must be kept synchronized, invalidation must reach all replicas

## follow_up_weak - Foundation Check
**Question**: "What's the fundamental difference between partitioning and replication?"
