---
id: virtual-network-veth-pairs
day: 2
tags: [networking, virtual-interfaces, veth, containers]
related_stories:
  - docker-networking
---

# Virtual Ethernet (veth) Pairs

## question
What best describes a veth (virtual Ethernet) pair used in container networking?

## options
- A) A hardware network adapter that supports virtualization extensions
- B) A bidirectional virtual cable connecting two network namespaces
- C) A protocol for encrypting traffic between containers
- D) A load balancer that distributes traffic across multiple containers

## answer
B

## explanation
A veth pair acts like a virtual Ethernet cable with two ends. Packets sent to one end appear at the other end. Docker uses veth pairs to connect containers to the host's bridge network, with one end in the container namespace and the other end on the host.

## hook
How does Docker ensure network isolation between containers while allowing them to communicate?
