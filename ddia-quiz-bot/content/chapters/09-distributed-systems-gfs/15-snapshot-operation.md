---
id: ch09-snapshot-operation
day: 1
tags: [gfs, snapshot, copy-on-write]
related_stories: []
---

# GFS Snapshot Operation

## question
How does GFS implement file and directory tree snapshots efficiently?

## options
- A) By immediately copying all data to new locations
- B) Using copy-on-write at the chunk level
- C) By maintaining multiple versions of all chunks
- D) Through incremental backup to tape

## answer
B

## explanation
GFS implements snapshots using copy-on-write. When a snapshot is created, the master revokes outstanding leases, logs the operation, and duplicates metadata. Actual chunk copying only happens when chunks are modified after the snapshot, making it space and time efficient.

## hook
How can you create an instant copy of a petabyte of data without actually copying anything?
