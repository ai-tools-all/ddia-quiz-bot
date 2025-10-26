---
id: spanner-snapshot-reads
day: 3
tags: [spanner, snapshot, reads, timestamps]
---

# Consistent Snapshot Reads

## question
How do read-only transactions in Spanner ensure a consistent view across shards?

## options
- A) They acquire read locks on all keys involved
- B) They read from leaders only to get the very latest values
- C) They read data as of a chosen timestamp and replicas delay until up-to-date through that timestamp
- D) They execute on a single shard to avoid cross-shard issues

## answer
C

## explanation
Read-only transactions are assigned a timestamp and read the most recent committed versions with timestamps ≤ that value. Replicas track freshness and may wait until they have applied all updates through the requested timestamp, enabling low-latency local reads with a consistent snapshot.

## hook
How can replicas serve consistent reads locally without contacting remote data centers?
