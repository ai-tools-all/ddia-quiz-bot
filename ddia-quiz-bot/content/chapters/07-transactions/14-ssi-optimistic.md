---
id: ch07-ssi-optimistic
day: 14
tags: [transactions, ssi, serializable-snapshot-isolation, optimistic-control]
related_stories: []
---

# Serializable Snapshot Isolation (SSI)

## question
How does Serializable Snapshot Isolation (SSI) detect conflicts compared to two-phase locking?

## options
- A) Locks all data before reading
- B) Detects conflicts at commit time based on read/write sets
- C) Prevents all concurrent transactions
- D) Requires manual conflict resolution

## answer
B

## explanation
SSI is an optimistic concurrency control mechanism. Transactions proceed without locking, tracking what they read and write. At commit time, the database checks if the transaction's reads are still valid (no concurrent writes to read data) and its writes don't conflict with concurrent reads. If conflicts are detected, the transaction aborts. This provides serializability with better performance than 2PL.

## hook
Why did PostgreSQL 9.1 add SSI as its serializable isolation implementation?
