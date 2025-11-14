---
id: virtual-network-iptables-nat
day: 5
tags: [networking, iptables, nat, masquerading]
related_stories:
  - docker-networking
---

# iptables NAT and Masquerading

## question
What is the primary difference between DNAT (Destination NAT) and SNAT (Source NAT/Masquerading)?

## options
- A) DNAT modifies the destination IP for inbound traffic, while SNAT modifies the source IP for outbound traffic
- B) DNAT is used for IPv4 and SNAT is used for IPv6 traffic
- C) DNAT provides encryption while SNAT provides compression
- D) DNAT works at Layer 2 while SNAT works at Layer 3

## answer
A

## explanation
DNAT (in PREROUTING chain) rewrites the destination IP address for incoming packets, allowing port forwarding to containers. SNAT/Masquerading (in POSTROUTING chain) rewrites the source IP address for outgoing packets, allowing containers on private networks to access the Internet using the host's IP.

## hook
Why does connection tracking (conntrack) need to remember both DNAT and SNAT transformations?
