---
id: ch05-consistent-prefix-reads
day: 7
tags: [replication, consistency, causal-consistency, ordering]
related_stories: []
---

# Consistent Prefix Reads

## question
In a chat application, User A asks "What's the capital of France?" and User B replies "Paris". Some users see the answer before the question. What consistency guarantee is violated?

## options
- A) Linearizability
- B) Read-after-write consistency
- C) Consistent prefix reads
- D) Monotonic writes

## answer
C

## explanation
Consistent prefix reads guarantee that if a sequence of writes happens in a certain order, anyone reading those writes will see them in the same order. This is particularly important in partitioned databases where different partitions may have different replication lag. Without this guarantee, causally related events can appear in the wrong order.

## hook
How do version vectors help maintain consistent prefix reads?
