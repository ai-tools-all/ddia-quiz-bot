---
id: ch06-partition-skew
day: 12
tags: [partitioning, skew, hot-spots, performance]
related_stories: []
---

# Partition Skew

## question
Your analytics database partitioned by date shows severe skew - today's partition handles 90% of queries while historical partitions are idle. What's the most practical solution?

## options
- A) Repartition all data by a different key
- B) Create sub-partitions for hot dates using compound keys
- C) Add more replicas to all partitions
- D) Switch to hash partitioning

## answer
B

## explanation
Using compound partitioning keys (date + another dimension like user_id range or hash) allows you to split hot partitions into smaller chunks while maintaining query efficiency. This is more practical than repartitioning everything or losing date-based query optimization. This technique is commonly used in time-series databases where recent data is accessed more frequently.

## hook
How does Apache Druid handle the "hot recent data" problem in time-series partitioning?
