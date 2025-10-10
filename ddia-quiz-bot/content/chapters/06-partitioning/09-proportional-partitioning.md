---
id: ch06-proportional-partitioning
day: 9
tags: [partitioning, rebalancing, cassandra]
related_stories: []
---

# Partitioning Proportional to Nodes

## question
In Cassandra's partitioning approach (proportional to nodes), what happens when a new node joins the cluster?

## options
- A) All existing partitions are reshuffled equally
- B) The new node splits partition ranges with existing nodes
- C) A new partition is created only for the new node
- D) Nothing changes until manual rebalancing

## answer
B

## explanation
In partitioning proportional to nodes, each node owns a certain number of partitions. When a new node joins, it picks random existing partitions to split, taking half of the data from each split partition. This keeps the number of partitions proportional to the number of nodes, ensuring that partition size remains roughly consistent as the cluster grows.

## hook
Why did Cassandra move from random partition splitting to virtual nodes (vnodes)?
