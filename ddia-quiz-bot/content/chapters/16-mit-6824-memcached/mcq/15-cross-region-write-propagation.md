---
id: memcached-cross-region-write-propagation
day: 1
tags: [memcached, multi-region, replication-lag, writes]
---

# Cross-Region Write Propagation

## question
User in the East Coast region writes a post. The write goes to the West Coast primary MySQL, which then replicates to East Coast MySQL. What cache consistency issues arise, and how does the system handle the user seeing their own write?

## options
- A) East Coast memcached has stale data until async replication completes; user sees their write immediately because the front-end delete ensures a cache miss that fetches from (possibly stale) local MySQL
- B) The system synchronously replicates to all regions before acknowledging the write to ensure immediate consistency
- C) East Coast reads are redirected to West Coast primary until replication catches up
- D) The write is rejected if the user is not in the primary region

## answer
A

## explanation
This scenario reveals the interaction between regional replication and caching. The user's write goes to West Coast primary MySQL, and the front-end immediately issues a delete to the East Coast memcached (where the user is located). This ensures read-your-writes: the user's next read misses the cache and queries their local East Coast MySQL. However, there's a window where the local MySQL replica might not yet have the update due to async replication lag. The user might briefly see old data, but typically replication lag is short (seconds). The front-end delete guarantees bypassing the cache; the replication lag is a separate issue. Some systems solve this by tracking replication position or temporarily reading from primary for users who just wrote.

## hook
How could you extend the system to guarantee cross-region read-your-writes even during replication lag?
