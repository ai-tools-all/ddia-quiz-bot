---
id: debug-docker-overlay-network-fails
day: 7
tags: [networking, troubleshooting, docker, overlay-network, vxlan, swarm]
related_stories:
  - overlay-networks
  - docker-swarm
---

# Multi-Host Docker Overlay Network Communication Fails

## question
In a Docker Swarm cluster, containers on different hosts cannot communicate over an overlay network, but containers on the same host can reach each other. The overlay network is properly created. What is the most likely network-level cause?

## options
- A) The containers have overlapping IP addresses across hosts
- B) UDP port 4789 (VXLAN) is blocked by a firewall between the hosts
- C) The docker_gwbridge is misconfigured
- D) The overlay network's subnet conflicts with the host network

## answer
B

## explanation
Docker overlay networks use VXLAN encapsulation to tunnel Layer 2 frames over Layer 3 networks between hosts. VXLAN uses UDP port 4789 by default. If this port is blocked by firewalls (host firewall, network firewall, security groups), VXLAN packets cannot traverse hosts, breaking inter-host communication. Intra-host communication works because it doesn't need encapsulation. Docker handles IP allocation (A) to prevent conflicts. docker_gwbridge (C) provides external connectivity but doesn't handle overlay communication. Subnet conflicts (D) would cause issues even on the same host.

## hook
How would you use tcpdump to verify if VXLAN packets are being sent and received on UDP port 4789?
