---
id: debug-container-internet-connectivity
day: 1
tags: [debugging, docker, nat, iptables, routing]
related_stories:
  - network-debugging
---

# Debugging Container Internet Connectivity

## question
A Docker container (172.17.0.2) on the default bridge network cannot reach the internet. You run `iptables -t nat -L POSTROUTING` and see no MASQUERADE rule for 172.17.0.0/16. The host can reach the internet fine. What is the most likely problem?

## options
- A) The container's default gateway is not set correctly
- B) The Docker daemon didn't create the necessary NAT rule for outbound traffic
- C) The container's DNS server is misconfigured
- D) The veth pair connecting the container to the bridge is down

## answer
B

## explanation
For containers to reach the internet, Docker must set up a MASQUERADE (source NAT) rule in the POSTROUTING chain that translates the container's private IP (172.17.0.2) to the host's public IP. Without this rule, packets leave the container with a private source IP that cannot be routed back. The gateway and veth pair would cause different symptoms (no packets leaving at all), and DNS issues would only affect name resolution, not all internet connectivity.

## hook
How would you verify if packets are leaving the container but not being translated by NAT?
