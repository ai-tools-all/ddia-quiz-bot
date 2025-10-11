---
id: craq-subjective-L7-002
type: subjective
level: L7
category: baseline
topic: craq
subtopic: regulatory-strategy
estimated_time: 12-15 minutes
---

# question_title - CRAQ Regulatory and Governance Blueprint

## main_question - Core Question
"Formulate a regulatory strategy for CRAQ operating in jurisdictions with conflicting data residency, auditing, and privacy requirements (e.g., GDPR, HIPAA, PCI). Leverage DDIA's data governance discussions to define architectural, operational, and organizational controls." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Data Residency Controls**: Region-specific chains, sovereign data plane, policy-driven routing
- **Auditing & Compliance**: Immutable logs, cryptographic attestation, audit workflows
- **Privacy & Access**: Fine-grained access control, data minimization, retention policies
- **Organizational Governance**: Compliance steering committee, shared responsibility model

### expected_keywords
- Primary keywords: data residency, compliance, auditing, governance, privacy
- Technical terms: sovereign region, cryptographic log, RBAC, retention policy

## peripheral_concepts - Nice to Have (40%)
- **Incident Response**: Regulatory breach handling, notification timelines
- **Vendor Risk Management**: Third-party audits, certifications
- **Change Management**: Controlled schema evolution with backward compatibility
- **Customer Transparency**: Compliance SLAs, dashboards

### bonus_keywords
- Implementation: data localization, masking, anonymization, attestation service, policy engine
- Scenarios: cross-border replication, right-to-be-forgotten requests, compliance audit
- Trade-offs: global consistency vs residency, compliance cost vs agility

## sample_excellent - Example Excellence
"Deploy CRAQ chains per jurisdiction with policy-driven routing to enforce residency; global metadata remains in a compliance-aware control plane. Every tail commit emits an auditable record signed with Merkle proofs. Implement fine-grained RBAC and data retention policies aligned with GDPR and HIPAA, with automated workflows for deletion requests. Establish a compliance steering committee overseeing change management and maintain certifications (SOC2, PCI). This reflects DDIA's guidance on governance and data lifecycle management." 

## sample_acceptable - Minimum Acceptable
"Use region-specific CRAQ chains with policy routing, keep auditable tail logs, enforce access controls and retention policies, and set up a compliance governance group—aligning with DDIA's data governance advice." 

## common_mistakes - Watch Out For
- Ignoring conflicting regulations or assuming one-size-fits-all
- No governance structure or accountability
- Missing deletion/right-to-be-forgotten processes
- Not referencing DDIA's data governance principles

## follow_up_excellent - Depth Probe
**Question**: "How do you certify cryptographic logs for auditors without leaking sensitive data?"
- **Looking for**: Merkle proofs, selective disclosure, zero-knowledge options
- **Red flags**: Sharing raw logs indiscriminately

## follow_up_partial - Guided Probe  
**Question**: "Which teams must be involved in the compliance steering committee?"
- **Hint embedded**: Legal, security, SRE, product, data governance
- **Concept testing**: Cross-functional leadership

## follow_up_weak - Foundation Check
**Question**: "Why can't you use the same house rules when renting apartments in different countries?"
- **Simplification**: Regulatory difference analogy
- **Building block**: Local compliance nuance

## bar_raiser_question - L7→L8 Challenge
"Propose a certification program that allows customers to plug custom compliance policies into CRAQ without forking the platform." 

### bar_raiser_concepts
- Compliance extensibility, policy plugins
- Certification, ecosystem governance
- Customer empowerment vs platform control
- Standardization vs customization

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 6-8 min answer + 6-7 min discussion
- **Common next topics**: Policy engines, compliance automation, cross-border operations
