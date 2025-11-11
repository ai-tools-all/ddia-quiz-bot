---
id: memcached-transformation-caching
day: 1
tags: [memcached, look-aside, transformation, flexibility]
---

# Caching Transformed Data

## question
A Facebook page displays a user's friend list sorted by recent activity, computed by joining user data with activity logs. How does look-aside caching benefit this scenario compared to look-through caching?

## options
- A) Look-aside allows the front-end to cache the pre-computed, sorted friend list rather than just raw database rows, avoiding repeated expensive computation
- B) Look-aside provides better consistency guarantees for join operations
- C) Look-through caching cannot handle joins across multiple tables
- D) Look-aside caching has lower latency for complex queries

## answer
A

## explanation
Look-aside caching's key advantage is flexibility: the front-end controls what gets cached and can store arbitrary transformations of database data. In this scenario, the PHP code can execute the expensive query (joining users and activity logs), sort the results, apply filtering logic, and cache the final result under a synthetic key like "user:123:friends_by_activity". Subsequent requests retrieve this pre-computed result directly. With look-through caching, the cache layer would only store raw rows and couldn't cache the computed transformation—every request would repeat the join and sort. This pattern is pervasive at Facebook: caching rendered HTML fragments, aggregated counts, personalized feeds, etc.

## hook
What invalidation strategy would you use for caches of transformed/aggregated data that depend on multiple database tables?
