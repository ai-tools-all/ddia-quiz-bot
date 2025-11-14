---
id: virtual-network-packet-flow-outbound
day: 7
tags: [networking, packet-flow, containers, nat]
related_stories:
  - docker-networking
---

# Container Outbound Packet Flow

## question
In what order does a packet travel when a container makes an outbound HTTP request to the Internet?

## options
- A) Container eth0 → veth → docker0 → iptables → host eth0 → gateway
- B) Container eth0 → docker0 → veth → host eth0 → iptables → gateway
- C) Container eth0 → iptables → veth → docker0 → host eth0 → gateway
- D) Container eth0 → host eth0 → docker0 → iptables → veth → gateway

## answer
A

## explanation
Outbound packets flow: container's eth0 (172.17.0.2) → through veth pair → docker0 bridge (172.17.0.1) → iptables POSTROUTING (SNAT/masquerade to host IP) → host's eth0 → default gateway → Internet. The packet's source IP is rewritten from the container's private IP to the host's public IP.

## hook
What information must iptables track to ensure return packets reach the correct container?
