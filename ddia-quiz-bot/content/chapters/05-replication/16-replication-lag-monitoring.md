---
id: ch05-replication-lag-monitoring
day: 16
tags: [replication, monitoring, operations, lag]
related_stories: []
---

# Monitoring Replication Lag

## question
Your monitoring shows a follower is 5 seconds behind the leader. What is the most concerning scenario?

## options
- A) The lag is steady at 5 seconds
- B) The lag is decreasing from 10 seconds
- C) The lag is increasing from 1 second
- D) The lag occasionally spikes to 5 seconds then returns to near-zero

## answer
C

## explanation
Increasing replication lag indicates the follower cannot keep up with the write rate and is falling further behind. This could lead to unbounded lag and eventual failure. Steady lag might be acceptable if within SLA. Decreasing lag shows recovery. Occasional spikes might be normal during batch operations. The trend matters more than the absolute value.

## hook
What metrics besides lag time would help diagnose replication problems?
