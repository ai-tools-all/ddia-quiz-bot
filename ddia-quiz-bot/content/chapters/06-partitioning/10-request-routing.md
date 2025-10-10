---
id: ch06-request-routing
day: 10
tags: [partitioning, request-routing, service-discovery]
related_stories: []
---

# Request Routing

## question
Which request routing approach requires all nodes to maintain the complete partition-to-node mapping?

## options
- A) Routing through a central routing tier
- B) Client-side partition awareness with gossip protocol
- C) External coordination service (like ZooKeeper)
- D) Random node contact with forwarding

## answer
B

## explanation
In gossip protocol-based routing (used by Cassandra and Riak), all nodes maintain the complete cluster state including partition-to-node mappings. Nodes share state changes through gossip protocol, allowing any node to route requests correctly. While this adds complexity and eventual consistency to cluster state, it eliminates single points of failure in routing.

## hook
How does Redis Cluster achieve request routing without a separate coordination service?
