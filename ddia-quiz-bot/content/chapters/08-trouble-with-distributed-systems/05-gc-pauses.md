---
id: ch08-gc-pauses
day: 5
tags: [garbage-collection, pauses, jvm]
related_stories: []
---

# Garbage Collection Pauses

## question
What is the primary danger of garbage collection pauses in distributed systems with leader election?

## options
- A) Memory leaks
- B) Increased CPU usage
- C) A node may be declared dead while it's just paused, causing split-brain
- D) Slower application performance

## answer
C

## explanation
During a stop-the-world GC pause, a node cannot respond to any requests. If the pause exceeds the timeout threshold, other nodes may declare it dead and elect a new leader. When the paused node resumes, it may still think it's the leader, resulting in split-brain where two nodes believe they're the leader simultaneously. This can cause data corruption or inconsistency.

## hook
What happens when a leader node experiences a 30-second GC pause?
