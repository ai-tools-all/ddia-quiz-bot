---
id: farm-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: system-design
subtopic: trade-offs
estimated_time: 12-15 minutes
---

# question_title - FaRM vs Spanner Trade-offs

## main_question - Core Question
"Compare and contrast FaRM and Spanner across the dimensions of latency, consistency, fault tolerance, and geographic distribution. For each dimension, explain the fundamental design choices that drive the differences and the implications for real-world deployment. Under what application scenarios would you choose FaRM over Spanner, and vice versa?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Latency: FaRM ~58μs (single datacenter, RDMA), Spanner ~10ms (geo-distributed, WAN)
- Consistency: both provide serializability, but different mechanisms (OCC+versions vs TrueTime+2PL)
- Fault tolerance: FaRM F+1 replicas in one DC (power failure tolerant), Spanner Paxos across regions (disaster tolerant)
- Geographic distribution: FaRM single DC, Spanner multi-region
- Use cases: FaRM for ultra-low latency within region, Spanner for global data and disaster recovery

### expected_keywords
- latency, geo-replication, disaster recovery, serializability, RDMA, WAN, CAP theorem, trade-off

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- FaRM capacity limited by single DC NVRAM, Spanner scales globally
- FaRM's OCC abort rate under contention vs Spanner's wound-wait
- Operational complexity: FaRM needs special hardware (RDMA, NVRAM)
- Cost models differ significantly
- Hybrid approaches: FaRM within region, cross-region async or eventual consistency

### bonus_keywords
- capacity, contention, wound-wait, operational cost, hybrid architecture, hardware requirements

## sample_excellent - Example Excellence
"FaRM optimizes for single-DC latency (~58μs) via RDMA and NVRAM, sacrificing geographic distribution; Spanner prioritizes global consistency across continents (~10ms) via TrueTime and Paxos, sacrificing latency. FaRM tolerates DC power failures but not regional disasters; Spanner survives multi-region failures. Choose FaRM for: high-frequency trading, in-memory caching layers, latency-critical services within one region. Choose Spanner for: global user bases, regulatory data residency, disaster recovery requirements, workloads tolerating millisecond latencies. Hybrid: FaRM per-region with cross-region async replication."

## sample_acceptable - Minimum Acceptable
"FaRM is faster (microseconds) but only works in one datacenter. Spanner is slower (milliseconds) but works globally. Use FaRM when you need speed in one location, Spanner when you need global availability."

## common_mistakes - Watch Out For
- Claiming FaRM doesn't provide serializability
- Missing the hardware dependency (RDMA, NVRAM) as a trade-off factor
- Not distinguishing fault tolerance within DC vs across regions
- Oversimplifying to "FaRM is always faster"

## follow_up_excellent - Depth Probe
**Question**: "Design a hybrid system using FaRM for local transactions and Spanner for global coordination. What consistency model would you provide?"
- **Looking for**: Local serializability within FaRM regions, eventual or causal consistency across regions, conflict resolution strategies, compensating transactions

## follow_up_partial - Guided Probe
**Question**: "What happens to FaRM's latency advantage if you need to replicate across continents?"
- **Hint embedded**: WAN latency dominates, negating RDMA benefits

## follow_up_weak - Foundation Check
**Question**: "What does geographic distribution mean and why does it affect latency?"
