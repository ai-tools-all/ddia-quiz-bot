---
id: cops-causal-consistency-guarantees
day: 1
tags: [cops, causal-consistency, guarantees, transitivity]
---

# Causal Consistency Guarantees

## question
What guarantee does COPS provide through its causal consistency model?

## options
- A) All clients see all writes in the exact same global order (linearizability)
- B) If operation A causally precedes B, then all clients observe A before B
- C) Writes are immediately visible at all replicas within 100ms
- D) Conflicting writes are automatically merged using CRDTs

## answer
B

## explanation
COPS enforces causal consistency: if operation A causally precedes B (through same-client sequencing or cross-client read-write chains), then all clients observe A before B. Dependencies are transitive—a write depending on a read of an earlier write inherits all of that earlier write's dependencies, forming causal chains across multiple clients and keys. This sits between eventual consistency (too weak) and linearizability (too strong), offering a practical middle ground for geo-replication.

## hook
How does causal consistency differ from eventual consistency and linearizability?
