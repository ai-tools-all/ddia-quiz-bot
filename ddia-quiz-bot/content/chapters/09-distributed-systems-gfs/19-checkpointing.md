---
id: ch09-checkpointing
day: 1
tags: [gfs, checkpointing, master-recovery]
related_stories: []
---

# Master Checkpointing

## question
Why does the GFS master create periodic checkpoints of its metadata?

## options
- A) To provide snapshots for users
- B) To reduce recovery time after master crashes
- C) To backup data to remote locations
- D) To compress metadata for storage efficiency

## answer
B

## explanation
The master periodically creates checkpoints of its metadata state. These checkpoints allow faster recovery after a crash - instead of replaying the entire operation log from the beginning, the master can load the latest checkpoint and only replay subsequent operations.

## hook
How do you quickly rebuild the phone book for millions of files after a crash?
