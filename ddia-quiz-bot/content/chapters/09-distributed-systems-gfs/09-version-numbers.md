---
id: ch09-version-numbers
day: 1
tags: [gfs, version-numbers, consistency]
related_stories: []
---

# Version Numbers in GFS

## question
What is the primary purpose of version numbers for chunks in GFS?

## options
- A) To implement multi-version concurrency control
- B) To detect stale replicas after chunk server failures
- C) To track file modification history
- D) To implement snapshot isolation

## answer
B

## explanation
GFS uses version numbers to detect stale replicas. When the master grants a new lease, it increments the chunk version number. Chunk servers that were offline during this will have outdated version numbers, allowing the system to identify and garbage collect stale replicas.

## hook
How do you know if a chunk server that just came back online has outdated data?
