---
id: craq-read-throughput
day: 13
tags: [craq, chain-replication, consistency]
related_stories:
  - vmware-ft
---

# CRAQ Read Scaling

## question
What performance advantage does CRAQ add on top of chain replication while still retaining linearizability?

## options
- A) Allowing multiple heads to accept writes concurrently without extra coordination
- B) Enabling clients to read from any replica without giving up linearizable semantics
- C) Eliminating the need for an external configuration manager to coordinate replicas

## answer
B

## explanation
The lecture notes that CRAQ optimizes chain replication by letting reads hit any replica, delivering an N-way throughput boost while preserving the strong ordering guarantees that chain replication already provides.

## hook
How does CRAQ ensure replicas only serve up-to-date data even while permitting distributed reads?
