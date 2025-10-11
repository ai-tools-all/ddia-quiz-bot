---
id: craq-subjective-L6-005
type: subjective
level: L6
category: baseline
topic: craq
subtopic: evolution-roadmap
estimated_time: 9-12 minutes
---

# question_title - CRAQ Evolution Roadmap

## main_question - Core Question
"Outline a two-year roadmap for evolving CRAQ in an enterprise environment. Include phases for observability, automation, cost optimization, and regulatory compliance, referencing DDIA's evolution and maintenance chapters." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Phase Planning**: Incremental rollout (e.g., observability first, automation second)
- **Instrumentation**: Expand metrics/logging/tracing aligned with DDIA operations guidance
- **Automation**: Introduce self-healing reconfiguration, autoscaling
- **Governance**: Address compliance, data retention, schema evolution

### expected_keywords
- Primary keywords: roadmap, observability, automation, compliance, cost
- Technical terms: phased rollout, feature flag, autoscaling, retention policy

## peripheral_concepts - Nice to Have (40%)
- **Stakeholder Alignment**: Product, SRE, compliance teams
- **Risk Mitigation**: Canary releases, rollback plans
- **Cost Modeling**: Evaluate replica sizing, hardware tiers
- **Feedback Loops**: Use postmortems and error budgets

### bonus_keywords
- Implementation: OKRs, milestones, change management, data catalog
- Scenarios: regulatory audit, hardware refresh, feature expansion
- Trade-offs: speed vs stability, innovation vs compliance

## sample_excellent - Example Excellence
"Year 1 focuses on visibility and safety: instrument dirty duration, tail lag, CDC backlog; build dashboards and alerting per DDIA operations best practices. Parallel effort adds runbooks and chaos drills. Year 2 introduces automation—autoscaling chains, automated rebalancing, and configuration manager failover scripts—guided by error budget feedback. We also address compliance by cataloging data, defining retention in CRAQ snapshots, and validating schema evolution with compatibility tests from DDIA's encoding chapter. Cost optimization includes moving cold chains to cheaper hardware and tuning chain lengths."

## sample_acceptable - Minimum Acceptable
"Plan phases: first improve observability (metrics, dashboards), then automate failover and scaling, and finally handle compliance and cost optimization—mirroring DDIA's maintenance advice." 

## common_mistakes - Watch Out For
- No phased plan or prioritization
- Ignoring compliance and schema evolution
- Lacking ties to DDIA's long-term maintenance recommendations
- Overlooking cost considerations

## follow_up_excellent - Depth Probe
**Question**: "What leading indicators tell you it's time to move from observability-focused phase to automation?"
- **Looking for**: MTTR trends, manual toil, error budget usage, incident volume
- **Red flags**: Arbitrary timeline

## follow_up_partial - Guided Probe  
**Question**: "How will regulatory requirements influence snapshot retention and audit logging?"
- **Hint embedded**: Data residency, retention policies, access audits
- **Concept testing**: Compliance awareness

## follow_up_weak - Foundation Check
**Question**: "Why plan renovations in stages instead of remodeling the whole house at once?"
- **Simplification**: Phased roadmap analogy
- **Building block**: Manage risk and resources

## bar_raiser_question - L6→L7 Challenge
"Integrate this roadmap with enterprise portfolio planning by mapping investments to business outcomes and risk reduction." 

### bar_raiser_concepts
- Portfolio management, ROI
- Risk quantification, business alignment
- Executive communication
- Value-focused prioritization

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: Portfolio planning, compliance, cost control
