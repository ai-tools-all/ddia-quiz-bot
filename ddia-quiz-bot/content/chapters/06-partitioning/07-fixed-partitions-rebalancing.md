---
id: ch06-fixed-partitions-rebalancing
day: 7
tags: [partitioning, rebalancing, operational]
related_stories: []
---

# Fixed Number of Partitions

## question
When using a fixed number of partitions strategy (like consistent hashing with virtual nodes), what is a key operational advantage?

## options
- A) Partitions are always perfectly balanced
- B) Rebalancing only moves entire partitions without splitting
- C) No configuration is required
- D) It works well with any data distribution

## answer
B

## explanation
With fixed number of partitions, you create many more partitions than nodes from the start (e.g., 1000 partitions for 10 nodes). When adding/removing nodes, you only move entire partitions between nodes without splitting or merging them. This simplifies operations and makes rebalancing predictable. The partition size remains constant, only their assignment to nodes changes.

## hook
Why does Riak use 64 partitions per node by default regardless of cluster size?
