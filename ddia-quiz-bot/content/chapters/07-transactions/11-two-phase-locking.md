---
id: ch07-two-phase-locking
day: 11
tags: [transactions, 2pl, locking, serializability]
related_stories: []
---

# Two-Phase Locking (2PL)

## question
What is the major performance drawback of two-phase locking (2PL) for achieving serializability?

## options
- A) High memory usage for lock tables
- B) Transactions must wait for locks, reducing concurrency
- C) Complex implementation requiring specialized hardware
- D) Incompatibility with indexes

## answer
B

## explanation
Two-phase locking ensures serializability but significantly reduces concurrency. Transactions must acquire locks before accessing data and hold them until commit. This means readers block writers, writers block readers, and transactions often wait in lock queues. This waiting can lead to poor performance under high concurrency and increased risk of deadlocks.

## hook
Why did many databases move away from 2PL to MVCC-based approaches?
