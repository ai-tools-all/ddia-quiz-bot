---
id: memcached-database-ground-truth
day: 1
tags: [memcached, architecture, consistency, database]
---

# Database as Ground Truth

## question
In Facebook's look-aside caching architecture, what is the fundamental role of the MySQL database versus memcached?

## options
- A) Both are equal peers in a distributed storage system, with memcached handling reads and MySQL handling writes
- B) MySQL is the authoritative source of truth providing durability and correctness; memcached is a performance optimization that can be reconstructed from MySQL at any time
- C) Memcached is the primary data store with MySQL serving as a backup for disaster recovery
- D) MySQL and memcached maintain independent copies synchronized through two-phase commit

## answer
B

## explanation
A critical architectural principle is that MySQL is the single source of truth—the authoritative, durable store that defines correct system state. Memcached is purely a performance optimization: a cache that can be entirely lost or filled with garbage, and the system remains correct (just slower). Any cached value can be reconstructed by querying MySQL. This principle simplifies consistency reasoning: when in doubt, the cache can be invalidated and data refetched from the database. Caches can have bounded staleness through TTLs, and entire cache clusters can be blown away and rebuilt. This contrasts with systems where cache and database are peers requiring complex synchronization.

## hook
What architectural principles does treating the database as ground truth enable versus a peer-to-peer cache-database design?
