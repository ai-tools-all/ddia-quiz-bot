---
id: ch08-quorum-accuracy
day: 9
level: L5
tags: [quorums, consensus, distributed-systems]
related_stories: []
---

# Quorum Systems and Network Delays

## question
In a 5-node cluster using quorum reads/writes (w=3, r=3), what happens if network delays cause 2 nodes to be slow but not dead?

## options
- A) The system continues normally with no impact
- B) Writes succeed but reads may timeout waiting for 3 responses
- C) Both reads and writes may experience increased latency but remain correct
- D) The system immediately fails

## answer
C

## explanation
With w=3 and r=3 in a 5-node system, operations need responses from 3 nodes. If 2 nodes are slow, the system can still achieve quorum using the 3 fast nodes, but latency increases as operations wait for the 3rd response. The system remains correct (w + r > n ensures overlap) but performance degrades. This illustrates how network delays affect distributed systems even when nodes haven't failed.

## hook
How do slow nodes affect quorum-based systems differently than failed nodes?
