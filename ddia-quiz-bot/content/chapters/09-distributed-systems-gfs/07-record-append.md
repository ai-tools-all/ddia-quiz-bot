---
id: ch09-record-append
day: 1
tags: [gfs, record-append, atomic-operations]
related_stories: []
---

# GFS Record Append

## question
What guarantee does GFS provide for record append operations?

## options
- A) Exactly-once delivery with strict ordering
- B) At-least-once delivery, possibly with duplicates
- C) At-most-once delivery with possible data loss
- D) Best-effort delivery with no guarantees

## answer
B

## explanation
GFS record append guarantees at-least-once append semantics. The data will be appended atomically at least once, but may appear multiple times in case of failures and retries. Applications must be prepared to handle duplicate records.

## hook
What happens when your append operation succeeds but the success message gets lost?
