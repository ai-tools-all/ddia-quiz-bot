---
id: ch09-primary-replica
day: 1
tags: [gfs, primary-replica, consistency]
related_stories: []
---

# Primary Replica in GFS

## question
What is the role of the primary replica in GFS write operations?

## options
- A) It stores the only copy of the data
- B) It serializes all mutations and assigns order to concurrent writes
- C) It performs data compression before storage
- D) It handles only read requests while secondaries handle writes

## answer
B

## explanation
The primary replica is designated by the master for each chunk and holds a lease. It determines the order of all mutations (writes) to that chunk and forwards them to secondary replicas in the same order, ensuring consistency across all replicas.

## hook
When multiple clients write to the same file simultaneously, who decides which write happens first?
