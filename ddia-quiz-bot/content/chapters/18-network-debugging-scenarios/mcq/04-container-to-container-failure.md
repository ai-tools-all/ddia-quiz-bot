---
id: debug-container-to-container
day: 4
tags: [debugging, docker, bridge, veth, layer2]
related_stories:
  - network-debugging
---

# Debugging Container-to-Container Communication

## question
Two containers (172.17.0.2 and 172.17.0.3) are on the docker0 bridge but cannot ping each other. `brctl show` confirms both veth interfaces are attached to docker0. What's the most likely cause?

## options
- A) The bridge is not forwarding packets between interfaces
- B) Each container's default gateway is misconfigured
- C) iptables FORWARD chain is dropping packets between containers
- D) The containers are in different network namespaces

## answer
C

## explanation
By default, the Linux kernel's iptables FORWARD chain policy might be set to DROP. Docker containers on the same bridge communicate through Layer 2 forwarding, but if iptables FORWARD rules block this traffic, packets are dropped. You can verify with `iptables -L FORWARD -v`. The gateway is only needed for external traffic. Network namespaces are expected (each container has its own), and the bridge should forward by default unless iptables blocks it.

## hook
How does a bridge decide whether to forward a packet between two interfaces?
