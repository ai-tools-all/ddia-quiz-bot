---
id: cops-subjective-L6-003
type: subjective
level: L6
category: advanced
topic: cops
subtopic: cross-datacenter-topology
estimated_time: 12-15 minutes
---

# question_title - Datacenter Topology and Replication Strategies

## main_question - Core Question
"You're deploying COPS for a global service with 5 datacenters (US-East, US-West, EU, Asia, Australia). Analyze how datacenter topology and network characteristics affect COPS performance. Design an optimized replication strategy considering: (1) replication lag, (2) dependency propagation paths, (3) partial replication, and (4) regional failures. Justify your design choices."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Network latency between DCs affects dependency propagation time
- Full replication: all DCs have all data (high consistency, high cost)
- Partial replication: some DCs have subsets (lower latency, complexity in routing)
- Dependency chains must propagate across full path (cascading delays)
- Regional failures require fallback paths without breaking causal consistency
- Write amplification: each put replicates to N-1 other DCs

### expected_keywords
- network topology, replication lag, full vs partial replication, dependency propagation, write amplification, failover

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Hub-and-spoke vs mesh topology for replication
- Batching dependencies to reduce overhead
- Regional affinity: clients prefer nearby DC
- Causal cuts: independent key subsets with separate replication
- Dependency compression across WAN
- Monitoring: track per-DC-pair propagation latency

### bonus_keywords
- hub-spoke topology, batching, regional affinity, causal cuts, WAN optimization, latency monitoring

## sample_excellent - Example Excellence
"**Design Approach:**

**1. Hybrid Topology:**
- **Regional hubs:** US-East, EU, Asia are primary hubs (full replication)
- **Spokes:** US-West, Australia partially replicate frequently-accessed keys
- Rationale: Reduces write amplification (5-way → 3-way for core), lower latency for reads in spoke regions, but requires dependency routing through hubs.

**2. Dependency Propagation:**
- **Direct mesh for hubs:** US-East ↔ EU ↔ Asia (lower latency for dependencies)
- **Hub-spoke for others:** US-West → US-East → other hubs
- Optimization: Batch dependencies every 10ms to reduce WAN round-trips (trade latency for throughput)
- Problem: Cascading delays worse for Australia (3-hop: Australia → Asia → EU → US-East)

**3. Partial Replication Strategy:**
- **Key affinity:** User data replicated to user's home region + nearest 2 DCs
- **Global data:** Catalogs, configs replicated everywhere
- **Dependency handling:** If put(X) depends on Y not local, either (a) fetch Y first, or (b) proxy put to Y's home DC
- Trade-off: Complexity in routing, but 60% reduction in storage and write amplification

**4. Regional Failure Handling:**
- **Fallback paths:** If Asia fails, Australia → US-West → US-East
- **Causal consistency risk:** Partial failure may block dependencies → use timeout degradation (after 5s, make visible with warning)
- **Monitoring:** Track per-DC-pair propagation SLAs, alert on >95th percentile breaches

**Metrics to Monitor:**
- P50/P95/P99 dependency wait time per DC-pair
- Cascading delay depth (max dependency chain hops)
- Write amplification ratio
- Visibility lag per region"

## sample_acceptable - Minimum Acceptable
"Full replication to all 5 DCs provides best consistency but high cost. Could use partial replication with regional affinity (users' data stays near home region). Network latency affects dependency propagation—closer DCs converge faster. For failures, route through alternate DCs but may break causal consistency during partitions. Monitor replication lag and dependency wait times per DC."

## common_mistakes - Watch Out For
- Not addressing partial replication option
- Missing dependency propagation path analysis
- Not discussing failure scenarios
- Ignoring write amplification costs

## follow_up_excellent - Depth Probe
**Question**: "How would you implement 'causal cuts'—identifying independent key subsets that can replicate separately to reduce cross-datacenter traffic?"
- **Looking for**: Static partitioning, dynamic analysis of access patterns, dependency graph analysis

## follow_up_partial - Guided Probe
**Question**: "If Australia clients frequently read data written in US-East, but with 200ms latency, how would you optimize this?"
- **Hint embedded**: Caching, prefetching, or full replication trade-offs

## follow_up_weak - Foundation Check
**Question**: "What is 'write amplification' and why does it increase with more datacenters?"
