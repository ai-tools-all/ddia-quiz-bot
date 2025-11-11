---
id: memcached-look-aside-pattern
day: 1
tags: [memcached, caching, architecture, design-patterns]
---

# Look-Aside Caching Pattern

## question
Why does Facebook use a look-aside caching pattern where front-ends manage the cache-database relationship, rather than a look-through pattern where the cache itself fetches from the database?

## options
- A) Look-through caching is more complex to implement at scale
- B) Look-aside allows front-ends to cache arbitrary transformations of database records, not just raw rows
- C) Look-through caching requires more network round trips
- D) Look-aside caching provides better consistency guarantees

## answer
B

## explanation
Look-aside caching decouples the cache from the database, giving front-end PHP code the flexibility to store transformed or aggregated data in memcached rather than being limited to raw database rows. The front-end checks memcached first, and on a miss, fetches from MySQL, applies any necessary transformations, and populates the cache. This pattern allows caching of computed results, joins, and application-specific data structures.

## hook
What are the trade-offs between giving clients control over caching logic versus centralizing it in the cache layer?
