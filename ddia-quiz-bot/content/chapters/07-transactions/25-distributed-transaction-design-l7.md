---
id: ch07-distributed-transaction-design-l7
day: 25
level: L7
tags: [transactions, distributed-systems, architecture, principal-engineer, scalability]
related_stories: []
---

# Distributed Transaction Architecture at Scale

## question
You're architecting a global e-commerce platform processing 1M+ orders/day across 50+ microservices (inventory, payment, shipping, fraud detection, loyalty points, tax calculation). Each service has its own database. As Principal Engineer, design a transaction coordination strategy that ensures: (1) consistency for financial operations, (2) < 2 second order processing, (3) graceful degradation when services fail, (4) ability to handle Black Friday traffic spikes (10x normal), and (5) compliance with regional data regulations. Consider the trade-offs between 2PC, saga orchestration, and event sourcing approaches.

## expected_concepts
- Hybrid transaction patterns (critical vs non-critical paths)
- Saga orchestration vs choreography trade-offs
- Compensating transaction design
- Idempotency and exactly-once semantics
- Circuit breakers and graceful degradation
- Event sourcing with CQRS
- Distributed tracing and observability
- Regional data sovereignty constraints

## answer
The solution requires a hybrid approach with different patterns for different consistency requirements:

Core Architecture: (1) **Critical path transactions** (payment, inventory) use synchronous saga orchestration with immediate compensations. Order service acts as orchestrator, maintaining state machine for each order. (2) **Non-critical path** (recommendations, analytics) use asynchronous event streaming with eventual consistency. (3) **Financial reconciliation** runs periodic batch jobs comparing event logs across services.

Transaction Patterns: Implement **stepped transactions** with checkpoints: (1) Reservation phase - soft locks on inventory/payment authorization (timeout-based release). (2) Confirmation phase - hard commits after all validations pass. (3) Compensation phase - automated rollback for any failures with dead letter queues for manual intervention.

Performance Optimization: (1) **Service mesh with circuit breakers** - fail fast on degraded services. (2) **Read replicas and caching** - reduce database load for read operations. (3) **Bulkheading** - isolate high-priority transactions from batch operations. (4) **Adaptive timeouts** - adjust based on current system load.

Scalability Strategy: (1) **Horizontal partitioning** by region/customer segment. (2) **Priority queues** for high-value customers during peak load. (3) **Graceful degradation** - disable non-essential services (recommendations) under extreme load. (4) **Pre-scaling** for known events with gradual warm-up.

Compliance Handling: (1) **Regional transaction coordinators** ensuring data doesn't cross boundaries. (2) **Encrypted event payloads** with regional key management. (3) **Audit log segregation** per jurisdiction.

Critical Insight: Don't try to solve distributed transactions with technology alone. Design business processes to be naturally partition-tolerant. For example, allow "eventual inventory" for non-critical items while maintaining strict consistency for limited-quantity items.

## hook
How would Amazon handle an order that involves marketplace sellers, Prime shipping, and international customs?

## follow_up
Your e-commerce platform acquires a competitor with a monolithic architecture and different transaction semantics. You need to operate both systems in parallel during a 2-year migration, maintaining transaction consistency across both platforms while gradually moving customers. Design the transaction bridge that handles dual writes, conflict resolution, and provides a unified view for customer service. Consider rollback scenarios if the migration fails partway.

## follow_up_answer
This requires a sophisticated transaction bridging layer with careful orchestration:

Bridge Architecture: (1) **Transaction Proxy Layer** intercepts all transaction requests, determining routing based on customer migration status. (2) **Dual-write coordinator** for customers mid-migration, maintaining consistency across both systems. (3) **Conflict resolution service** with business rules for handling divergent states.

Migration Patterns: (1) **Shadow mode** - new system processes transactions in parallel without side effects, comparing results for validation. (2) **Canary migrations** - move low-risk customer segments first (low transaction volume, simple use cases). (3) **Reversible migrations** - maintain bi-directional sync for rollback capability with version vectors tracking change origin.

Consistency Strategy: (1) **Write-through cache** as source of truth during migration, both systems read from cache but write to their respective stores. (2) **Saga decomposition** - break monolith transactions into steps matching microservice boundaries. (3) **Event replay capability** - rebuild state in either system from authoritative event log.

Dual-Write Handling: (1) **Versioned writes** with vector clocks to detect concurrent modifications. (2) **Deterministic conflict resolution** - last-write-wins for non-financial data, manual queue for financial conflicts. (3) **Compensating transaction chains** - if one system fails after both committed, run compensation in successful system.

Unified View: (1) **GraphQL federation layer** presenting consistent API regardless of backend. (2) **Real-time ETL pipelines** synchronizing data for analytics and customer service tools. (3) **Transaction genealogy tracking** - maintain lineage of transactions across systems for debugging.

Rollback Strategy: (1) **Incremental rollback** by customer segment, not all-or-nothing. (2) **State snapshot points** before major migration milestones. (3) **Parallel run period** where old system remains authoritative while validating new system.

Key Insight: The bridge isn't just technical - it's organizational. Establish clear ownership boundaries, runbooks for conflict resolution, and escalation paths. The hardest part isn't handling normal operations but managing edge cases and failures that span both systems. Build extensive monitoring and alerting for transaction divergence early.
