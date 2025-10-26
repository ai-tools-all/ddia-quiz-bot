---
id: spanner-architecture-replication
day: 1
tags: [spanner, paxos, replication, shards]
---

# Spanner Replication Architecture

## question
In Spanner, how are writes to a shard replicated across data centers?

## options
- A) Each replica accepts writes independently and resolves conflicts with last-write-wins
- B) Writes go through the shard’s Paxos leader and are replicated to a majority of replicas
- C) Writes are broadcast to all replicas and must be acknowledged by all to commit
- D) Writes are sent only to the nearest replica for low latency

## answer
B

## explanation
Each shard is a Paxos group with a single leader handling writes. The leader replicates log entries to followers and commits once a majority acknowledges, providing fault tolerance and progress despite a minority of slow or failed replicas.

## hook
Why does a single Paxos leader per shard help Spanner make progress during failures?
