---
id: ch08-network-congestion
day: 17
tags: [network-congestion, queueing, performance]
related_stories: []
---

# Network Congestion and Queueing

## question
Where do network delays primarily come from in datacenter networks?

## options
- A) Speed of light limitations
- B) Queueing at switches, routers, and network interfaces
- C) Cable length
- D) Packet encryption/decryption

## answer
B

## explanation
In datacenter networks, propagation delay (speed of light) is negligible due to short distances. Most delays come from queueing: at network switches when multiple packets compete for the same output port, at the sender's network interface when the link is busy, and at the receiver's network buffer. During congestion, packets can spend far more time queued than actually traveling. This variability in queueing delay makes network timing unpredictable.

## hook
Why can a packet take 1ms or 100ms to cross the same datacenter network?
