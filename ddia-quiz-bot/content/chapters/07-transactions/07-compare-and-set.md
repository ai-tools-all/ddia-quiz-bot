---
id: ch07-compare-and-set
day: 7
tags: [transactions, compare-and-set, optimistic-locking, concurrency-control]
related_stories: []
---

# Compare-and-Set Operations

## question
A compare-and-set (CAS) operation updates a value only if it hasn't changed since it was read. What type of concurrency control does this represent?

## options
- A) Pessimistic locking
- B) Optimistic concurrency control
- C) Two-phase locking
- D) Deadlock prevention

## answer
B

## explanation
Compare-and-set is a form of optimistic concurrency control. It assumes conflicts are rare and proceeds without locking. At commit time, it checks if the value changed since reading. If unchanged, the update succeeds; if changed, the operation fails and typically retries. This avoids the overhead of locking in low-contention scenarios.

## hook
Why do modern CPUs provide CAS as a hardware instruction for lock-free data structures?
