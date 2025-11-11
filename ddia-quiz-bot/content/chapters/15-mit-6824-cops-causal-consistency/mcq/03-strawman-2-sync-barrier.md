---
id: cops-strawman-2-sync-barrier
day: 1
tags: [cops, sync-barrier, design-alternatives, latency]
---

# Strawman 2: Explicit Synchronization

## question
Why does COPS reject Strawman 2 (explicit sync barriers) as a solution to causal consistency?

## options
- A) Sync barriers require complex vector clock implementations that are hard to maintain
- B) Sync barriers force writes to wait for cross-datacenter propagation, sacrificing the goal of local write latency
- C) Sync barriers don't actually solve the causal ordering problem
- D) Sync barriers increase storage overhead by requiring additional metadata

## answer
B

## explanation
Strawman 2 introduces a sync barrier (similar to fsync) that forces a write to wait until it propagates to all data centers before returning. While this solves the ordering problem, it sacrifices the goal of local writes—clients must wait for cross-datacenter round-trips, negating the performance benefits and reducing fault tolerance. Systems like Spanner and Facebook's Memcache already incur similar write latencies, so this approach offers no improvement over existing solutions.

## hook
What's the fundamental trade-off between local writes and strong consistency guarantees?
