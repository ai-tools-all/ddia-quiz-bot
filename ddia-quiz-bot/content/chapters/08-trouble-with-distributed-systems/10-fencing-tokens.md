---
id: ch08-fencing-tokens
day: 10
level: L5
tags: [fencing, distributed-locking, safety]
related_stories: []
---

# Fencing Tokens

## question
How do fencing tokens prevent split-brain in distributed locking systems?

## options
- A) They prevent nodes from acquiring locks
- B) They ensure only requests with the latest token are accepted by the storage system
- C) They make locks permanent
- D) They speed up lock acquisition

## answer
B

## explanation
Fencing tokens are monotonically increasing numbers issued with each lock acquisition. When a client acquires a lock, it gets token n. The storage system only accepts writes with the highest token seen so far. If a node with an old lock (token n-1) tries to write after a pause, its write is rejected because the storage has seen token n. This prevents a delayed node from corrupting data after its lock has been reassigned.

## hook
What happens when a "zombie" node wakes up thinking it still holds a lock?
