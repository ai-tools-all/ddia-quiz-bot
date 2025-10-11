---
id: craq-subjective-L7-004
type: subjective
level: L7
category: baseline
topic: craq
subtopic: sustainability-cost
estimated_time: 12-15 minutes
---

# question_title - CRAQ Sustainability and Cost Leadership

## main_question - Core Question
"Develop a sustainability and cost leadership program for CRAQ that reduces total cost of ownership and environmental footprint while maintaining DDIA-grade reliability. Include hardware strategy, workload placement, and financial governance." 

## core_concepts - Must Mention (60%)
- **Hardware Optimization**: Energy-efficient instances, workload-aware sizing, hardware acceleration
- **Workload Placement**: Move cold chains to greener regions/off-peak compute, maintain SLA awareness
- **Financial Governance**: Chargeback/showback, cost transparency, investment prioritization
- **Reliability Guardrails**: Ensure cost savings do not erode availability/consistency guarantees

### expected_keywords
- Primary keywords: sustainability, TCO, workload placement, guardrails, governance
- Technical terms: energy footprint, autoscaling, tiered storage, chargeback

## peripheral_concepts - Nice to Have (40%)
- **Metrics**: Carbon per request, cost per transaction, efficiency KPIs
- **Process**: Green capacity planning, procurement policy, supplier audits
- **Innovation**: ARM/x86 mix, liquid cooling, hardware offload for replication
- **Risk**: SLA impact, regulatory sustainability mandates

### bonus_keywords
- Implementation: energy dashboards, carbon-aware scheduling, reserved instances, dynamic rightsizing
- Scenarios: peak holiday traffic, regulatory ESG reporting, data center outages
- Trade-offs: cost vs redundancy, green energy availability vs latency

## sample_excellent - Example Excellence
"Adopt energy-efficient hardware (ARM-based nodes for cold chains, SSD tiering) and schedule non-critical workloads in carbon-friendly regions using carbon-aware placement policies. Implement chargeback showing cost/carbon per request, tied to product teams' error budgets. Guardrails ensure CRAQ's reliability isn't compromised: production chains maintain high-availability tiers while cold chains move to lower-cost regions. Publish ESG reports based on dirty-duration and workload metrics to demonstrate efficiency improvements, aligning with DDIA's emphasis on monitoring and operational governance." 

## sample_acceptable - Minimum Acceptable
"Use energy-efficient hardware, move cold CRAQ chains to greener/cheaper regions when SLA allows, implement cost transparency with guardrails so savings don’t reduce reliability—following DDIA's operational governance." 

## common_mistakes - Watch Out For
- Cost cutting without reliability guardrails
- Ignoring sustainability metrics or regulatory reporting
- No governance or accountability
- Not referencing DDIA's operations/monitoring guidance

## follow_up_excellent - Depth Probe
**Question**: "How do you quantify whether moving a chain to a greener region affects latency SLOs?"
- **Looking for**: Latency modeling, instrumentation, pilot programs
- **Red flags**: Assuming no impact

## follow_up_partial - Guided Probe  
**Question**: "What incentives encourage product teams to participate in cost/carbon reduction?"
- **Hint embedded**: Chargeback, executive scorecards, shared OKRs
- **Concept testing**: Organizational levers

## follow_up_weak - Foundation Check
**Question**: "Why compare the fuel efficiency of cars in a fleet?"
- **Simplification**: Efficiency metrics analogy
- **Building block**: Measure to manage

## bar_raiser_question - L7→L8 Challenge
"Propose an industry alliance to set carbon efficiency benchmarks for distributed databases, with CRAQ as the reference implementation." 

### bar_raiser_concepts
- Industry-wide sustainability standards
- Benchmark governance
- Strategic positioning, thought leadership
- Influence beyond the company

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 6-8 min answer + 6-7 min discussion
- **Common next topics**: ESG reporting, cost governance, energy-aware scheduling
