---
id: cops-lamport-clock-ordering
day: 1
tags: [cops, lamport-clock, version-ordering, concurrent-writes]
---

# Lamport Clock Version Ordering

## question
Two clients at different data centers concurrently update key X. DC1 assigns version 42, DC2 assigns version 41. Both updates have no causal relationship. What happens when DC1 receives DC2's update?

## options
- A) DC1 rejects DC2's update because 41 < 42
- B) DC1 keeps both versions and returns them as siblings
- C) DC1 applies DC2's update only if it has later physical timestamp
- D) DC1 discards DC2's update because 41 < 42 (last-writer-wins)

## answer
D

## explanation
COPS uses Lamport clocks for version numbers and applies last-writer-wins (LWW) for concurrent updates to the same key. When DC1 (version 42) receives DC2's update (version 41), it compares version numbers and discards the lower version. This means DC2's update is permanently lost. LWW is simple but can lose data—applications needing better conflict resolution must use CRDTs or custom merge logic on top of COPS.

## hook
What are the practical implications of losing updates in a shopping cart or collaborative editing scenario?
