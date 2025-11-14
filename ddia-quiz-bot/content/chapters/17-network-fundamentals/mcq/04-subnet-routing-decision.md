---
id: network-subnet-routing
day: 4
tags: [networking, subnet, routing, subnet-mask]
related_stories:
  - network-basics
---

# Subnet Mask and Routing Decision

## question
Your PC has IP 192.168.1.10 with subnet mask 255.255.255.0 (/24). You want to send a packet to 192.168.2.50. What routing decision will your PC make?

## options
- A) Send directly to 192.168.2.50 using ARP, since both are in the 192.168.x.x range
- B) Send the packet to the default gateway because the destination is in a different subnet
- C) Broadcast the packet to all devices on the network to find 192.168.2.50
- D) Drop the packet because cross-subnet communication is not allowed

## answer
B

## explanation
The PC applies the subnet mask to both its own IP and the destination IP. Source: 192.168.1.10 AND 255.255.255.0 = 192.168.1.0 (network). Destination: 192.168.2.50 AND 255.255.255.0 = 192.168.2.0 (network). Since the networks differ (192.168.1.0 ≠ 192.168.2.0), the PC will forward the packet to its default gateway rather than attempting direct communication via ARP.

## hook
What would happen if the subnet mask was 255.255.0.0 (/16) instead?
