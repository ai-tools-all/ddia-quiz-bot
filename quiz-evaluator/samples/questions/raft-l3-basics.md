---
id: raft-l3-basics
title: "Basic Understanding of Raft"
level: "L3"
category: "Distributed Systems"
---

## Question
What is the Raft consensus algorithm and what problem does it solve in distributed systems?

## Core Concepts
- Purpose of consensus algorithms
- Basic leader-follower model
- Fault tolerance
- Data consistency across nodes

## Peripheral Concepts
- Comparison with other consensus algorithms
- Use cases and applications
- Performance characteristics

## Sample Excellent Answer
The Raft consensus algorithm is a protocol designed to achieve consensus in distributed systems, solving the fundamental problem of maintaining consistent state across multiple nodes even in the presence of failures.

The problem it solves: In distributed systems, multiple nodes need to agree on shared state or values. This becomes challenging when nodes can fail, messages can be delayed or lost, and there's no global clock. Without consensus, different nodes might have different views of the system state, leading to inconsistencies.

Raft solves this through a leader-based approach where:
1. One node is elected as the leader through a democratic voting process
2. All state changes go through the leader, providing a single source of truth
3. The leader replicates changes to follower nodes
4. Changes are only committed after a majority of nodes acknowledge them
5. The system can tolerate failures of up to (n-1)/2 nodes in an n-node cluster

This ensures all nodes eventually have the same ordered sequence of state changes, maintaining consistency even when some nodes fail or recover.

## Sample Acceptable Answer
Raft is a consensus algorithm used in distributed systems to ensure multiple nodes agree on the same data. It solves the problem of keeping data consistent across different servers even when some servers fail. Raft works by electing one node as a leader who manages all changes, and the leader makes sure all other nodes (followers) get the same updates in the same order.

## Evaluation Rubric
Understanding of Problem: 30%
Explanation of Solution: 30%
Technical Accuracy: 25%
Clarity: 15%
