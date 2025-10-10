---
id: ch07-transaction-overhead
day: 15
tags: [transactions, performance, overhead, trade-offs]
related_stories: []
---

# Transaction Overhead

## question
An e-commerce site's product view counter doesn't use transactions, accepting occasional inaccuracy. What is this design pattern called?

## options
- A) Eventually consistent design
- B) Best-effort delivery
- C) Sacrificing consistency for performance
- D) Event sourcing

## answer
C

## explanation
This is deliberately sacrificing consistency for performance. For non-critical data like view counts, the overhead of transactions (locking, coordination, potential retries) isn't justified. Accepting occasional lost updates or slight inaccuracies is a valid trade-off when the business impact is minimal but the performance gain is significant. This pattern is common for metrics, analytics, and approximate counts.

## hook
What types of data in YouTube or Netflix can tolerate inconsistency vs requiring strict transactions?
