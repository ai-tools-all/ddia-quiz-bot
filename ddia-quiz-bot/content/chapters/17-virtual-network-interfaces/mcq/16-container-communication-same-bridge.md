---
id: virtual-network-container-same-bridge
day: 16
tags: [networking, docker, bridge, layer2]
related_stories:
  - docker-networking
---

# Container-to-Container Communication

## question
When Container A (172.17.0.2) sends a packet to Container B (172.17.0.3) on the same docker0 bridge, does the packet go through iptables NAT?

## options
- A) Yes, all packets must go through NAT to leave a container
- B) No, the bridge forwards the packet directly at Layer 2 without NAT
- C) Only if the containers are on different subnets
- D) Yes, but only through PREROUTING, not POSTROUTING

## answer
B

## explanation
Containers on the same bridge network communicate directly at Layer 2. The packet goes: Container A's eth0 → veth → docker0 bridge → veth → Container B's eth0. The bridge uses MAC address lookup to forward packets. NAT is only needed when crossing network boundaries (e.g., container to Internet).

## hook
How does Container A discover Container B's MAC address in the first place?
