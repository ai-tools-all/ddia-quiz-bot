---
id: ch06-partitioning-purpose
day: 1
tags: [partitioning, scalability, distributed-systems]
related_stories: []
---

# Purpose of Partitioning

## question
What is the primary purpose of partitioning (sharding) data in distributed systems?

## options
- A) To reduce storage costs by compressing data
- B) To scale beyond the capacity limits of a single machine
- C) To improve data security through isolation
- D) To simplify database queries

## answer
B

## explanation
Partitioning (also known as sharding) is primarily used to achieve scalability by distributing data across multiple machines. When datasets grow too large for a single machine's disk, memory, or processing power, partitioning allows you to spread the load across multiple nodes, with each partition handling a subset of the total data.

## hook
How do you handle a dataset when it no longer fits on your biggest server?
