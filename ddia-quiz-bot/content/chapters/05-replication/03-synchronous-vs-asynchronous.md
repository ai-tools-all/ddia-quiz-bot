---
id: ch05-sync-vs-async-replication
day: 3
tags: [replication, synchronous, asynchronous, latency, durability]
related_stories: []
---

# Synchronous vs Asynchronous Replication

## question
What is the main trade-off between synchronous and asynchronous replication?

## options
- A) Synchronous is cheaper but slower
- B) Asynchronous guarantees data durability but has higher latency
- C) Synchronous guarantees data durability but blocks writes if a follower is down
- D) There is no significant difference between them

## answer
C

## explanation
Synchronous replication guarantees that data is safely stored on at least one follower before confirming the write, ensuring no data loss if the leader fails. However, if the synchronous follower is down or slow, writes are blocked. Asynchronous replication doesn't block writes but risks data loss if the leader fails before replication completes.

## hook
Why might you choose semi-synchronous replication over fully synchronous?
