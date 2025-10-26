---
id: spanner-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: spanner
subtopic: geo-deployment
estimated_time: 12-15 minutes
---

# question_title - Geo-Distributed Deployment Trade-offs

## main_question - Core Question
"You’re deploying Spanner globally. Discuss leader placement, replica counts, and quorum selection to balance latency and availability. How do these choices interact with commit-wait and read locality? Consider a region outage scenario."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Paxos majority quorums tolerate minority failures
- Leader proximity affects write latency; followers enable local reads
- Commit-wait adds fixed latency bound by TrueTime uncertainty
- Trade-off: more regions improve availability but can raise quorum latency

### expected_keywords
- Primary: quorum, leader placement, latency, availability
- Technical: majority, cross-region RTT, region outage, failover

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Read routes to nearest fresh replica; fallback to leader if lagging
- Dynamic leader moves to optimize write path
- Witness/readonly replicas vs voting members
- Cost vs SLO trade-offs (tail latency budgets)

### bonus_keywords
- multi-region topology, follower reads, freshness SLOs, failover drills

## sample_excellent - Example Excellence
"Use 5 replicas across 3+ regions; place the leader near write-heavy clients to minimize leader RTT. Majority quorum (3) survives one region loss. Commit-wait adds ~ε to writes regardless of geography, but Paxos round-trip dominates; minimize cross-ocean hops for the leader’s quorum. Serve reads from local followers when their applied-through ≥ requested ts; otherwise route to leader or wait. During a region outage, re-elect a leader where a majority remains and keep read SLOs by rerouting to fresh replicas."

## sample_acceptable - Minimum Acceptable
"Choose a majority quorum and place the leader to reduce latency. Followers serve local reads if fresh; otherwise wait or go to the leader."

## common_mistakes - Watch Out For
- Assuming all replicas must ack every write
- Ignoring leader placement impact on latency
- Confusing commit-wait with replication delay

## follow_up_excellent - Depth Probe
**Question**: "When would you add non-voting replicas, and what are their trade-offs?"
- **Looking for**: Read scaling vs quorum size, failover behavior

## follow_up_partial - Guided Probe  
**Question**: "How does a region failure affect quorum and leader election?"
- **Hint embedded**: Majority semantics

## follow_up_weak - Foundation Check
**Question**: "What does a majority quorum guarantee in Paxos?"
