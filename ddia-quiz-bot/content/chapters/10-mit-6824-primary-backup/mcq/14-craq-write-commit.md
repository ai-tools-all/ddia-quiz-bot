---
id: craq-write-commit
day: 14
tags: [craq, chain-replication, writes]
related_stories:
  - vmware-ft
---

# CRAQ Write Commitment Point

## question
When does CRAQ’s chain replication pipeline consider a client write complete?

## options
- A) As soon as the head logs the update locally
- B) Once a majority of replicas acknowledge receipt
- C) After the tail replica applies the update and responds

## answer
C

## explanation
CRAQ inherits chain replication’s discipline: the write is passed head-to-tail, and only the tail’s acknowledgement signals that every replica has applied the update, so that is when the client can be told the write committed.

## hook
What safeguards prevent duplicate client requests when a write must be retried during failover?
