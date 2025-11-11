---
id: memcached-regional-pool
day: 1
tags: [memcached, regional-pool, memory-efficiency, cold-data]
---

# Regional Pool for Cold Data

## question
Why does Facebook use a regional pool of memcached servers shared across clusters for infrequently accessed data?

## options
- A) To improve cache hit rates by centralizing all data in one location
- B) To avoid wasteful memory replication of cold data across all clusters while dedicating per-cluster RAM to hot keys
- C) To reduce network latency by placing cold data closer to the database
- D) To simplify cluster management by separating hot and cold data into different systems

## answer
B

## explanation
Replicating all data across all clusters would waste RAM on cold (rarely accessed) keys that don't need the parallelism of replication. Facebook uses a regional pool—a shared memcached layer accessible to all clusters within a region—for infrequently accessed data. This avoids duplicating cold data across every cluster's memory while allowing per-cluster RAM to be dedicated to hot keys that benefit from replication for throughput. The regional pool acts as a shared second-tier cache: miss in local cluster, check regional pool, then query database. This architecture optimizes memory efficiency against serving capacity.

## hook
How would you determine the threshold for classifying a key as "hot" versus "cold" for this architecture?
