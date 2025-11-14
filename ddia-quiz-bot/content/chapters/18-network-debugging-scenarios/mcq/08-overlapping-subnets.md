---
id: debug-overlapping-subnets
day: 8
tags: [debugging, docker, routing, subnets, networks]
related_stories:
  - network-debugging
---

# Debugging Overlapping Network Subnets

## question
You create two custom Docker bridge networks: net1 (172.20.0.0/16) and net2 (172.20.0.0/16). Container A on net1 cannot reach Container B on net2 by IP, even though both are running. What's the problem?

## options
- A) Docker doesn't allow containers on different networks to communicate
- B) Both networks use the same subnet, creating routing ambiguity
- C) The veth pairs are not properly configured
- D) iptables is blocking traffic between the two bridges

## answer
B

## explanation
When two networks use overlapping subnets (both 172.20.0.0/16), the kernel's routing table cannot determine which network interface a packet should use to reach an IP in that range. This creates routing ambiguity. The first network's route wins, so packets destined for the second network might be routed to the first network's bridge instead. Docker networks CAN communicate if properly configured (via routing or overlay), veth pairs are automatic, and iptables would show different symptoms.

## hook
How does the kernel decide which route to use when multiple routes match the same destination?
