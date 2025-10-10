---
id: ch06-key-range-partitioning
day: 2
tags: [partitioning, key-range, range-queries]
related_stories: []
---

# Key Range Partitioning

## question
When using key range partitioning, what is a significant advantage compared to hash partitioning?

## options
- A) It completely eliminates hot spots
- B) It enables efficient range queries
- C) It requires less memory
- D) It provides perfect load distribution

## answer
B

## explanation
Key range partitioning assigns a continuous range of keys to each partition (like encyclopedia volumes A-C, D-F, etc.). This preserves the sort order of keys, making range queries efficient since all keys in a range will be on the same partition or adjacent partitions. Hash partitioning, while better for load distribution, scatters adjacent keys across partitions, making range queries inefficient.

## hook
Why do time-series databases often prefer range partitioning over hash partitioning?
