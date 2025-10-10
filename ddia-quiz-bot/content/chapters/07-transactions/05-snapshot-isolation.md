---
id: ch07-snapshot-isolation
day: 5
tags: [transactions, snapshot-isolation, mvcc, consistency]
related_stories: []
---

# Snapshot Isolation

## question
How does snapshot isolation (MVCC) handle long-running read queries while other transactions are modifying data?

## options
- A) Blocks all write transactions until the read completes
- B) Reads see a consistent snapshot from when the query started
- C) Returns an error if data changes during the read
- D) Reads see all changes as they happen

## answer
B

## explanation
Snapshot isolation using Multi-Version Concurrency Control (MVCC) allows read queries to see a consistent snapshot of the database from the moment the transaction began. Writers don't block readers and readers don't block writers. Each transaction reads from its own consistent snapshot, preventing phenomena like non-repeatable reads while maintaining good performance.

## hook
How does PostgreSQL implement MVCC without requiring a separate version store like some databases?
