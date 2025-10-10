---
id: ch09-consistency-model
day: 1
tags: [gfs, consistency, weak-consistency]
related_stories: []
---

# GFS Consistency Model

## question
How does GFS characterize its consistency model for concurrent writes?

## options
- A) Strong consistency - all clients see the same data immediately
- B) Relaxed consistency - defined but possibly inconsistent regions
- C) Eventual consistency - converges over time
- D) Causal consistency - preserves cause-effect relationships

## answer
B

## explanation
GFS provides a relaxed consistency model where file regions can be "defined" (all replicas have the same data) but not necessarily "consistent" (all clients see the same data). After concurrent writes, regions may be undefined or defined but inconsistent.

## hook
Can a file system be useful even if different readers might see different data?
