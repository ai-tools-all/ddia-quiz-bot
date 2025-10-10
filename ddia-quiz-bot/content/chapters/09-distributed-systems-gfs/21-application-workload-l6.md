---
id: ch09-application-workload-l6
day: 1
tags: [gfs, workload, system-design, L6]
related_stories: []
---

# GFS Application Workload (L6)

## question
GFS was optimized for specific workload patterns at Google. Which design trade-off would be LEAST appropriate for a system handling millions of small (< 1KB) files with random access patterns?

## options
- A) Large chunk size of 64MB
- B) Client-side metadata caching
- C) Append-only write optimization
- D) Three-way replication

## answer
A

## explanation
The 64MB chunk size would be highly inefficient for tiny files, creating severe hot spots and wasting space. Small files would typically fit in a single chunk, causing all requests to hit the same chunk servers. A system for small files would need smaller block sizes and different metadata management strategies.

## hook
If GFS is perfect for Google's web crawl data, what workload would make it perform terribly?
