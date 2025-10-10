---
id: raft-l5-consensus
title: "Raft Consensus Algorithm"
level: "L5"
category: "Distributed Systems"
---

## Question
Explain how the Raft consensus algorithm achieves distributed consensus. Include the roles of leader election, log replication, and safety properties in your answer.

## Core Concepts
- Leader election mechanism
- Log replication process
- Commitment rules
- Safety properties (election safety, leader append-only, log matching)
- Term numbers and their purpose

## Peripheral Concepts
- Split vote handling
- Client interaction model
- Configuration changes
- Log compaction
- Performance optimizations

## Sample Excellent Answer
The Raft consensus algorithm achieves distributed consensus by decomposing the problem into three sub-problems: leader election, log replication, and safety.

Leader Election: Raft uses terms (logical time periods) and randomized timeouts to elect leaders. Nodes start as followers, become candidates when they timeout, and request votes from other nodes. A candidate becomes leader upon receiving votes from a majority. The randomized timeouts help prevent split votes.

Log Replication: The leader receives client commands, appends them to its log, and replicates entries to followers through AppendEntries RPCs. Followers acknowledge receipt, and the leader commits entries once a majority has replicated them. The leader then notifies followers of committed entries.

Safety Properties: Raft ensures safety through several mechanisms:
1. Election Safety: Only one leader per term through majority voting
2. Leader Append-Only: Leaders never overwrite their logs
3. Log Matching: If two logs contain an entry with the same index and term, all preceding entries are identical
4. Leader Completeness: Committed entries appear in all future leaders' logs
5. State Machine Safety: All servers apply the same commands in the same order

The election restriction ensures only nodes with up-to-date logs can become leaders, preventing data loss.

## Sample Acceptable Answer
Raft achieves consensus through a leader-based approach. One node is elected as the leader using majority voting and term numbers. The leader handles all client requests and replicates log entries to follower nodes. Entries are considered committed when a majority of nodes have stored them. Safety is ensured through the requirement that leaders must have all committed entries and followers only accept logs from valid leaders with higher term numbers.

## Evaluation Rubric
Technical Accuracy: 35%
Completeness of Core Concepts: 30%
Understanding of Safety Properties: 20%
Clarity of Explanation: 10%
Use of Examples: 5%
