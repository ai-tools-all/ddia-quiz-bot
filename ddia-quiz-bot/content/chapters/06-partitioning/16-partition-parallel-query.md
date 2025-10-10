---
id: ch06-partition-parallel-query
day: 16
tags: [partitioning, query-execution, parallel-processing]
related_stories: []
---

# Parallel Query Execution

## question
In a partitioned database, what type of query benefits most from parallel execution across partitions?

## options
- A) Point lookups by primary key
- B) Aggregations across the entire dataset
- C) Updates to a single row
- D) Index lookups for a specific value

## answer
B

## explanation
Aggregation queries (like COUNT, SUM, AVG across all data) benefit greatly from parallel execution because the work can be distributed across all partitions, with each partition computing partial results that are then combined. Point lookups and single-row operations typically touch only one partition, so parallelization doesn't help. Full-dataset operations see near-linear speedup with partition count.

## hook
How does Apache Spark achieve sub-second response times for aggregations over terabytes of data?
