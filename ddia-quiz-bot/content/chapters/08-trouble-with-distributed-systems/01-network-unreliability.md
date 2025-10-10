---
id: ch08-network-unreliability
day: 1
tags: [network, unreliability, distributed-systems]
related_stories: []
---

# Network Unreliability

## question
What is the most fundamental assumption you should make about network communication in distributed systems?

## options
- A) Networks are reliable and packets always arrive
- B) Networks are unreliable and packets may be lost, delayed, or duplicated
- C) Network failures are rare and can be ignored
- D) Packet ordering is always preserved

## answer
B

## explanation
In distributed systems, networks are inherently unreliable. Packets may be lost due to network congestion, router failures, or cable issues. They may be delayed indefinitely or arrive multiple times. This is why protocols like TCP exist to provide reliability guarantees on top of unreliable networks.

## hook
Why can't distributed systems assume the network is reliable?
