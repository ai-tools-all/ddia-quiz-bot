---
id: ch07-event-sourcing-transactions
day: 20
tags: [transactions, event-sourcing, consistency, append-only]
related_stories: []
---

# Event Sourcing and Transactions

## question
How does event sourcing naturally avoid many traditional transaction conflicts?

## options
- A) Events are processed faster than regular updates
- B) Append-only writes don't conflict like updates do
- C) Events are automatically serialized by timestamp
- D) Event stores don't support concurrent access

## answer
B

## explanation
Event sourcing stores state changes as an append-only log of events rather than updating records in place. Since events are only appended (never modified), there are no update conflicts. Multiple events can be written concurrently without interference. Consistency is achieved by ordered event replay. This eliminates lost updates and many concurrency issues inherent in update-in-place systems.

## hook
How does Kafka's append-only log enable high-throughput transactional messaging?
