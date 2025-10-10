---
id: ch05-failover-challenges
day: 17
tags: [replication, failover, split-brain, high-availability]
related_stories: []
---

# Failover Challenges

## question
During an automatic failover from a failed leader to a follower, what is the "split-brain" problem?

## options
- A) The database is split into multiple partitions
- B) Two nodes both believe they are the leader, accepting writes independently
- C) The network is partitioned into two segments
- D) Read and write operations are split between nodes

## answer
B

## explanation
Split-brain occurs when the original leader is temporarily unreachable (e.g., network partition) but not actually dead. The system promotes a follower to leader, but then the original leader comes back, resulting in two active leaders accepting writes. This causes data divergence and conflicts. Solutions include fencing tokens, STONITH (Shoot The Other Node In The Head), or requiring a majority quorum for leader election.

## hook
How does a generation number or epoch help prevent split-brain scenarios?
