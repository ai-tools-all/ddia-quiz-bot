---
id: cops-context-size-growth
day: 1
tags: [cops, context, metadata-overhead, optimization]
---

# Context Size and Metadata Overhead

## question
A COPS client performs 1000 get operations on different keys, then performs a single put. What is the primary concern with this access pattern?

## options
- A) The put will fail because COPS limits context to 100 entries
- B) The context sent with the put contains 1000 key-version pairs, creating large metadata overhead
- C) The 1000 gets violate causal consistency because they're not atomic
- D) The put will be delayed until all 1000 dependencies propagate to all datacenters

## answer
B

## explanation
The client context accumulates all key-version pairs from gets, so after 1000 reads, the put's dependency list contains 1000 entries. This creates significant metadata overhead—both in network transmission and storage at replicas. Each replica must check all 1000 dependencies before making the write visible. In practice, COPS needs optimizations like: (1) garbage collecting old versions from context, (2) compacting dependencies (keeping only latest version per key), or (3) using dependency snapshots. This is a key practical challenge in implementing COPS.

## hook
How would you design a context pruning mechanism that preserves causal correctness while reducing metadata size?
