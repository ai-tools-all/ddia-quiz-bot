---
id: craq-configuration-manager
day: 15
tags: [craq, chain-replication, coordination]
related_stories:
  - vmware-ft
---

# CRAQ Configuration Control

## question
Why does CRAQ deploy a separate configuration manager alongside the replication chain?

## options
- A) To randomly reshuffle read requests for load balancing
- B) To keep all replicas in agreement about chain membership and avoid split-brain takeovers
- C) To compress replication logs before they traverse the network

## answer
B

## explanation
The lecture emphasizes that chain-based systems like CRAQ rely on an external authority—often itself run with Raft or Paxos—to declare which nodes form the head-to-tail chain so that independent nodes cannot diverge into conflicting replicas during network partitions.

## hook
What failure scenarios must the configuration manager guard against to keep the chain coherent?
