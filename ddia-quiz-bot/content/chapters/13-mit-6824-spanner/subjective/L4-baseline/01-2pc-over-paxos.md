---
id: spanner-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: spanner
subtopic: 2pc-over-paxos
estimated_time: 8-10 minutes
---

# question_title - 2PC over Paxos Participants

## main_question - Core Question
"Describe how Spanner uses two-phase commit (2PC) over Paxos-replicated shards. Why does this reduce classic 2PC blocking? Walk through what happens if the coordinator crashes during commit."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Each shard is a Paxos group with a leader
- Prepare/commit decisions are logged via Paxos (replicated write-ahead)
- Recovery by reading Paxos logs; decisions survive single-node failures
- Coordinator often co-located with a participant; any replica can recover state

### expected_keywords
- Primary: Paxos, 2PC, coordinator, prepare, commit
- Technical: majority quorum, replicated log, recovery

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Distinguish from 3PC and its assumptions
- Coordinator selection / re-election
- Idempotence of retries and commit decisions
- Interaction with lock holding and timestamp assignment

### bonus_keywords
- leader election, write-ahead log, durable vote, blocking conditions

## sample_excellent - Example Excellence
"Each participant is a Paxos group. On prepare, the leader appends a prepare record via Paxos; a majority makes the vote durable. On commit, the commit record is also replicated. If the coordinator crashes after sending some commits, a recovering coordinator or even participants can consult their Paxos logs to finish the decision consistently. Because votes/decisions are replicated, single-node failures don’t force indefinite blocking as in classic 2PC where the coordinator’s state can be lost."

## sample_acceptable - Minimum Acceptable
"Spanner runs 2PC across Paxos groups so decisions are replicated. If the coordinator dies, logs let the system finish the transaction."

## common_mistakes - Watch Out For
- Claiming Paxos lets a single replica decide unilaterally
- Ignoring the need for majority quorums
- Saying 2PC becomes non-blocking in all cases

## follow_up_excellent - Depth Probe
**Question**: "What’s the outcome if a participant times out after prepare but before commit, and then recovers?"
- **Looking for**: Read replicated log; commit/abort decision; idempotent application

## follow_up_partial - Guided Probe  
**Question**: "How does logging prepare in Paxos change recovery vs single-node participants?"
- **Hint embedded**: Durability of votes

## follow_up_weak - Foundation Check
**Question**: "What are the two phases in 2PC and what do they achieve?"
