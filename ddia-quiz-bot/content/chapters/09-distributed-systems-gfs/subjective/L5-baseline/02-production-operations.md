---
id: gfs-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: gfs
subtopic: production-operations
estimated_time: 10-12 minutes
---

# question_title - Operating GFS at Scale

## main_question - Core Question
"You're running a 10,000-node GFS cluster. Describe your operational playbook for handling: daily hardware failures, performance degradation, and capacity planning. What metrics and automation would you implement?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Failure as Normal**: Statistical inevitability at scale
- **Automated Recovery**: Self-healing without human intervention
- **Monitoring Hierarchy**: System, service, and business metrics
- **Capacity Forecasting**: Growth prediction and procurement cycles

### expected_keywords
- Primary keywords: MTTR, automation, monitoring, capacity, SLOs
- Technical terms: failure domains, blast radius, runbooks

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Chaos Engineering**: Proactive failure injection
- **Cost Optimization**: Utilization vs headroom trade-offs
- **Performance Regression**: Detection and attribution
- **Incident Management**: Escalation and post-mortems
- **Data Lifecycle**: Hot/warm/cold migration
- **Compliance**: Audit trails, data governance

### bonus_keywords
- Tools: Prometheus, Grafana, PagerDuty, Terraform
- Practices: SRE, DevOps, GitOps, Infrastructure as Code
- Metrics: P99 latency, IOPS, throughput, error budgets

## sample_excellent - Example Excellence
"Operating 10K nodes requires treating failures as routine: Daily ops: 1) Automated failure detection via heartbeats with immediate re-replication triggers. Expect ~30 disk failures/day (1% annual rate). 2) Health scoring per chunk server considering disk SMART stats, network errors, CPU temperature for predictive maintenance. Performance: 1) Tiered monitoring - business KPIs (request latency), service metrics (replication lag), system metrics (disk I/O). 2) Automated load balancing moving hot chunks when detecting skew. 3) Canary deployments with automatic rollback on regression. Capacity: 1) Model growth using historical trends plus business forecasts. 2) Maintain 20% headroom, trigger procurement at 70% capacity. 3) Just-in-time hardware delivery synchronized with demand. Key automation: Self-healing replication, automatic rebalancing, predictive failure replacement. Critical metrics: chunk distribution variance, replica lag P99, metadata operation latency, cost per GB. Use error budgets to balance reliability vs feature velocity."

## sample_acceptable - Minimum Acceptable
"With 10,000 nodes, I'd expect about 30 failures daily. Implement automated monitoring to detect failures quickly and trigger re-replication. Track key metrics like storage utilization, replication lag, and request latency. Plan capacity based on growth trends maintaining 20-30% headroom."

## common_mistakes - Watch Out For
- Underestimating failure frequency
- Manual intervention for common failures
- Missing business-level metrics
- No proactive capacity planning

## follow_up_excellent - Depth Probe
**Question**: "A critical bug causes 1% of writes to corrupt data silently. How do you detect, assess impact, and recover?"
- **Looking for**: Checksums, audit systems, blast radius assessment, recovery strategies
- **Red flags**: Not considering silent corruption possibility

## follow_up_partial - Guided Probe
**Question**: "You mentioned 30 failures per day. Walk me through the first 5 minutes after a chunk server fails."
- **Hint embedded**: Detection, assessment, prioritization
- **Concept testing**: Automation flow understanding

## follow_up_weak - Foundation Check
**Question**: "If you managed a fleet of 100 delivery trucks, how would you handle daily breakdowns?"
- **Simplification**: Physical fleet management
- **Building block**: Scale operations thinking

## bar_raiser_question - L5→L6 Challenge
"Design a 'self-driving' GFS that requires zero human operators. What AI/ML would you implement for autonomous operation?"

### bar_raiser_concepts
- Predictive failure models
- Automated root cause analysis
- Self-optimizing placement algorithms
- Anomaly detection and response
- Capacity planning ML models
- Cost optimization algorithms

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 5-7 min discussion
- **Common next topics**: SRE practices, ML for systems, autonomous infrastructure
