---
id: ch09-consistency-trade-offs-l6
day: 1
tags: [gfs, consistency, trade-offs, L6]
related_stories: []
---

# GFS Consistency Trade-offs (L6)

## question
Why did GFS choose a relaxed consistency model instead of strong consistency? What system design benefit did this provide?

## options
- A) Reduced storage costs by allowing different replica versions
- B) Simplified failure recovery and improved availability and performance
- C) Eliminated the need for primary replicas
- D) Made the system easier to understand for developers

## answer
B

## explanation
GFS's relaxed consistency model allows the system to continue operating during failures and network partitions without complex coordination protocols. This improves availability and performance at scale. Applications must handle inconsistencies, but this trade-off was acceptable for Google's workloads like MapReduce that could tolerate or work around weak consistency.

## hook
Would you choose a "sometimes incorrect but always available" system over an "always correct but sometimes unavailable" one?
