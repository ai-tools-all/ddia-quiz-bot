---
id: virtual-network-overlay-networks
day: 18
tags: [networking, docker, overlay, swarm, multi-host]
related_stories:
  - docker-networking
---

# Overlay Network Purpose

## question
What problem do Docker overlay networks solve that bridge networks cannot?

## options
- A) Better performance for high-throughput applications
- B) Container communication across multiple physical hosts
- C) Support for more than 256 containers per network
- D) Automatic SSL/TLS encryption for HTTP traffic

## answer
B

## explanation
Overlay networks enable container-to-container communication across multiple Docker hosts, typically used in Docker Swarm or Kubernetes clusters. They create a virtual network layer that spans multiple physical machines, using VXLAN encapsulation to tunnel traffic between hosts while maintaining the same subnet for containers.

## hook
How does an overlay network handle routing packets between containers on different hosts?
