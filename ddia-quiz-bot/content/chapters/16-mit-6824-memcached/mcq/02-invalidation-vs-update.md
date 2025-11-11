---
id: memcached-invalidation-protocol
day: 1
tags: [memcached, invalidation, consistency, race-conditions]
---

# Invalidation vs Update Protocol

## question
Why does Facebook's memcached system use invalidation (delete) on writes instead of updating (set) the cache with new values?

## options
- A) Delete operations are faster than set operations in memcached
- B) Concurrent writes could leave stale values in cache if two clients' sets arrive out of order relative to their database commits
- C) Invalidation uses less network bandwidth
- D) MySQL replication log only supports delete operations

## answer
B

## explanation
Invalidation is preferred because concurrent writes can cause consistency problems with updates. If two clients write to the same key, their memcached sets might arrive in a different order than their database commits completed, leaving a stale value cached. By deleting the key instead, subsequent reads will miss the cache and fetch the correct, fresh value from the database. The delete forces cache misses that naturally resolve to consistent state.

## hook
How do delete-based invalidation and set-based updates differ in their consistency guarantees when multiple clients write concurrently?
