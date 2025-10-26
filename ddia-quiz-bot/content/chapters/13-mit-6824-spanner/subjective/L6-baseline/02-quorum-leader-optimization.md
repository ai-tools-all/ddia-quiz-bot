---
id: spanner-subjective-L6-002
type: subjective
level: L6
category: baseline
topic: spanner
subtopic: quorum-optimization
estimated_time: 12-15 minutes
---

# question_title - Quorum and Leader Optimization under Shifting Load

## main_question - Core Question
"Traffic is split 70/30 between Region A and B with strict write SLOs and RPO=0. Propose a replica layout (5 or 7 replicas), leader placement, and quorum path to minimize P50/P99 write latency while preserving availability under a single-region outage. Explain commit-wait’s role in the budget."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Majority quorum selection; tolerate minority failures
- Leader near dominant writers to reduce Paxos RTT
- Commit-wait adds ~ε regardless of geography; Paxos round-trip dominates
- Availability trade-offs with 5 vs 7 replicas and region spread

### expected_keywords
- quorum, leader placement, RTT, availability, tail latency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Non-voting/witness replicas for read scale without quorum inflation
- Dynamic leader moves as diurnal traffic shifts
- Client routing for read locality and freshness SLOs

### bonus_keywords
- region outage, election, follower reads, quorum geometry, cost

## sample_excellent - Example Excellence
"Use 5 replicas across 3+ regions: 2 in A, 2 in B, 1 in C (tie-breaker). Place leader in A. Quorum (3) is satisfied with A(leader)+B+A/B nearest, avoiding extra ocean hops; during A outage, B+C still form majority and elect a leader. Commit-wait (~ε) is added once; Paxos RTT (leader↔quorum) dominates, so co-locate quorum close to leader to reduce P50/P99. If traffic shifts, move leader to B during off-hours."

## sample_acceptable - Minimum Acceptable
"Pick a majority-friendly spread, place the leader near main writers, and ensure a majority remains after one region fails. Commit-wait is a small fixed add-on; Paxos RTT is the main lever."

## common_mistakes - Watch Out For
- Requiring acknowledgments from all replicas
- Confusing commit-wait with replication time
- Ignoring leader placement in tail latency

## follow_up_excellent - Depth Probe
**Question**: "When would you choose 7 replicas over 5, and how would you arrange quorums?"
- **Looking for**: Higher availability/read distribution vs write path length

## follow_up_partial - Guided Probe  
**Question**: "What happens to quorum and leader after A fails?"
- **Hint embedded**: Majority availability and re-election

## follow_up_weak - Foundation Check
**Question**: "What does a majority quorum guarantee?"
