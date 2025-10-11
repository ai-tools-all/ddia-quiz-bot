---
id: craq-subjective-L3-005
type: subjective
level: L3
category: baseline
topic: craq
subtopic: configuration-manager
estimated_time: 5-7 minutes
---

# question_title - Why CRAQ Needs a Configuration Manager

## main_question - Core Question
"Explain the role of the configuration manager in CRAQ. How does it prevent split-brain situations, and which DDIA coordination building block does it resemble?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **External Authority**: Separate service defines head/tail membership
- **Split-Brain Prevention**: Stops multiple chains from acting independently
- **Replica Liveness Tracking**: Detects failed nodes and reconfigures chain
- **DDIA Link**: Mirrors Zookeeper-style coordination services

### expected_keywords
- Primary keywords: configuration manager, membership, split-brain, failover
- Technical terms: Zookeeper, leader election, Raft-backed metadata, view change

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Failure Domains**: Manager typically replicated with quorum protocol
- **Change Propagation**: Pushes new configuration downstream
- **Interaction with Dirty/Clean**: Ensures consistent metadata during swaps
- **Comparison to Primary-Backup Controllers**: Similar to VMware FT metadata host

### bonus_keywords
- Implementation: epoch number, lease, fencing token
- Scenarios: network partition, tail replacement, head promotion
- Trade-offs: additional latency, dependency on consensus layer

## sample_excellent - Example Excellence
"CRAQ deliberately outsources membership to a configuration manager so that only one chain definition exists at any given time. That manager—often itself a Raft or Paxos cluster like the Zookeeper service described in DDIA—tracks which replicas are alive, assigns the head and tail, and tells clients where to read and write. If the chain were to promote a new head locally during a partition, we could end up with two conflicting tails (split brain). The configuration manager hands out monotonically increasing epochs and fencing tokens so only the authorized chain processes requests. It's the same coordination primitive Kleppmann calls out for leader election and lock management." 

## sample_acceptable - Minimum Acceptable
"CRAQ uses a separate configuration manager, much like Zookeeper, to decide which nodes are the head and tail. That avoids split brain because everyone follows the external authority when a node fails."

## common_mistakes - Watch Out For
- Believing the chain can self-elect safely without outside help
- Ignoring epoch/fencing ideas from DDIA coordination chapter
- Forgetting the need to inform clients of topology changes
- Assuming configuration manager handles data replication itself

## follow_up_excellent - Depth Probe
**Question**: "Walk through how the configuration manager would handle a head failure while a write is mid-flight."
- **Looking for**: Detect failure, choose successor, bump epoch, ensure dirty replicas reconcile, client retry with new head
- **Red flags**: Allowing old head to resume without fencing

## follow_up_partial - Guided Probe  
**Question**: "What happens if two configuration managers accidentally run at once?"
- **Hint embedded**: Conflicting chain definitions
- **Concept testing**: Importance of consensus-backed manager

## follow_up_weak - Foundation Check
**Question**: "In a relay race, who decides the running order? Why is a coach-like authority helpful?"
- **Simplification**: Chain membership analogy
- **Building block**: Central coordinator prevents chaos

## bar_raiser_question - L3→L4 Challenge
"Compare the configuration manager in CRAQ to the metadata management in Dynamo-style systems from DDIA. When would you favor one over the other for membership control?"

### bar_raiser_concepts
- Static hashing vs explicit ordering
- Consistent hashing ring vs single chain list
- Conflict resolution vs strict serialization
- Operational complexity trade-offs

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Epoch handling, failover simulation, consensus under the hood
