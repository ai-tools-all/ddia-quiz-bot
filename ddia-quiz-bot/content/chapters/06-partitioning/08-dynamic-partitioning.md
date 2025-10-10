---
id: ch06-dynamic-partitioning
day: 8
tags: [partitioning, rebalancing, dynamic-partitioning]
related_stories: []
---

# Dynamic Partitioning

## question
In dynamic partitioning, when does a partition typically split?

## options
- A) When the node CPU usage exceeds threshold
- B) When the partition size exceeds a configured maximum
- C) At fixed time intervals
- D) When query latency increases

## answer
B

## explanation
Dynamic partitioning automatically splits partitions when they grow beyond a configured size threshold (e.g., 10GB). Similarly, partitions can be merged if they shrink below a minimum threshold. This approach adapts to data volume: few partitions when data is small, more partitions as data grows. HBase and MongoDB use this strategy.

## hook
How does MongoDB decide which partition key range to split when a chunk grows too large?
