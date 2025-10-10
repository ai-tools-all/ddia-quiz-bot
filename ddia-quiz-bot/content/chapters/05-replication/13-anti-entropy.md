---
id: ch05-anti-entropy
day: 13
tags: [replication, anti-entropy, repair, consistency]
related_stories: []
---

# Anti-Entropy and Repair

## question
What is the purpose of anti-entropy processes in distributed databases?

## options
- A) To compress data and save storage space
- B) To detect and repair inconsistencies between replicas
- C) To prevent unauthorized access to data
- D) To optimize query performance

## answer
B

## explanation
Anti-entropy processes continuously or periodically compare data between replicas and repair any inconsistencies found. This includes mechanisms like read repair (fixing inconsistencies detected during reads) and active anti-entropy (background processes using Merkle trees to efficiently compare and sync replicas). These processes ensure eventual consistency even when nodes miss updates.

## hook
How do Merkle trees make anti-entropy repair more efficient than comparing every record?
