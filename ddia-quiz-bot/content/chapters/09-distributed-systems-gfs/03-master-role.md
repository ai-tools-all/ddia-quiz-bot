---
id: ch09-master-role
day: 1
tags: [gfs, master-server, architecture]
related_stories: []
---

# GFS Master Server Role

## question
What is the primary responsibility of the GFS master server?

## options
- A) Storing all file data and serving client read/write requests
- B) Managing metadata, chunk locations, and coordinating system-wide activities
- C) Performing data compression and deduplication
- D) Executing MapReduce jobs on stored data

## answer
B

## explanation
The GFS master maintains all file system metadata including the namespace, access control information, mapping from files to chunks, and current locations of chunks. It also controls system-wide activities like chunk lease management, garbage collection, and chunk migration.

## hook
How can a single master server manage petabytes of data across thousands of machines?
