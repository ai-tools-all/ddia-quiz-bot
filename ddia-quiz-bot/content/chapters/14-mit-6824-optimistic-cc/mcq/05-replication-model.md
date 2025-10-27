---
id: occ-replication-model
day: 5
tags: [replication, primary-backup, fault-tolerance]
---

# Replication Model

## question
Which statement best describes replication and failure tolerance in this single–data center design?

## options
- A) It uses Paxos majority quorums across regions to commit writes
- B) It uses primary-backup within one data center; with F+1 replicas per shard to tolerate F failures, a single surviving replica can recover data after crashes
- C) It requires acknowledgments from all replicas to commit
- D) It cannot tolerate any server failure because data is only in volatile memory

## answer
B

## explanation
The system replicates within one data center using primary-backup pairs and non-volatile RAM. To tolerate F failures, maintain F+1 replicas. It does not rely on cross-region majorities.

## hook
How do durability and availability differ in single-DC primary-backup vs quorum systems?
