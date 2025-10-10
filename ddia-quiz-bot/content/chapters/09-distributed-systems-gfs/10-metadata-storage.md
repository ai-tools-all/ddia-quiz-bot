---
id: ch09-metadata-storage
day: 1
tags: [gfs, metadata, master-server]
related_stories: []
---

# GFS Metadata Storage

## question
Where does the GFS master store its metadata, and why?

## options
- A) On disk only for persistence
- B) In memory for performance, with operation log on disk
- C) Distributed across chunk servers
- D) In a separate database system

## answer
B

## explanation
The GFS master keeps all metadata in memory for fast access and operations. It persists metadata changes through an operation log on disk and periodic checkpoints. This design enables quick metadata operations while maintaining durability.

## hook
How can a single server's RAM hold the metadata for petabytes of storage?
