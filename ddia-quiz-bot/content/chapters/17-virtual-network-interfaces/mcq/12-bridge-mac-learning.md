---
id: virtual-network-bridge-mac-learning
day: 12
tags: [networking, bridge, mac-address, switching]
related_stories:
  - docker-networking
---

# Bridge MAC Learning

## question
How does a Linux bridge (like docker0) know which interface to forward packets to?

## options
- A) It broadcasts every packet to all connected interfaces
- B) It maintains a MAC address learning table mapping MAC addresses to interfaces
- C) It uses DNS to resolve container names to interfaces
- D) It relies on iptables rules to determine the destination

## answer
B

## explanation
A Linux bridge behaves like a Layer 2 switch, learning which MAC addresses are reachable through which interfaces. When a packet arrives, the bridge examines the source MAC and records which interface it came from. For forwarding, it looks up the destination MAC in its learning table.

## hook
What happens when the bridge receives a packet for a MAC address it hasn't learned yet?
