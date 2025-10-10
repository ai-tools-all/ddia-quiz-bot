---
id: ch06-consistent-hashing
day: 11
tags: [partitioning, consistent-hashing, load-balancing]
related_stories: []
---

# Consistent Hashing

## question
What problem does consistent hashing primarily solve in distributed systems?

## options
- A) Ensuring ACID transactions across partitions
- B) Minimizing data movement when nodes are added/removed
- C) Providing strong consistency guarantees
- D) Eliminating network partitions

## answer
B

## explanation
Consistent hashing minimizes the amount of data that needs to be moved when nodes are added or removed from the cluster. Instead of rehashing all keys (which would happen with simple modulo hashing), consistent hashing ensures that only ~1/n of the keys need to be relocated when adding/removing a node in an n-node cluster. This makes cluster resizing operations much more efficient.

## hook
Why do CDNs (Content Delivery Networks) heavily rely on consistent hashing?
