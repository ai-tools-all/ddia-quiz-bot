---
id: network-multi-hop-routing
day: 8
tags: [networking, routing, mac-address, ip-address, ttl]
related_stories:
  - network-basics
---

# Multi-Hop Routing

## question
When your laptop (192.168.1.10, MAC: AA:BB:CC:DD:EE:01) sends a packet to Google (8.8.8.8) that passes through 3 routers before reaching the destination, which statement is correct about the packet headers?

## options
- A) The source MAC remains AA:BB:CC:DD:EE:01 and source IP remains 192.168.1.10 throughout the entire journey
- B) The source IP stays 192.168.1.10 but the source MAC changes at each router hop, and TTL decrements
- C) Both source MAC and source IP remain unchanged, only the destination addresses change
- D) The source IP changes at each router due to NAT, but the source MAC remains constant

## answer
B

## explanation
As the packet traverses routers: (1) The Layer 3 IP addresses remain constant end-to-end (source: 192.168.1.10, destination: 8.8.8.8) unless NAT is involved at the home router. (2) The Layer 2 MAC addresses change at every hop because each link is a separate broadcast domain - the source MAC becomes the outgoing router interface's MAC, and the destination MAC becomes the next hop's MAC. (3) TTL (Time To Live) decrements at each router to prevent routing loops. This illustrates why we need both Layer 2 (local) and Layer 3 (end-to-end) addressing.

## hook
What would happen if TTL reached zero before the packet reached its destination?
