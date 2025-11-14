---
id: virtual-network-docker-bridge
day: 3
tags: [networking, docker, bridge, default-network]
related_stories:
  - docker-networking
---

# Docker Bridge Network

## question
What IP address range does Docker's default bridge network (docker0) typically use?

## options
- A) 10.0.0.0/8
- B) 192.168.0.0/16
- C) 172.17.0.0/16
- D) 169.254.0.0/16

## answer
C

## explanation
Docker's default bridge network uses the 172.17.0.0/16 subnet, with the bridge interface getting 172.17.0.1 and containers receiving IPs like 172.17.0.2, 172.17.0.3, etc. This is a private network segment isolated from the external network.

## hook
How do containers on this private network reach the Internet?
