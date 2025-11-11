---
id: memcached-regional-architecture
day: 1
tags: [memcached, replication, multi-datacenter, consistency]
---

# Regional Replication Architecture

## question
How does Facebook's regional replication architecture balance latency and consistency?

## options
- A) All reads and writes are synchronously replicated across regions for strong consistency
- B) Reads are local within each region while all writes flow to the primary region's MySQL master, accepting seconds of replication lag
- C) Each region independently handles reads and writes with eventual reconciliation
- D) Writes are load-balanced across regional MySQL masters for better write throughput

## answer
B

## explanation
Facebook deploys multiple regions (e.g., West Coast primary, East Coast secondary), each with complete data replicas. All writes flow to the primary region's MySQL master to ensure a single source of truth. Asynchronous MySQL replication propagates updates to secondary regions with seconds of lag. Reads remain local within each region (from both memcached and MySQL), exploiting the read-heavy workload while tolerating brief cross-region inconsistency. This design prioritizes low-latency local reads over global consistency.

## hook
Why is it acceptable for users far from the primary region to see slightly stale data in a social media context?
