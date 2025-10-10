---
id: ch09-garbage-collection
day: 1
tags: [gfs, garbage-collection, storage-management]
related_stories: []
---

# GFS Garbage Collection

## question
How does GFS handle deleted files and orphaned chunks?

## options
- A) Immediate deletion upon request
- B) Lazy garbage collection during regular scans
- C) Manual cleanup by administrators
- D) Never deletes data (append-only system)

## answer
B

## explanation
GFS uses lazy garbage collection. Deleted files are renamed to hidden names and removed after a few days. The master regularly scans and identifies orphaned chunks (not referenced by any file) and instructs chunk servers to delete them.

## hook
Why might "delete" not actually delete your data right away in a distributed system?
