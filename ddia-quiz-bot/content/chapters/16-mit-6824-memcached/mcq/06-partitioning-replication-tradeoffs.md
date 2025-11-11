---
id: memcached-partitioning-replication
day: 1
tags: [memcached, partitioning, replication, hot-keys]
---

# Partitioning vs Replication Trade-offs

## question
How does Facebook's memcached system use both partitioning and replication to handle different data access patterns?

## options
- A) All data is uniformly partitioned across memcached servers within each cluster
- B) Data is partitioned within clusters for capacity, but popular keys are replicated across clusters to handle hot spots
- C) Replication is only used for disaster recovery, not performance
- D) Partitioning is used for read scalability while replication handles write scalability

## answer
B

## explanation
Memcached employs both strategies strategically: sharding (partitioning) within clusters distributes data for capacity and balanced load, but cannot alleviate hot spots where a single key receives massive traffic. Popular keys are replicated across multiple clusters, allowing parallel serving and multiplying throughput for hot items. Additionally, a regional pool caches infrequently accessed items shared across clusters, avoiding wasteful replication of cold data while dedicating per-cluster RAM to hot keys. This hybrid approach balances memory efficiency with performance.

## hook
Why can't sharding alone solve the hot key problem? What fundamental limitation does replication address?
