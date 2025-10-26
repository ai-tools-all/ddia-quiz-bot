---
id: spanner-2pc-over-paxos
day: 2
tags: [spanner, transactions, 2pc, paxos]
---

# Two-Phase Commit Over Paxos

## question
Why does Spanner run two-phase commit (2PC) over Paxos-replicated participants instead of single servers?

## options
- A) To avoid the need for a transaction coordinator altogether
- B) To ensure that prepare/commit decisions survive crashes without blocking progress
- C) To allow any replica to unilaterally commit if others are slow
- D) To eliminate the latency of the prepare phase entirely

## answer
B

## explanation
In Spanner, each participant shard is a Paxos group. 2PC messages (prepare/commit) are logged through Paxos so decisions are replicated and recoverable after crashes, avoiding classic 2PC blocking due to a failed coordinator or participant.

## hook
How does replicating 2PC logs change the failure modes of distributed commits?
