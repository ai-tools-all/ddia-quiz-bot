---
id: ch04-schema-governance-l7
day: 25
level: L7
tags: [schema-evolution, governance, organizational-design, principal-engineer]
related_stories: []
---

# Schema Evolution and Organizational Governance

## question
Your company has 200+ microservices using Protocol Buffers, and you're observing an increasing number of production incidents caused by incompatible schema changes. Teams claim they follow the "additive changes only" rule, yet breakages still occur. As a Principal Engineer, diagnose the root cause and propose a systemic solution.

## expected_concepts
- Schema registry and centralized validation
- CI/CD integration for compatibility checks
- Organizational Conway's Law effects
- Producer-consumer contract testing
- Breaking change detection automation
- Cultural issues around schema ownership

## answer
The root cause is typically not technical ignorance but systemic organizational failures: (1) No central schema registry means no automated compatibility verification before deployment, (2) Teams lack visibility into who consumes their schemas, (3) Default values and optional fields create semantic breaking changes that pass syntactic validation, (4) Rolling deployments without proper sequencing (deploy consumers before producers for backward compatibility, producers before consumers for forward compatibility). 

Solution architecture: Implement a schema registry (Confluent Schema Registry, Buf Schema Registry) as the single source of truth, with CI/CD gates that reject incompatible changes. Enforce semantic versioning with automated compatibility checks (BACKWARD, FORWARD, FULL). Create producer-consumer contract tests. Most critically, shift culture: schemas are APIs and require the same governance - published schemas are contracts, breaking them requires major version bumps and coordinated migration plans.

Second-order insight: The real problem is treating internal service schemas as implementation details rather than public contracts. This reflects deeper organizational issues about service boundaries and ownership.

## hook
Why do schema registries reduce incidents more than better documentation?

## follow_up
Your schema registry now prevents obvious breaking changes, but you're still seeing runtime errors where services misinterpret data due to semantic mismatches (e.g., a field changed from "price in dollars" to "price in cents" without changing the field name or type). How do you systematically prevent semantic breaking changes?

## follow_up_answer
Semantic breaking changes require cultural and process solutions, not just technical ones: (1) Enforce field naming conventions that include units/semantics (price_cents, not just price), (2) Use typed quantities pattern (message Money { int64 amount_cents; string currency; }), (3) Require ADR (Architecture Decision Records) for any field semantic changes with migration plans, (4) Implement consumer-driven contract tests that verify behavior, not just schema compatibility, (5) Establish schema review process for changes to high-traffic schemas, (6) Use feature flags for gradual semantic migrations. The key insight: schemas encode syntax, but comments and team communication encode semantics - you need both in your governance model.
