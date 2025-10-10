---
id: ch06-geo-partitioning
day: 20
tags: [partitioning, geo-distribution, compliance, latency]
related_stories: []
---

# Geographic Partitioning

## question
A global service partitions data by geographic region. Besides compliance requirements, what is another significant benefit?

## options
- A) Reduced storage costs
- B) Simplified backup procedures
- C) Lower latency for regional queries
- D) Elimination of network partitions

## answer
C

## explanation
Geographic partitioning places data close to where it's accessed, significantly reducing latency for queries that are naturally regional (e.g., users typically access their regional data). This locality of reference improves performance and user experience. While compliance (data sovereignty) is often the primary driver, the latency benefits are substantial for geographically distributed user bases.

## hook
How does Netflix optimize content delivery by combining geographic partitioning with CDN strategies?
