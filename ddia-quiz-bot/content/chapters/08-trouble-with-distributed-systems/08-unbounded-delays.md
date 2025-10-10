---
id: ch08-unbounded-delays
day: 8
tags: [network-delays, asynchronous-systems, distributed-systems]
related_stories: []
---

# Unbounded Network Delays

## question
Why are network delays considered "unbounded" in distributed systems?

## options
- A) Because networks are infinitely fast
- B) Because there's no guaranteed maximum time for message delivery
- C) Because delays are always the same
- D) Because networks don't have bandwidth limits

## answer
B

## explanation
In asynchronous networks (like the internet and most datacenters), there's no upper bound on how long a message might take to arrive. A packet could be delayed indefinitely due to network congestion, queuing, retransmissions, or routing issues. This unbounded delay makes it impossible to distinguish between a very slow message and a lost message, which is why timeouts are probabilistic, not deterministic.

## hook
Can you guarantee a message will arrive within 1 second in a distributed system?
