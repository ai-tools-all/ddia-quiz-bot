---
id: spanner-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: spanner
subtopic: 2pc-failures
estimated_time: 9-11 minutes
---

# question_title - 2PC over Paxos: Failure Recovery Analysis

## main_question - Core Question
"Walk through a transaction that reaches prepare on all participants and then the coordinator crashes. Explain how Spanner finishes the decision using Paxos logs at participants. Under what conditions can progress still be blocked?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Prepare and commit records replicated via Paxos at each shard
- Durable votes/decisions allow recovery by any new leader
- Idempotent re-issue of commit/abort based on replicated state
- Progress depends on majority quorums at participants

### expected_keywords
- Paxos log, majority quorum, prepare record, commit record, recovery

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Coordinator re-election/relocation and how it discovers state
- Timeouts and retry semantics across groups
- Classic 2PC blocking vs reduced blocking in Spanner

### bonus_keywords
- write-ahead, durable vote, in-doubt transaction, leader failover

## sample_excellent - Example Excellence
"On prepare, each participant’s leader appends a prepare entry via Paxos; a majority makes the ‘yes’ vote durable. If the coordinator crashes post-prepare, a new coordinator consults participants. Since votes are in Paxos logs, leaders can confirm prepared=yes, and if a commit record was replicated anywhere, participants converge on commit idempotently. Blocking only occurs if a participant’s Paxos group lacks a majority (e.g., region outage), since neither state nor new entries can be confirmed."

## sample_acceptable - Minimum Acceptable
"Votes and decisions are in Paxos logs so recovery can finish the transaction. You still need a majority at each group to read/append the recovery entries."

## common_mistakes - Watch Out For
- Claiming no blocking is possible in Spanner
- Letting a single follower decide without a quorum
- Ignoring idempotency when resending commit

## follow_up_excellent - Depth Probe
**Question**: "If some participants logged commit and others only prepared, how do you reconcile?"
- **Looking for**: Commit wins via replicated evidence; idempotent apply

## follow_up_partial - Guided Probe  
**Question**: "Why is a majority quorum required even during recovery reads?"
- **Hint embedded**: Leader change and Paxos safety

## follow_up_weak - Foundation Check
**Question**: "What are the two phases of 2PC and their purpose?"
