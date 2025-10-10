---
id: ch09-hot-spots
day: 1
tags: [gfs, hot-spots, performance, load-balancing]
related_stories: []
---

# Handling Hot Spots in GFS

## question
What problem can arise with small files in GFS due to the large chunk size?

## options
- A) Increased metadata overhead
- B) Hot spots when many clients access the same small file
- C) Wasted storage space
- D) Slower read performance

## answer
B

## explanation
Small files may consist of only one or a few chunks. If many clients access the same small file simultaneously, the chunk servers storing those chunks become hot spots. This is a trade-off of the large chunk size design decision.

## hook
What happens when a million users try to download the same 1KB configuration file?
