---
id: ch05-sloppy-quorums
day: 12
tags: [replication, sloppy-quorums, hinted-handoff, availability]
related_stories: []
---

# Sloppy Quorums and Hinted Handoff

## question
What is the main trade-off when using sloppy quorums with hinted handoff?

## options
- A) Lower latency but higher storage costs
- B) Better availability but potential to read stale data
- C) Stronger consistency but lower throughput
- D) Simpler implementation but less fault tolerance

## answer
B

## explanation
Sloppy quorums increase write availability by allowing writes to nodes outside the normal set when the designated nodes are unavailable. These temporary writes are handed off to the correct nodes later (hinted handoff). This improves availability but breaks the quorum intersection guarantee - reads might not see the latest writes since they might be temporarily stored on nodes outside the read quorum.

## hook
Why might an e-commerce site choose sloppy quorums over strict quorums?
