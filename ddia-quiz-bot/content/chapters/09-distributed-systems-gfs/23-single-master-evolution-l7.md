---
id: ch09-single-master-evolution-l7
day: 1
tags: [gfs, architecture-evolution, distributed-systems, L7]
related_stories: []
---

# Evolution Beyond Single Master (L7)

## question
GFS's single-master design was eventually replaced by Colossus at Google. Which architectural change would best address GFS's scalability limitations while preserving its simplicity benefits?

## options
- A) Fully decentralized peer-to-peer architecture with eventual consistency
- B) Multiple masters with sharded namespace and distributed metadata management
- C) Blockchain-based consensus for all metadata operations
- D) Client-managed metadata with no centralized coordination

## answer
B

## explanation
Multiple masters with sharded namespace (used in Colossus and similar to HDFS Federation) maintains the simplicity of centralized coordination within each shard while scaling horizontally. This preserves GFS's benefits (simple consistency, easy debugging) while removing the single master bottleneck. Each master handles a portion of the namespace independently.

## hook
How would you evolve a system that worked perfectly for 1PB but struggles at 1000PB?
