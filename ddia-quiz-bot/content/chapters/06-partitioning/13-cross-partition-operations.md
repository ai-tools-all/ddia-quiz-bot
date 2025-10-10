---
id: ch06-cross-partition-operations
day: 13
tags: [partitioning, transactions, distributed-systems]
related_stories: []
---

# Cross-Partition Operations

## question
Why are cross-partition transactions typically more expensive than single-partition transactions?

## options
- A) They require more CPU processing power
- B) They need coordination protocols like 2PC or consensus
- C) They use more disk space
- D) They require special hardware

## answer
B

## explanation
Cross-partition transactions require coordination protocols like two-phase commit (2PC) or consensus algorithms to ensure atomicity across multiple partitions. This introduces network round-trips, potential for blocking, and coordination overhead. Single-partition transactions can be handled entirely within one node using local ACID guarantees, making them much faster and simpler.

## hook
How does CockroachDB provide ACID transactions across partitions without traditional 2PC?
