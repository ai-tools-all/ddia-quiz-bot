---
id: ch05-multi-leader-replication
day: 8
tags: [replication, multi-leader, multi-master, conflicts]
related_stories: []
---

# Multi-Leader Replication Use Cases

## question
In which scenario is multi-leader replication most beneficial?

## options
- A) Single datacenter with high read load
- B) Multiple datacenters with users in different geographic regions
- C) Applications requiring strong consistency
- D) Systems with infrequent writes

## answer
B

## explanation
Multi-leader replication shines in multi-datacenter deployments where each datacenter has its own leader. This reduces write latency (users write to their nearest datacenter), provides better fault tolerance (each datacenter can operate independently), and reduces expensive cross-datacenter network traffic. The trade-off is dealing with write conflicts.

## hook
What happens when the same user updates their profile from two different continents simultaneously?
