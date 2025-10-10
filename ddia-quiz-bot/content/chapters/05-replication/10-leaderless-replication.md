---
id: ch05-leaderless-replication
day: 10
tags: [replication, leaderless, dynamo, cassandra, quorums]
related_stories: []
---

# Leaderless Replication

## question
In a leaderless replication system like Cassandra, how are writes typically handled?

## options
- A) Writes go to a randomly selected node that coordinates replication
- B) Writes are sent to multiple nodes in parallel, success requires a quorum
- C) Writes are queued until all nodes are available
- D) Writes rotate through nodes in round-robin fashion

## answer
B

## explanation
In leaderless replication, the client sends writes to multiple replicas in parallel. A write is considered successful when a quorum (e.g., majority) of nodes acknowledge it. Similarly, reads query multiple nodes and return the most recent value based on version numbers. This approach avoids single points of failure but requires careful quorum configuration.

## hook
What happens when a node comes back online after missing some writes?
