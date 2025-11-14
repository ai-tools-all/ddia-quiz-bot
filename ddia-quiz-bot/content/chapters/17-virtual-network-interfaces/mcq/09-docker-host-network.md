---
id: virtual-network-host-mode
day: 9
tags: [networking, docker, host-network, performance]
related_stories:
  - docker-networking
---

# Docker Host Network Mode

## question
A container running with `--network host` offers maximum performance but at what cost?

## options
- A) The container cannot access localhost services on the host
- B) The container loses network isolation and shares the host's network stack, including ports
- C) The container requires manual iptables configuration for all connections
- D) The container can only communicate with other containers using host network mode

## answer
B

## explanation
Host network mode removes the network namespace boundary, making the container share the host's network stack directly. This eliminates NAT and veth pair overhead (maximum performance) but also eliminates network isolation. The container sees all host interfaces and can conflict with host ports.

## hook
When would the performance benefit of host networking outweigh the security concerns?
