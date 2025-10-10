---
id: ch07-serial-execution
day: 13
tags: [transactions, serial-execution, performance, stored-procedures]
related_stories: []
---

# Actual Serial Execution

## question
Some modern databases (like Redis, VoltDB) achieve serializability by actually executing transactions serially on a single thread. What enables this approach to be practical?

## options
- A) Faster network connections
- B) RAM becoming cheap enough to keep all data in memory
- C) Better query optimizers
- D) Multi-core processors

## answer
B

## explanation
Serial execution becomes practical when the entire dataset fits in RAM, eliminating disk I/O waits that would make single-threaded execution too slow. Combined with stored procedures (reducing network round trips) and partitioning (allowing parallel execution across partitions), serial execution can achieve high throughput while guaranteeing serializability without locks or coordination overhead.

## hook
How does Redis achieve 100K+ operations per second with single-threaded execution?
