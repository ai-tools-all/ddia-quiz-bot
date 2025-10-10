---
id: ch08-network-partition
day: 6
tags: [network-partition, split-brain, distributed-systems]
related_stories: []
---

# Network Partitions

## question
During a network partition, what problem occurs when both sides of the partition continue accepting writes?

## options
- A) Performance degradation
- B) Split-brain leading to divergent data
- C) Increased latency
- D) Memory overflow

## answer
B

## explanation
When a network partition occurs, nodes on each side may continue operating independently. If both sides accept writes, they will have divergent data states (split-brain). When the partition heals, these divergent states must be reconciled, which may require complex conflict resolution or could result in data loss. This is why many systems choose to sacrifice availability during partitions to maintain consistency.

## hook
What happens when half your database cluster can't talk to the other half?
