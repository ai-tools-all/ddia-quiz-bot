---
id: ch07-optimistic-vs-pessimistic
day: 18
tags: [transactions, concurrency-control, optimistic, pessimistic]
related_stories: []
---

# Optimistic vs Pessimistic Concurrency

## question
When would pessimistic locking be preferred over optimistic concurrency control?

## options
- A) Read-heavy workloads with rare writes
- B) High contention with frequent conflicts on the same data
- C) Geographically distributed systems
- D) Systems with many short transactions

## answer
B

## explanation
Pessimistic locking is preferred when conflict probability is high. With frequent conflicts, optimistic approaches waste work due to many transaction aborts and retries. Pessimistic locking prevents conflicts upfront by acquiring locks, avoiding wasted work. It's better when the cost of retrying (due to complex computations or side effects) exceeds the cost of waiting for locks.

## hook
Why do airline reservation systems still use pessimistic locking for seat selection?
