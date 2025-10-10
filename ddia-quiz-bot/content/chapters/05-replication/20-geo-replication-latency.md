---
id: ch05-geo-replication-latency
day: 20
tags: [replication, geo-distribution, latency, multi-region]
related_stories: []
---

# Geographic Replication and Latency

## question
A company has users in US-West, Europe, and Asia. They use single-leader replication with the leader in US-West. European users complain about slow writes. What's the best solution?

## options
- A) Add more followers in Europe
- B) Increase network bandwidth between regions
- C) Switch to multi-leader replication with a leader in each region
- D) Cache more data in Europe

## answer
C

## explanation
With single-leader replication, all writes must go to US-West, causing high latency for European and Asian users (potentially 100-200ms round trip). Multi-leader replication places a leader in each region, allowing local writes with low latency. The trade-off is handling write conflicts between regions, but for many applications, the latency improvement justifies the added complexity.

## hook
How would you handle conflicts when users in different regions update the same data?
