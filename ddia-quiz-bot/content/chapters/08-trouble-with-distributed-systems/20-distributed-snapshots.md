---
id: ch08-distributed-snapshots
day: 20
level: L5
tags: [snapshots, consistency, distributed-state]
related_stories: []
---

# Distributed Snapshots

## question
Why is taking a consistent snapshot of a distributed system challenging compared to a single-node system?

## options
- A) Distributed systems have more data
- B) There's no single point in time when all nodes can simultaneously capture their state
- C) Network bandwidth limitations
- D) Storage capacity issues

## answer
B

## explanation
In a single-node system, you can atomically capture the entire state at one instant. In distributed systems, nodes can't coordinate to snapshot at exactly the same time due to clock skew and message delays. Messages in flight during the snapshot can cause inconsistencies - they might be included in neither the sender's nor receiver's snapshot, or included in both. Algorithms like Chandy-Lamport create consistent snapshots by carefully tracking messages crossing the snapshot boundary.

## hook
How do you photograph a moving distributed system without blur?
