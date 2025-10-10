---
id: ch05-chain-replication
day: 15
tags: [replication, chain-replication, consistency, performance]
related_stories: []
---

# Chain Replication

## question
What is a key advantage of chain replication over traditional primary-backup replication?

## options
- A) Lower storage requirements
- B) Simpler implementation
- C) Strong consistency with good read throughput by separating read and write paths
- D) Better write throughput

## answer
C

## explanation
Chain replication arranges nodes in a chain where writes go to the head and propagate through to the tail, while reads are served from the tail. This provides strong consistency (tail has all committed writes) while distributing read load. The tail can serve all reads without coordinating with the head, unlike primary-backup where the primary handles both reads and writes.

## hook
Why might chain replication struggle with node failures in the middle of the chain?
