---
id: ch05-consensus-vs-replication
day: 19
tags: [replication, consensus, distributed-systems, coordination]
related_stories: []
---

# Consensus vs Replication

## question
What is the key difference between replication and consensus in distributed systems?

## options
- A) Replication is about data copies; consensus is about agreeing on values
- B) Replication requires more nodes than consensus
- C) Consensus is faster than replication
- D) They are essentially the same concept

## answer
A

## explanation
Replication focuses on maintaining copies of data across nodes for availability and performance. Consensus is about getting multiple nodes to agree on a single value or decision (like which node is the leader). While replication often uses consensus (e.g., for leader election or commit decisions), they solve different problems. Consensus algorithms like Raft or Paxos ensure agreement despite failures.

## hook
Why do some systems use consensus for replication while others avoid it?
