---
id: virtual-network-iptables-chains
day: 10
tags: [networking, iptables, packet-flow, netfilter]
related_stories:
  - docker-networking
---

# iptables Chain Traversal

## question
An external client sends a packet to host port 8080, which is published from a container's port 80. Through which iptables chains does the packet traverse?

## options
- A) INPUT → FORWARD → OUTPUT
- B) PREROUTING → FORWARD → POSTROUTING
- C) PREROUTING → INPUT → OUTPUT → POSTROUTING
- D) INPUT → OUTPUT

## answer
B

## explanation
The packet arrives at the host's network interface and enters PREROUTING (where DNAT changes destination to container IP). After routing decision (forward to container, not local), it goes through FORWARD chain (firewall rules). Finally, it exits via POSTROUTING chain before reaching the container. The packet never touches INPUT or OUTPUT chains since it's not destined for a local process.

## hook
How would the chain traversal differ if the destination was a process running directly on the host?
