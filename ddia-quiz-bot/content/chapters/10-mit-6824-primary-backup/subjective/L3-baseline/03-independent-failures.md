---
id: primary-backup-subjective-L3-003
type: subjective
level: L3
category: baseline
topic: primary-backup
subtopic: failure-independence
estimated_time: 5-7 minutes
---

# question_title - Independent vs Correlated Failures

## main_question - Core Question
"Why does primary-backup replication require failures to be independent? Give examples of correlated failures that would defeat replication, and explain how to mitigate them."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Independence Assumption**: Failures in replicas should be uncorrelated
- **Correlated Failures**: Events affecting multiple replicas simultaneously
- **Mitigation Strategy**: Physical/logical separation of replicas

### expected_keywords
- Primary keywords: independent, correlated, simultaneous, separation
- Technical terms: datacenter, power failure, natural disaster

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Manufacturing Defects**: Same hardware batch failures
- **Geographic Distribution**: Different cities/regions for replicas
- **Power Grid Separation**: Independent power sources
- **Network Segmentation**: Different network paths/providers

### bonus_keywords
- Examples: earthquakes, floods, building fires, city-wide outages
- Solutions: multi-datacenter, cloud regions, diverse hardware
- Risk factors: same rack, same building, same vendor

## sample_excellent - Example Excellence
"Primary-backup replication assumes failures are independent because it can only tolerate one replica failing at a time. Correlated failures defeat this by taking down multiple replicas simultaneously. Examples include: earthquakes destroying all servers in the same datacenter, city-wide power outages affecting all local replicas, manufacturing defects in servers bought from the same batch causing simultaneous hardware failures, or building fires destroying all equipment in one location. To mitigate these risks, we need physical separation - placing replicas in different datacenters, ideally in different cities or regions with independent power grids and network providers. We should also consider hardware diversity, using different server models or vendors for replicas to avoid systematic hardware defects affecting all copies."

## sample_acceptable - Minimum Acceptable
"Replication needs failures to be independent so they don't all fail at once. Correlated failures like earthquakes, power outages affecting a whole building, or all servers overheating in the same rack would take down all replicas. To prevent this, put replicas in different locations - different racks, buildings, or even cities."

## common_mistakes - Watch Out For
- Focusing only on hardware failures, missing environmental causes
- Not explaining why correlation defeats replication
- Suggesting software diversity (doesn't work for identical service)
- Underestimating correlation risks in same location

## follow_up_excellent - Depth Probe
**Question**: "How would you balance the trade-off between geographic distribution for failure independence and the increased network latency it creates?"
- **Looking for**: Asynchronous vs synchronous replication, consistency models
- **Red flags**: Not recognizing the latency impact

## follow_up_partial - Guided Probe  
**Question**: "Your company bought 1000 identical servers. You use two for primary-backup replication. Six months later, both fail the same day. What likely happened?"
- **Hint embedded**: Manufacturing defect manifesting
- **Concept testing**: Understanding systematic failures

## follow_up_weak - Foundation Check
**Question**: "If you have important data on your laptop, would you back it up to another laptop sitting right next to it? Why or why not?"
- **Simplification**: Physical proximity risks
- **Building block**: Understanding correlation through location

## bar_raiser_question - L3→L4 Challenge
"A bank runs primary-backup replication. They have three scenarios: (1) Both replicas in same datacenter, different racks, (2) Replicas in two datacenters in San Francisco, (3) One replica in SF, one in New York. Rank these by failure independence and explain the trade-offs of each."

### bar_raiser_concepts
- Spectrum of independence levels
- Regional disaster risks
- Network latency implications
- Cost vs availability trade-offs

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: CAP theorem, geo-replication, disaster recovery
