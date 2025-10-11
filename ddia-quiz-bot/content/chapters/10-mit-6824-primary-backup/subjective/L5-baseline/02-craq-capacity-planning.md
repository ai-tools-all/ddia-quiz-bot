---
id: craq-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: craq
subtopic: capacity-planning
estimated_time: 8-10 minutes
---

# question_title - Capacity Planning for CRAQ Chains

## main_question - Core Question
"You are running CRAQ across three regions with heterogeneous latency and workload skew. Design a capacity plan that meets a 5 ms p95 read latency SLO and a 50 ms p95 write latency SLO while maintaining the linearizable semantics highlighted in DDIA. Discuss shard allocation, chain length, and the interplay with clean propagation." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Chain Length vs Latency Trade-off**: More replicas increase propagation delay
- **Shard Placement Strategy**: Partition data to keep hot shards on shorter chains
- **Region-Aware Routing**: Geo-proximate reads from clean replicas, writes centralized or tiered
- **DDIA Connection**: Balancing partitioning (sharding) and replication semantics

### expected_keywords
- Primary keywords: chain length, shard, latency SLO, propagation delay
- Technical terms: consistent hashing, hot shard isolation, geo-distribution, head locality

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Adaptive Chain Sizing**: Short chains for critical tables, longer for archival data
- **Ingress Rate Limiting**: Keep write volume within propagation budget
- **Hybrid Deployment**: Cross-region read caches backed by CRAQ
- **Monitoring Tie-in**: Use replication lag metrics to adjust capacity

### bonus_keywords
- Implementation: autoscaling, predictive analytics, load shedding, traffic steering
- Scenarios: flash sale, regional outage, nightly batch job
- Trade-offs: operational complexity, cost, failover planning

## sample_excellent - Example Excellence
"To meet 5 ms read SLO, I'd keep each user-facing shard on a three-node chain within a single region, replicating asynchronously to secondary chains for DR. Writes enter at the regional head, traverse two local hops, and reach the tail within ~3×2 ms plus ack time, fitting the 50 ms budget. Hot shards get dedicated chains to avoid propagation contention; cold shards can share longer chains. Routing picks the nearest clean replica, updating the configuration manager to prefer low-latency nodes. This combines DDIA's sharding guidance with CRAQ's clean/dirty semantics." 

## sample_acceptable - Minimum Acceptable
"Shorter CRAQ chains give lower write latency, so I'd shard workloads so that hot shards stay on three-node regional chains and route reads to the nearest clean replica. This keeps us inside the 5 ms/50 ms SLOs while following DDIA's advice on balancing partitioning with replication guarantees." 

## common_mistakes - Watch Out For
- Using one global chain regardless of geo latency
- Ignoring effect of chain length on propagation
- Forgetting to link plan to DDIA partitioning guidance
- No discussion of monitoring or adjustments

## follow_up_excellent - Depth Probe
**Question**: "How would you adapt this plan if one region experiences sustained network jitter doubling the RTT?"
- **Looking for**: Temporary traffic shift, dynamic chain rebalancing, SLO impact analysis
- **Red flags**: Leaving plan static regardless of conditions

## follow_up_partial - Guided Probe  
**Question**: "Which metrics tell you a shard needs to be split or its chain shortened?"
- **Hint embedded**: Dirty duration, tail lag, write queue length
- **Concept testing**: Data-driven scaling

## follow_up_weak - Foundation Check
**Question**: "Why do theme parks open extra queues when lines get long?"
- **Simplification**: Shard/chain scaling analogy
- **Building block**: Reducing wait time by splitting load

## bar_raiser_question - L5→L6 Challenge
"Propose an automated control loop that continuously tunes CRAQ chain lengths using DDIA's error budget framework."

### bar_raiser_concepts
- Feedback loops, SLO-based scaling
- Predictive modeling of propagation delay
- Cost vs reliability trade-offs
- Automation and guardrails

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Adaptive routing, chaos testing, predictive autoscaling
