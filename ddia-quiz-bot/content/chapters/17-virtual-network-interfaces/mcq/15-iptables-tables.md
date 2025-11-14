---
id: virtual-network-iptables-tables
day: 15
tags: [networking, iptables, firewall, nat]
related_stories:
  - docker-networking
---

# iptables Tables

## question
Which iptables table is used for Network Address Translation (NAT) operations like port forwarding?

## options
- A) filter table
- B) nat table
- C) mangle table
- D) raw table

## answer
B

## explanation
The nat table handles all Network Address Translation operations. It contains PREROUTING (for DNAT/port forwarding), POSTROUTING (for SNAT/masquerading), and OUTPUT chains. The filter table handles packet filtering (allow/deny), mangle modifies packet headers, and raw handles connection tracking exemptions.

## hook
In what order does a packet traverse the different iptables tables?
