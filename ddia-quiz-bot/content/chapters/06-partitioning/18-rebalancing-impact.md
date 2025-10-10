---
id: ch06-rebalancing-impact
day: 18
tags: [partitioning, rebalancing, operations, performance]
related_stories: []
---

# Rebalancing Impact

## question
During partition rebalancing in a production system, what is the primary risk to monitor?

## options
- A) Data corruption from partial transfers
- B) Network saturation affecting normal queries
- C) Running out of disk space
- D) Loss of partition metadata

## answer
B

## explanation
Rebalancing involves moving large amounts of data between nodes, which can saturate network bandwidth and impact normal query performance. This is why many systems implement rate limiting for rebalancing operations. Data corruption is prevented by protocols, disk space is planned, and metadata is replicated, but network congestion directly affects user-facing performance.

## hook
How does LinkedIn's Kafka handle rebalancing millions of partitions without impacting real-time stream processing?
