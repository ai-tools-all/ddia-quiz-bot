---
id: craq-subjective-L7-003
type: subjective
level: L7
category: baseline
topic: craq
subtopic: merger-integration
estimated_time: 12-15 minutes
---

# question_title - CRAQ Integration Strategy After a Merger

## main_question - Core Question
"Your company acquires a firm running a leaderless Dynamo-style system. Develop an integration strategy that harmonizes both stacks into a unified CRAQ-based platform, preserving customer SLAs and leveraging DDIA's insights on heterogeneous system evolution." 

## core_concepts - Must Mention (60%)
- **Migration Phasing**: Dual-write or shadow-read strategy to minimize downtime
- **Consistency Harmonization**: Translate Dynamo's eventual consistency semantics into CRAQ clean/dirty model
- **Customer Communication**: SLA maintenance, staged cutovers, transparency
- **Technical Debt Management**: Sunset legacy features, provide compatibility APIs

### expected_keywords
- Primary keywords: migration, dual-write, shadow read, compatibility, SLA
- Technical terms: consistency bridge, conflict resolution, idempotency, phased rollout

## peripheral_concepts - Nice to Have (40%)
- **Risk Mitigation**: Canary migrations, kill switches, rollback plans
- **Data Quality**: Reconciliation pipelines, conflict repair
- **Org Structure**: Joint architecture board, knowledge transfer
- **Tooling**: Compatibility SDKs, data validators

### bonus_keywords
- Implementation: CDC bridge, contract tests, feature flags, traffic shaping
- Scenarios: high-availability accounts, multi-region customers, compliance obligations
- Trade-offs: speed vs stability, platform simplification vs feature parity

## sample_excellent - Example Excellence
"Adopt a dual-write strategy: Dynamo continues serving customers while we replicate writes into CRAQ using CDC. We build a consistency bridge that maps Dynamo's vector clocks into CRAQ's clean/dirty metadata, resolving conflicts before tail commit. Shadow reads validate parity; once error budgets show confidence, we cut traffic over region by region, offering compatibility APIs. Communication plans assure customers of SLA continuity, and joint architecture boards manage technical debt. This aligns with DDIA's advice on phased evolution of heterogeneous systems." 

## sample_acceptable - Minimum Acceptable
"Use dual writes and shadow reads to migrate Dynamo workloads into CRAQ, reconcile conflicts before tail commit, cut over gradually with strong customer communication, and provide compatibility APIs—matching DDIA's phased evolution approach." 

## common_mistakes - Watch Out For
- Big-bang migration without staging
- Ignoring conflict semantics translation
- No plan for customer communication or SLAs
- Not referencing DDIA's heterogeneous evolution guidance

## follow_up_excellent - Depth Probe
**Question**: "How do you handle Dynamo-specific features (e.g., per-item TTL) when moving to CRAQ?"
- **Looking for**: Emulation layer, feature parity decisions, roadmap commitments
- **Red flags**: Dropping features without mitigation

## follow_up_partial - Guided Probe  
**Question**: "Which metrics tell you it's safe to complete a regional cutover?"
- **Hint embedded**: Error rate parity, tail lag vs Dynamo latency, reconciliation backlog
- **Concept testing**: Data-driven go/no-go

## follow_up_weak - Foundation Check
**Question**: "Why move roommates gradually into a new house rather than all at once?"
- **Simplification**: Phased migration analogy
- **Building block**: Manage disruption

## bar_raiser_question - L7→L8 Challenge
"Design contractual guarantees for customers that bridge eventual-to-strong consistency transitions without service credits." 

### bar_raiser_concepts
- Contract negotiation, SLA evolution
- Customer trust, risk sharing
- Business and technical alignment
- Migration guarantees

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 6-8 min answer + 6-7 min discussion
- **Common next topics**: Migration governance, dual-write complexity, customer success
