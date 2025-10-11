---
id: craq-subjective-L4-007
type: subjective
level: L4
category: baseline
topic: craq
subtopic: observability
estimated_time: 6-8 minutes
---

# question_title - Observability for CRAQ Cleanliness

## main_question - Core Question
"Design the key observability metrics and alerts you would implement to ensure CRAQ maintains the service-level guarantees discussed in DDIA's operations chapter. Highlight how these metrics map to the system's unique concepts." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Dirty Duration Tracking**: Time between write arrival and clean flag
- **Tail Lag**: Difference between tail commit index and head index
- **Configuration Manager Health**: Epoch agreement, leader election status
- **DDIA Link**: Aligns with measuring replication lag, SLA dashboards, and error budgets

### expected_keywords
- Primary keywords: dirty duration, tail lag, epoch, alerting, SLO
- Technical terms: histogram, percentile, health check, error budget

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Client-Facing Symptoms**: Read fallback rate, dirty-read rejections
- **Retry Storm Detection**: Spikes in duplicate suppression hits
- **Integration with Incident Response**: Runbook triggers from metrics
- **Comparative Baseline**: Similar to follower lag and commit latency in DDIA

### bonus_keywords
- Implementation: Prometheus counter, RED metrics, saturation, burn rate
- Scenarios: cross-region failure, configuration manager stalled, network congestion
- Trade-offs: alert fatigue vs sensitivity, aggregated vs per-shard metrics

## sample_excellent - Example Excellence
"I'd export dirty duration histograms for each replica, tail lag (head index minus tail index), configuration manager heartbeat latency, and duplicate request rate. Alerts would trigger if dirty duration p99 exceeds 2× our SLA or if tail lag stays above 100 entries for more than a minute, similar to follower lag alerting in DDIA. Configuration manager health mirrors the operations chapter's advice: we monitor epoch monotonicity and leader re-election churn."

## sample_acceptable - Minimum Acceptable
"Track how long replicas stay dirty, the lag between head and tail, and the configuration manager's health. Alert if those exceed thresholds, just like DDIA suggests for replication lag and coordination services."

## common_mistakes - Watch Out For
- Monitoring only CPU or memory without CRAQ semantics
- Ignoring configuration manager metrics
- No threshold translation to SLA impact
- Lacking linkage to DDIA operational best practices

## follow_up_excellent - Depth Probe
**Question**: "How would you correlate dirty duration spikes with downstream client latency to prove causation?"
- **Looking for**: Distributed tracing, client logs, statistical correlation
- **Red flags**: Assuming causation without data

## follow_up_partial - Guided Probe  
**Question**: "What is the first mitigation when tail lag alert fires?"
- **Hint embedded**: Investigate tail health, network, configuration manager status
- **Concept testing**: Operational runbook

## follow_up_weak - Foundation Check
**Question**: "Why do pilots track the distance between planes instead of only their speed?"
- **Simplification**: Lag metric analogy
- **Building block**: Relative positioning matters

## bar_raiser_question - L4→L5 Challenge
"Create an error budget policy that incorporates dirty duration and duplicate suppression metrics to decide when to slow rollout of new features."

### bar_raiser_concepts
- Linking SLOs to release gates
- Error budget consumption tied to replication health
- Feedback loop into change management
- Inspired by DDIA operational governance guidance

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Incident response, automation, chaos testing
