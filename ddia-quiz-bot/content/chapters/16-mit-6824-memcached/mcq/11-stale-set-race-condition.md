---
id: memcached-stale-set-race
day: 1
tags: [memcached, race-condition, stale-data, ordering]
---

# Stale Set Race Condition

## question
Consider two clients (A and B) concurrently updating the same key. Client A writes value 1 to MySQL, then sets it in memcached. Client B writes value 2 to MySQL, then sets it in memcached. The MySQL commits complete in order (A then B), but the memcached sets arrive out of order (B's set arrives before A's set). What is the outcome and why does Facebook use deletes instead?

## options
- A) Memcached ends up with value 1 (stale), but deletes avoid this by forcing cache misses that fetch the correct value 2 from the database
- B) Memcached correctly has value 2 because sets include timestamps that resolve ordering
- C) The race is harmless because MySQL transactions ensure serializability
- D) Memcached rejects out-of-order sets using version numbers

## answer
A

## explanation
With set-based updates, network timing can cause memcached to cache a stale value: B's set arrives first (value 2), then A's delayed set overwrites it with value 1, even though the database correctly has value 2. This leaves stale data cached until timeout. By using delete-based invalidation instead, both clients send deletes (which are idempotent and order-independent). Subsequent readers miss the cache and fetch the correct value from the database, which always has the right answer. Deletes convert a potential long-lived inconsistency into a brief cache miss.

## hook
Could you use versioned cache entries (compare-and-set) to safely use set-based updates? What would be the complexity?
