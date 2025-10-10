---
id: ch05-quorum-conditions
day: 11
tags: [replication, quorums, consistency, availability]
related_stories: []
---

# Quorum Read and Write Conditions

## question
In a system with n=5 replicas, w=3 (write quorum), what is the minimum read quorum (r) to guarantee you'll read the latest value?

## options
- A) r = 1
- B) r = 2
- C) r = 3
- D) r = 5

## answer
C

## explanation
To guarantee reading the latest value, the condition w + r > n must be satisfied. With n=5 and w=3, we need r > 2, so minimum r = 3. This ensures that the read quorum and write quorum overlap on at least one node that has the latest value. With r=3 and w=3, there's guaranteed to be at least one node in common (3+3 > 5).

## hook
What happens to availability when you increase quorum requirements?
