---
id: ch05-single-leader-replication
day: 2
tags: [replication, single-leader, master-slave, consistency]
related_stories: []
---

# Single-Leader Replication

## question
In single-leader (master-slave) replication, how are writes typically handled?

## options
- A) Writes can go to any replica and are synchronized later
- B) Writes go only to the leader, which then replicates to followers
- C) Writes are sent to all replicas simultaneously
- D) Writes are load-balanced across all replicas

## answer
B

## explanation
In single-leader replication, all writes must go through the single leader node. The leader processes the write and then sends the data change to all follower replicas. Clients can read from any replica (leader or follower), but writes must go through the leader to maintain consistency.

## hook
What happens to writes when the single leader fails?
