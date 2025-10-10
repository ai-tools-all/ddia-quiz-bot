---
id: ch09-chunk-size
day: 1
tags: [gfs, chunks, storage-design]
related_stories: []
---

# GFS Chunk Size

## question
What is the standard chunk size in GFS and why was this size chosen?

## options
- A) 4KB - to match typical page size
- B) 1MB - for efficient network transfer
- C) 64MB - to reduce metadata overhead and network overhead
- D) 1GB - to minimize chunk server interactions

## answer
C

## explanation
GFS uses 64MB chunks, which is much larger than typical file system block sizes. This large chunk size reduces the amount of metadata the master needs to store, reduces network overhead for large sequential reads/writes, and reduces the number of client-master interactions.

## hook
Why would using 64MB blocks instead of 4KB blocks make a distributed file system more efficient?
