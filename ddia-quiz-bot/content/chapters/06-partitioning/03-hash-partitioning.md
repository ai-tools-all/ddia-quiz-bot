---
id: ch06-hash-partitioning
day: 3
tags: [partitioning, hash-partitioning, load-balancing]
related_stories: []
---

# Hash Partitioning

## question
What is the main trade-off when using hash partitioning instead of key range partitioning?

## options
- A) Higher storage costs but better performance
- B) Better load distribution but loss of efficient range queries
- C) Increased network traffic but reduced latency
- D) More complex setup but easier maintenance

## answer
B

## explanation
Hash partitioning uses a hash function to determine partition placement, which typically provides better load distribution by randomly scattering keys across partitions. However, this destroys the sort order of keys, making range queries inefficient as they must query all partitions. This is the classic trade-off: uniform load distribution versus range query capability.

## hook
How does Cassandra handle the trade-off between hash and range partitioning?
