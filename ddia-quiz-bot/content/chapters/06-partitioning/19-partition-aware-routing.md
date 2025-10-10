---
id: ch06-partition-aware-routing
day: 19
tags: [partitioning, client-optimization, request-routing]
related_stories: []
---

# Partition-Aware Clients

## question
What is the main advantage of partition-aware clients over proxy-based routing?

## options
- A) Simpler client implementation
- B) Better security through isolation
- C) Lower latency by avoiding extra network hops
- D) Automatic failover handling

## answer
C

## explanation
Partition-aware clients maintain partition mappings and route requests directly to the correct node, eliminating the extra network hop through a proxy or routing tier. This reduces latency and removes a potential bottleneck. However, it requires more complex client logic and coordination to keep partition mappings updated across all clients.

## hook
Why do high-frequency trading systems always use partition-aware clients despite the added complexity?
