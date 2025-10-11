---
id: craq-subjective-L3-007
type: subjective
level: L3
category: baseline
topic: craq
subtopic: network-partition
estimated_time: 5-7 minutes
---

# question_title - CRAQ Under Network Partitions

## main_question - Core Question
"During a network partition, the configuration manager can still reach the head but not the tail. Explain why CRAQ prefers to halt writes rather than allow the reachable nodes to continue independently. Tie your answer to DDIA's discussion of the CAP theorem." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Consistency over Availability**: CRAQ chooses linearizability when partitioned
- **Tail Dependence**: Without tail acks, replicas stay dirty, blocking reads/writes
- **Configuration Manager Fencing**: Prevents split brain by pausing write authority
- **CAP Reference**: Aligns with choosing CP in DDIA's taxonomy

### expected_keywords
- Primary keywords: partition, consistency, availability trade-off, fencing, dirty state
- Technical terms: CAP theorem, CP system, network link failure, quorum equivalence

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Read Fallback**: Clients may still read from tail-side if accessible
- **Partial Availability**: Serve stale data only if explicitly configured (not default)
- **Operational Runbooks**: Similar to fail-stop runbooks in DDIA
- **Comparison to Dynamo**: Dynamo would pick availability with eventual consistency

### bonus_keywords
- Implementation: circuit breaker, write suspension, client retry policy
- Scenarios: datacenter partition, cross-region link loss, maintenance windows
- Trade-offs: downtime vs data correctness, customer communication

## sample_excellent - Example Excellence
"CRAQ relies on the tail's acknowledgment to mark data clean. If the manager can't see the tail, all new writes remain dirty. Rather than serve potentially divergent histories, the configuration manager fences the head (revokes its epoch) and stops accepting writes. This mirrors the CP stance in DDIA's CAP discussion: we sacrifice availability during partition to keep a single linearizable history. Systems like Dynamo make the opposite choice, but CRAQ is explicitly optimized for strong semantics." 

## sample_acceptable - Minimum Acceptable
"If the tail is unreachable, CRAQ can't mark entries clean, so it stops writes to avoid inconsistency. That matches the CAP discussion where a CP system preserves consistency during partitions." 

## common_mistakes - Watch Out For
- Claiming CRAQ stays available with stale reads by default
- Missing the role of tail acknowledgment in halting writes
- Confusing CRAQ's choice with Dynamo-style AP systems
- Ignoring configuration manager fencing tokens

## follow_up_excellent - Depth Probe
**Question**: "What operational steps could you take to shorten the partition window while still honoring the CP decision?"
- **Looking for**: Rapid detection, reroute via redundant links, promote new tail with majority agreement
- **Red flags**: Allowing isolated head to proceed alone

## follow_up_partial - Guided Probe  
**Question**: "If the head side wanted to keep serving reads, what restrictions would protect linearizability?"
- **Hint embedded**: Only previously clean data, no new writes
- **Concept testing**: Safe read set

## follow_up_weak - Foundation Check
**Question**: "If two people are holding opposite ends of a rope and it snaps, should one keep pulling as if nothing happened?"
- **Simplification**: Partition analogy
- **Building block**: Need both ends connected

## bar_raiser_question - L3→L4 Challenge
"Contrast CRAQ's CP decision with the leaderless quorum systems in DDIA's partitioning chapter. Under what customer requirements would you choose one design over the other?"

### bar_raiser_concepts
- Customer SLAs on consistency vs availability
- Conflict resolution costs
- Operational complexity of quorum repair
- Business trade-offs highlighted in DDIA

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Quorum design, replication modes, customer SLA mapping
