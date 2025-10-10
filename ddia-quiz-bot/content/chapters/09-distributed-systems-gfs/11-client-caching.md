---
id: ch09-client-caching
day: 1
tags: [gfs, caching, client-design]
related_stories: []
---

# Client Caching in GFS

## question
What does the GFS client cache, and what does it NOT cache?

## options
- A) Caches both file data and metadata
- B) Caches file data but not metadata
- C) Caches metadata but not file data
- D) Doesn't cache anything

## answer
C

## explanation
GFS clients cache metadata (like chunk locations) to reduce master load, but they do not cache file data. This is because most applications stream through huge files or have working sets too large to benefit from caching, and it simplifies coherence.

## hook
Why would a file system choose NOT to cache the actual file data?
