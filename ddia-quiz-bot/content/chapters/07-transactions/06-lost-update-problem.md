---
id: ch07-lost-update-problem
day: 6
tags: [transactions, lost-update, concurrency, race-conditions]
related_stories: []
---

# Lost Update Problem

## question
Two users simultaneously increment a counter that starts at 100. User A reads 100, adds 1; User B reads 100, adds 1. Both write back their results. What's the final value and what is this problem called?

## options
- A) 102 - Race condition
- B) 101 - Lost update
- C) 102 - Successful execution
- D) 100 - Write skew

## answer
B

## explanation
This is the classic lost update problem. Both transactions read the same initial value (100), calculate new values independently (101), and write back. The second write overwrites the first, "losing" one update. The counter ends at 101 instead of 102. Solutions include atomic operations, explicit locking, or compare-and-set operations.

## hook
How do databases like Redis solve this with atomic INCR operations?
