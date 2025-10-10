---
id: ch09-concurrent-mutations
day: 1
tags: [gfs, concurrency, mutations, consistency]
related_stories: []
---

# Concurrent Mutations in GFS

## question
What happens when multiple clients concurrently write to the same region of a GFS file?

## options
- A) One write succeeds, others fail with errors
- B) Writes are serialized by locking
- C) The region becomes undefined with mixed fragments
- D) The system crashes to prevent corruption

## answer
C

## explanation
When concurrent writes (not record appends) target the same region, GFS doesn't guarantee serialization. The region becomes undefined - it may contain fragments from multiple writes mixed together. Applications must use record append or their own locking for consistency.

## hook
What happens when two people try to edit the same paragraph at exactly the same time?
