---
id: debug-iptables-forward-policy
day: 9
tags: [debugging, iptables, forward, docker, policy]
related_stories:
  - network-debugging
---

# Debugging iptables FORWARD Policy

## question
After installing Docker on a system with a strict firewall, existing containers suddenly cannot reach the internet. Running `iptables -L FORWARD` shows policy is DROP with no ACCEPT rules. What's the issue?

## options
- A) Docker's NAT rules were not created in the POSTROUTING chain
- B) The FORWARD chain's default policy blocks container traffic that must be forwarded
- C) The container's veth interface is not connected to docker0
- D) The host's default gateway is unreachable

## answer
B

## explanation
Container traffic to the internet must be forwarded from the container's namespace through the host's network stack. With FORWARD policy set to DROP and no explicit ACCEPT rules, this forwarding is blocked. Docker typically adds rules to allow forwarding for its networks, but a strict firewall might override these. The NAT rules (option A) are separate and handle address translation after forwarding. veth connectivity (C) would prevent any communication, and gateway issues (D) would affect the host too.

## hook
What's the difference between the FORWARD chain and the POSTROUTING chain in iptables?
