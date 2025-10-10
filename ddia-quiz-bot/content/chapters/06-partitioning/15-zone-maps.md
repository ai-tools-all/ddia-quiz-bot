---
id: ch06-zone-maps
day: 15
tags: [partitioning, optimization, metadata]
related_stories: []
---

# Zone Maps and Partition Pruning

## question
How do zone maps (partition metadata) optimize query performance in partitioned systems?

## options
- A) By caching frequently accessed data
- B) By tracking min/max values per partition to skip irrelevant partitions
- C) By compressing partition data
- D) By replicating hot partitions

## answer
B

## explanation
Zone maps maintain metadata about each partition such as minimum and maximum values for columns. Query optimizers use this metadata to skip partitions that cannot contain relevant data (partition pruning). For example, if querying for orders from 2024 and a partition's zone map shows it only contains 2022 data, that partition can be skipped entirely without reading any data.

## hook
How much query performance improvement can partition pruning provide in a petabyte-scale data warehouse?
