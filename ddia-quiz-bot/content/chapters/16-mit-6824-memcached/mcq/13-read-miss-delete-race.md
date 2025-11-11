---
id: memcached-read-miss-delete-race
day: 1
tags: [memcached, race-condition, consistency, read-delete]
---

# Read-Miss-Delete Race

## question
Consider this sequence: (1) Client A reads key X from memcached—miss; (2) Client A queries MySQL, gets value=10; (3) Client B writes value=20 to MySQL for key X; (4) Client B deletes key X from memcached; (5) Client A (whose read is still in progress) sets X=10 in memcached. What is the result and how could this be prevented?

## options
- A) Cache has stale value 10; prevented by using version numbers or lease tokens that invalidate slow reads
- B) Cache correctly has value 20 because deletes have priority
- C) The race is impossible due to MySQL locking
- D) Cache has no value because the delete removes A's set

## answer
A

## explanation
This classic race leaves stale data cached: A's read started before B's write, but A's cache population happens after B's delete. The cache now has the old value 10 while the database has 20. Facebook's lease mechanism helps prevent this: when A gets a cache miss, it receives a lease (token). If B's delete invalidates that key, it can also invalidate A's outstanding lease. When A tries to set the cache with its (now stale) result, memcached rejects the set because the lease is invalidated. This forces A to discard its stale data, and the next reader will fetch the fresh value 20.

## hook
How long should lease tokens remain valid, and what factors influence this timeout?
