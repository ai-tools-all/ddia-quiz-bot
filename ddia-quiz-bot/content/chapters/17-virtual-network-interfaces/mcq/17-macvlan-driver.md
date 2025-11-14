---
id: virtual-network-macvlan-driver
day: 17
tags: [networking, docker, macvlan, drivers]
related_stories:
  - docker-networking
---

# MACVLAN Network Driver

## question
What unique capability does the macvlan network driver provide to Docker containers?

## options
- A) Encrypted communication between containers across different hosts
- B) Each container gets its own MAC address and appears as a physical device on the LAN
- C) Automatic load balancing across multiple network interfaces
- D) Support for running containers without root privileges

## answer
B

## explanation
The macvlan driver allows containers to have their own MAC addresses and appear as physical devices directly on the network. This is useful for legacy applications that expect to be directly on the LAN, monitoring tools that need promiscuous mode, or scenarios requiring direct Layer 2 access.

## hook
What are the downsides of using macvlan compared to standard bridge networking?
