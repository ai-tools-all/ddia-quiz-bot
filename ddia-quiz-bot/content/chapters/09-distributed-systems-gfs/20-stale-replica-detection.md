---
id: ch09-stale-replica-detection
day: 1
tags: [gfs, stale-replicas, fault-tolerance]
related_stories: []
---

# Detecting Stale Replicas

## question
How does a GFS client avoid reading stale data from an outdated chunk replica?

## options
- A) Always reads from the primary replica
- B) Checks chunk version number before reading
- C) Uses timestamps to detect old data
- D) Relies on eventual consistency

## answer
B

## explanation
Clients and chunk servers verify chunk version numbers during operations. The master includes chunk version numbers in its replies to clients. If a chunk server has a stale version (was offline during a write), reads from it will fail.

## hook
How do you know if the data you're reading is the latest version or yesterday's news?
