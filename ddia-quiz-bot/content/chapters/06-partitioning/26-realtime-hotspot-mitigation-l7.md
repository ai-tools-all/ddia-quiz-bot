---
id: ch06-realtime-hotspot-mitigation-l7
day: 26
level: L7
tags: [partitioning, hot-spots, real-time, principal-engineer, performance]
related_stories: []
---

# Real-Time Hot Spot Mitigation

## question
Your real-time bidding platform processes 10M requests/second, partitioned by ad campaign ID. During major events (Super Bowl, elections), certain campaigns experience 1000x traffic spikes, overwhelming their partitions. As Principal Engineer, design an adaptive partitioning system that: (1) detects hot spots within seconds, (2) redistributes load without losing in-flight requests, (3) maintains sub-10ms p99 latency, and (4) automatically scales down after spikes. Consider the cost of over-provisioning vs dynamic scaling.

## expected_concepts
- Dynamic partition splitting strategies
- Real-time load detection algorithms
- Consistent hashing with virtual nodes
- Request buffering during transitions
- Predictive scaling based on patterns
- Circuit breakers and backpressure
- Cost-aware resource allocation

## answer
The solution requires a multi-layered approach combining predictive and reactive strategies:

Detection layer: (1) **Sliding window metrics** per partition tracking QPS, latency percentiles, and queue depth with 1-second granularity. (2) **Adaptive thresholds** using exponential moving averages to detect anomalies (3x above baseline triggers alert). (3) **Predictive detection** using historical patterns and external signals (event calendars, trending topics).

Mitigation architecture: (1) **Virtual partition layer** - each physical partition hosts multiple virtual partitions (e.g., 100 virtual per physical). (2) **Dynamic splitting** - hot virtual partitions split into sub-partitions with deterministic hashing for consistent routing. (3) **Overflow pools** - pre-warmed standby capacity that hot partitions can instantly claim. (4) **Request shadowing** - gradually shift traffic to new partitions while draining old ones.

Implementation details: **Two-phase migration** - Phase 1: New partition shadows reads while old handles writes. Phase 2: Atomic switch using compare-and-swap on routing table. **Backpressure mechanism** - Clients implement exponential backoff when partition returns "splitting" status. **Smart batching** - Aggregate multiple campaigns into micro-batches during splits to maintain throughput.

Scale-down logic: Gradual consolidation using inverse of split logic. Keep capacity warm for 15 minutes post-spike (configurable based on cost tolerance). Maintain "heat map" of historical hot spots for faster future response.

Critical innovation: **Speculative pre-splitting** - ML model predicts hot spots 30 seconds ahead based on bid velocity changes, allowing proactive splitting before latency impact.

## hook
How do you handle cascading hot spots when migrating load creates new hot spots on destination partitions?

## follow_up
During a critical product launch, your system correctly handles the initial traffic spike through dynamic partitioning. However, 10 minutes in, a bug in the advertiser's code causes their campaigns to submit malformed requests that pass validation but cause expensive processing loops. These requests consume 100x normal CPU, and the hot spot detection focuses on QPS, not resource consumption. How do you evolve your system to handle resource-based hot spots while maintaining the real-time SLAs?

## follow_up_answer
This requires extending hot spot detection beyond traffic volume to resource consumption patterns:

Multi-dimensional detection: (1) **Resource accounting per request** - Track CPU time, memory allocation, and I/O operations per campaign at request level using lightweight profiling. (2) **Compound health scoring** - Combine QPS, CPU usage, memory pressure, and latency into unified partition health score using weighted formula: Health = 0.3*QPS_normalized + 0.4*CPU_normalized + 0.2*Memory_normalized + 0.1*Latency_normalized.

Request classification: (1) **Adaptive request profiling** - Sample 1% of requests for deep profiling, increase to 100% for campaigns showing anomalies. (2) **Request fingerprinting** - Hash request patterns to identify and cache expensive operations. (3) **Dynamic cost estimation** - ML model learns request cost based on features, predicts cost before execution.

Mitigation strategies: (1) **Resource-based rate limiting** - Instead of request count, limit by resource tokens where expensive requests consume more tokens. (2) **Quality of Service tiers** - Automatically downgrade campaigns causing resource spikes to batch processing tier. (3) **Circuit breaker enhancement** - Trip based on resource consumption rate, not just error rate.

Defensive architecture: (1) **Request sandboxing** - Execute suspicious requests in resource-constrained containers. (2) **Progressive rollback** - Automatically revert to previous campaign configuration if resource spike detected after update. (3) **Cost feedback loop** - Return resource cost to advertisers in response headers, enabling client-side optimization.

Long-term solution: Implement **admission control** based on learned request costs. New or updated campaigns go through graduated rollout with resource monitoring. Establish **resource budgets** per campaign with automatic throttling when exceeded.

Key insight: Traffic-based partitioning assumes uniform request cost - a dangerous assumption in multi-tenant systems. True hot spot mitigation must consider multiple resource dimensions and implement defense-in-depth against both intentional and accidental resource abuse.
