---
id: ch05-monotonic-reads
day: 6
tags: [replication, consistency, monotonic-reads, time-travel]
related_stories: []
---

# Monotonic Reads

## question
A user sees a comment on a post, refreshes the page, and the comment disappears. Later refresh brings it back. Which consistency guarantee would prevent this?

## options
- A) Strong consistency
- B) Read-after-write consistency
- C) Monotonic reads
- D) Consistent prefix reads

## answer
C

## explanation
Monotonic reads guarantee prevents "going back in time" - once a user has seen a particular state of the data, they won't see an earlier state on subsequent reads. This typically happens when reads are served by different replicas with varying replication lag. The solution is to ensure a user always reads from the same replica or from replicas with guaranteed monotonic timestamps.

## hook
Why might forcing users to read from the same replica create new problems?
