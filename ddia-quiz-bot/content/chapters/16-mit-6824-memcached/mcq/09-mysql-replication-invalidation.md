---
id: memcached-mysql-replication-invalidation
day: 1
tags: [memcached, mysql, replication-log, invalidation]
---

# MySQL Replication Log Invalidation

## question
Why does Facebook have MySQL servers monitor their replication logs and send invalidations to memcached, in addition to front-end invalidations?

## options
- A) To reduce front-end complexity by centralizing all invalidation logic in the database layer
- B) To ensure invalidations propagate even if front-end deletes fail due to network issues or bugs, providing a safety net for consistency
- C) To improve performance by batching invalidations at the database level
- D) To support cross-region replication where front-ends don't have access to all memcached clusters

## answer
B

## explanation
MySQL servers monitoring the replication log and sending asynchronous invalidations provide a crucial safety net. If a front-end fails to send a delete (due to network failure, timeout, bug in application code, or front-end crash), the stale data would remain cached indefinitely. By having MySQL servers independently send deletes based on what they observe in the replication stream, the system ensures eventual cache consistency even when front-end invalidations fail. This dual-invalidation approach (front-end for read-your-writes + MySQL async for reliability) balances immediate user feedback with long-term consistency.

## hook
What trade-offs exist between relying solely on database-driven invalidation versus the dual approach?
