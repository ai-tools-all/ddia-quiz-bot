---
id: ch06-multi-tenant-partitioning-l7
day: 25
level: L7
tags: [partitioning, multi-tenancy, architecture, principal-engineer, scalability]
related_stories: []
---

# Multi-Tenant Partitioning Strategy

## question
You're designing a B2B SaaS platform expecting 100K tenants ranging from 10 MB to 100 GB of data each. The top 1% of tenants generate 60% of revenue and traffic. As Principal Engineer, design a partitioning strategy that optimizes for: (1) tenant isolation for enterprise customers, (2) cost efficiency for small tenants, (3) ability to scale individual tenants, and (4) preventing noisy neighbor problems. Consider operational complexity and migration paths between tiers.

## expected_concepts
- Hybrid partitioning strategies (shared vs dedicated)
- Tenant classification and tier assignment
- Resource pools and quality of service
- Dynamic tenant migration between tiers
- Cost optimization vs performance isolation
- Monitoring and capacity planning per tier
- Noisy neighbor detection and mitigation

## answer
The solution requires a tiered partitioning strategy with automatic tier progression based on tenant characteristics:

Architecture: (1) **Pool-based partitioning** for small tenants (<1GB): Multiple tenants per partition with logical isolation, hash-partitioned by tenant_id for load distribution. (2) **Dedicated partitions** for medium tenants (1-10GB): Single tenant owns multiple partitions, enabling independent scaling. (3) **Dedicated clusters** for enterprise tenants (>10GB or SLA requirements): Complete physical isolation with custom configuration.

Tier assignment: Dynamic scoring based on data size, query rate, revenue tier, and compliance requirements. Implement automatic promotion/demotion with hysteresis to prevent flapping.

Key innovations: (1) **Virtual partitions** - logical abstraction over physical partitions enabling seamless migration. (2) **Resource tokens** - rate limiting per tenant preventing noisy neighbors in shared pools. (3) **Shadow writes** during migration - write to both old and new partitions during transition, switch reads after verification.

Operational model: Centralized schema management with tenant-specific migrations. Per-tier monitoring with different SLAs. Cost allocation model transparent to sales for pricing decisions.

Critical insight: Design for tenant mobility from day one. The biggest operational challenge isn't the initial partitioning but moving tenants between tiers as they grow or change usage patterns without downtime or data loss.

## hook
How do you handle a small tenant that suddenly goes viral and needs immediate scaling?

## follow_up
A major enterprise client with dedicated infrastructure demands data residency in specific regions for different subsidiaries, but they want unified global analytics across all their data. They also require the ability to dynamically move subsidiaries between regions based on changing regulations. How do you modify your partitioning strategy to support this while maintaining query performance and data consistency?

## follow_up_answer
This requires extending the partitioning strategy with geo-awareness and cross-region federation:

Multi-region architecture: (1) **Regional partition pools** with subsidiary-level partition assignment. Each subsidiary gets dedicated partitions within their required region. (2) **Global metadata layer** tracking subsidiary-to-region mappings with versioning for audit compliance. (3) **Cross-region read replicas** for analytics with field-level encryption for PII during replication.

Query federation: Implement a **smart query router** that decomposes queries into regional sub-queries for operational data (respecting residency) and assembles results. For analytics, maintain **materialized views** in a central analytics region with anonymized/aggregated data that complies with all regional regulations.

Dynamic migration: (1) **Dual-write period** - new region shadows writes while historical data migrates. (2) **Incremental catchup** using change data capture (CDC) to minimize switchover window. (3) **Atomic cutover** using distributed transaction to update metadata and switch primary region. (4) **Retention policies** for old region data based on compliance requirements.

Performance optimization: Regional query caches with cross-region invalidation. Predictive pre-warming of new region before migration. Query result caching at federation layer for repeated analytics.

Key insight: The real complexity isn't in the data movement but in maintaining a consistent global view while respecting regional boundaries. Build the abstraction layer between logical tenant structure and physical partition placement early - this flexibility becomes invaluable as regulations and requirements evolve.
