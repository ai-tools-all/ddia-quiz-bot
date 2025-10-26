---
id: spanner-external-consistency
day: 6
tags: [spanner, external-consistency, ordering]
---

# External Consistency

## question
Which guarantee best describes Spanner’s external consistency?

## options
- A) Transactions appear to execute in some serial order unrelated to real time
- B) If T1 commits before T2 starts in real time, then T2’s timestamp is greater and T2 sees T1’s effects
- C) Each shard is consistent, but cross-shard reads may see anomalies
- D) Only read-only transactions are externally consistent

## answer
B

## explanation
Using the start and commit-wait rules with TrueTime, Spanner ensures that real-time ordering is respected in the serialization order. Later-starting transactions receive later timestamps and observe earlier commits.

## hook
How does timestamp assignment align the serial order with wall-clock order?
