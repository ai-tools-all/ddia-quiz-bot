---
id: memcached-thundering-herd
day: 1
tags: [memcached, thundering-herd, lease, cache-miss]
---

# Thundering Herd Problem

## question
In Facebook's memcached system, what problem occurs when a popular key is deleted and how is it mitigated?

## options
- A) Multiple clients simultaneously miss the cache and all query the database, causing a thundering herd; mitigated by having clients acquire leases before fetching from the database
- B) The deletion propagates too slowly across clusters, causing inconsistency; mitigated by using synchronous deletes
- C) The database becomes overloaded with delete operations; mitigated by batching deletes
- D) Clients repeatedly retry failed delete operations; mitigated by making deletes idempotent

## answer
A

## explanation
When a popular key is deleted from memcached, many clients simultaneously experience cache misses and would all rush to query the database, causing a thundering herd that could overwhelm the database. Facebook mitigates this using a lease mechanism: when a client experiences a cache miss, it requests a lease (token) from memcached. Only the client holding the lease is allowed to fetch from the database and populate the cache. Other clients either wait briefly and retry (expecting the lease holder to populate the cache) or use stale data if available. This prevents duplicate expensive database queries.

## hook
How would you design a lease system that balances database protection with cache freshness requirements?
