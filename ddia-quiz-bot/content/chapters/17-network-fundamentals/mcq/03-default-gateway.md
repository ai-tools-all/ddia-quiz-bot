---
id: network-default-gateway
day: 3
tags: [networking, gateway, routing, router]
related_stories:
  - network-basics
---

# Default Gateway

## question
What is the role of a default gateway in a network?

## options
- A) It assigns IP addresses to devices using DHCP
- B) It handles traffic destined for networks outside the local subnet
- C) It converts MAC addresses to IP addresses using ARP
- D) It provides DNS resolution for domain names

## answer
B

## explanation
A default gateway (typically a router) handles traffic destined for IP addresses outside the local subnet. When a device determines that the destination IP is not in its local network (by applying the subnet mask), it forwards the packet to the default gateway, which then routes it toward the destination network. DHCP, ARP, and DNS are separate functions.

## hook
How does a device decide whether to use the default gateway or send a packet directly to another device?
