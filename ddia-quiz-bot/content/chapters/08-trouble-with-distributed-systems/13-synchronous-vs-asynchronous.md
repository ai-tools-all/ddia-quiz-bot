---
id: ch08-sync-vs-async-networks
day: 13
level: L5
tags: [network-models, synchronous, asynchronous]
related_stories: []
---

# Synchronous vs Asynchronous Network Models

## question
Why do most practical distributed systems assume an asynchronous network model rather than a synchronous one?

## options
- A) Asynchronous networks are faster
- B) Real networks have unbounded delays and no guaranteed delivery time
- C) Synchronous networks don't exist
- D) Asynchronous models are simpler to implement

## answer
B

## explanation
Real networks (internet, datacenter networks) are asynchronous - they provide no guarantees about maximum message delay or even delivery. While synchronous models (with bounded delays) would make distributed algorithms simpler and more powerful, they don't match reality. Systems must handle unbounded delays, message loss, and reordering. This is why consensus algorithms like Raft and Paxos are designed for asynchronous networks with failure detection via timeouts.

## hook
Why can't we just assume messages arrive within 1 second and build simpler systems?
