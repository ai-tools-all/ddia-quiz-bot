---
id: virtual-network-troubleshooting
day: 19
tags: [networking, troubleshooting, docker, debugging]
related_stories:
  - docker-networking
---

# Container Network Troubleshooting

## question
A container cannot reach the Internet, but can ping other containers on the same bridge. What is the most likely cause?

## options
- A) The docker0 bridge is down
- B) The container's DNS is misconfigured
- C) IP forwarding is disabled on the host or iptables POSTROUTING masquerade rule is missing
- D) The veth pair is disconnected

## answer
C

## explanation
If containers can communicate with each other but not reach the Internet, the issue is likely with NAT/forwarding. The host needs IP forwarding enabled (sysctl net.ipv4.ip_forward=1) and iptables MASQUERADE rules in POSTROUTING to translate container IPs to the host's IP for outbound traffic. Container-to-container works because it's local Layer 2 forwarding.

## hook
How would you verify that IP forwarding is enabled and check the iptables masquerade rules?
