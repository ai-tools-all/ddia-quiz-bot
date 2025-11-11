---
id: memcached-read-your-writes
day: 1
tags: [memcached, consistency, user-experience]
---

# Read-Your-Writes Consistency

## question
How does Facebook's memcached system guarantee read-your-writes consistency—that users immediately see their own updates?

## options
- A) By routing all reads from a user to the same MySQL master that handled their write
- B) By front-ends issuing deletes to memcached immediately after database writes, forcing cache misses that fetch fresh data
- C) By using synchronous replication between MySQL and memcached
- D) By temporarily disabling caching for users who recently performed writes

## answer
B

## explanation
Front-ends send writes to MySQL first, then immediately issue deletes to memcached for affected keys. This invalidation forces the user's next read to miss the cache and fetch the fresh value directly from the database. The front-end delete ensures that even if asynchronous MySQL invalidations are delayed or lost, the user who just wrote will see their update immediately. This pattern specifically addresses the critical UX requirement that users see their own actions reflected without delay.

## hook
What happens if the front-end delete fails due to a network issue? How does the system eventually recover consistency?
