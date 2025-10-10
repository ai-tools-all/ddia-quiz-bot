---
id: ch08-clock-skew-problems
day: 11
level: L5
tags: [clocks, clock-skew, ordering]
related_stories: []
---

# Clock Skew Problems

## question
Two nodes with 100ms clock skew write to a last-write-wins register using wall-clock timestamps. Node A writes at its time 12:00:00.500, then Node B writes at its time 12:00:00.400. What happens?

## options
- A) Node A's write wins because it happened first
- B) Node B's write wins because its timestamp is earlier
- C) Both writes are kept
- D) The system detects the conflict

## answer
B

## explanation
In last-write-wins using wall-clock timestamps, the write with the later timestamp wins, regardless of actual causality. Due to clock skew, Node B's clock is 100ms behind, so even though it writes after Node A, its timestamp (12:00:00.400) is earlier. Node A's write (12:00:00.500) will win, potentially losing Node B's update. This demonstrates why wall-clock timestamps are dangerous for ordering events in distributed systems.

## hook
Can a write from the past overwrite a write from the future?
