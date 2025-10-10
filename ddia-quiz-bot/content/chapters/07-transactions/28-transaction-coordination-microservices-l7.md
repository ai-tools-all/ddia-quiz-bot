---
id: ch07-transaction-coordination-microservices-l7
day: 28
level: L7
tags: [transactions, microservices, saga, principal-engineer, distributed-systems]
related_stories: []
---

# Transaction Coordination in Microservices

## question
You're redesigning Uber's ride booking flow spanning 20+ microservices: user, driver, pricing, routing, payment, fraud, insurance, promotion, notification, tracking, dispatch, surge, safety, rating, support, analytics, billing, tax, compliance, and partner services. A single ride involves complex orchestration: (1) match rider with driver, (2) calculate dynamic pricing, (3) authorize payment, (4) track ride progress, (5) process payment on completion, (6) handle mid-ride changes. Design a transaction coordination system that handles: partial failures gracefully, regulatory requirements (some regions require upfront payment), scale to 1M concurrent rides, maintain <500ms booking latency, and support debugging when things go wrong. Consider how to handle edge cases like driver cancellation after payment authorization.

## expected_concepts
- Saga orchestration vs choreography at scale
- State machine modeling for complex flows
- Compensating transactions with time bounds
- Distributed tracing and observability
- Circuit breakers and bulkheading
- Event sourcing for audit and replay
- Idempotency and exactly-once semantics
- Timeout and retry strategies with exponential backoff

## answer
The solution requires a sophisticated orchestration platform with intelligent state management:

Architecture Foundation: (1) **Hierarchical Saga Orchestration** - master saga for ride lifecycle, sub-sagas for payment, matching, tracking. Each saga has its own state machine with clear compensation paths. (2) **Event-sourced backbone** - all state changes as events, enabling replay and debugging. (3) **Service mesh with sidecars** - standardized communication, retry, circuit breaking.

State Machine Design: (1) **Ride states**: REQUESTED → MATCHING → MATCHED → PAYMENT_AUTHORIZING → EN_ROUTE → ARRIVED → IN_PROGRESS → COMPLETING → COMPLETED. (2) **Parallel states** for independent operations (fraud check + route calculation). (3) **Timeout transitions** - automatic state progression on timeouts (e.g., driver no-show). (4) **Regional variants** - different state machines for different regulatory requirements.

Transaction Coordination: (1) **Optimistic resource locking** - soft-lock driver for 30 seconds during matching. (2) **Two-phase booking** - reserve price and driver, confirm after payment auth. (3) **Compensating transaction matrix** - predefined compensation for each state+failure combination. (4) **Saga execution engine** - handles retries, timeouts, compensation triggering.

Failure Handling: (1) **Graceful degradation paths** - if surge pricing fails, fall back to standard pricing. (2) **Partial success handling** - complete ride even if rating service is down. (3) **Time-boxed compensations** - if refund fails, queue for manual processing. (4) **Dead letter queues per service** - capture and retry failed operations.

Performance Optimization: (1) **Predictive pre-warming** - pre-authorize payment during matching. (2) **Async non-critical paths** - analytics, notifications don't block ride flow. (3) **Regional saga executors** - reduce latency with geographic distribution. (4) **Caching with TTL** - cache driver locations, surge multipliers.

Edge Case Handling: (1) **Driver cancellation after payment**: immediate refund initiation, re-enter matching with priority. (2) **Mid-ride rerouting**: new payment authorization, differential pricing. (3) **Lost connectivity**: offline mode with eventual reconciliation. (4) **Regulatory holds**: escrow pattern for regions requiring upfront payment.

Observability: (1) **Distributed tracing** with correlation IDs across all services. (2) **Saga visualization dashboard** - real-time view of ride state progression. (3) **Anomaly detection** - alert on unusual state transitions or timing. (4) **Replay capability** - reproduce issues by replaying event stream.

Critical Insight: The complexity isn't in the happy path but in the combination of failures. Build a failure scenario matrix and test combinations. The system should be designed for debuggability from day one - when a ride goes wrong at 2 AM, you need to quickly understand which service in which state caused the issue.

## hook
How does Uber handle a ride where the driver's app crashes mid-trip but the rider's app is still tracking?

## follow_up
Your ride-hailing platform expands into food delivery, sharing 60% of the microservices but with different transaction semantics (multi-stop delivery, restaurant preparation time, thermal tracking). You need to support both models simultaneously, with some services handling both rides and deliveries. Design a unified transaction framework that accommodates both domains while allowing independent evolution. Consider crossover scenarios like a driver switching from rides to delivery mid-shift, and maintaining separate SLAs for each business line.

## follow_up_answer
This requires an extensible transaction framework with domain-specific adaptations:

Unified Framework Architecture: (1) **Abstract Saga Protocol** - base saga interface with domain-specific implementations. (2) **Service capability discovery** - services advertise which operations they support for each domain. (3) **Transaction context propagation** - headers indicating domain, priority, SLA requirements. (4) **Polymorphic service behavior** - same service endpoint, different logic based on context.

Domain Modeling: (1) **Shared core operations** - payment, user, notification, tracking used by both. (2) **Domain-specific operations** - kitchen integration for delivery, surge pricing for rides. (3) **Transaction inheritance hierarchy** - BaseTransaction → RideTransaction/DeliveryTransaction. (4) **Cross-domain state machines** - driver state spans both domains with transitions.

Multi-Domain Orchestration: (1) **Domain-aware saga router** - routes to appropriate orchestrator based on transaction type. (2) **Composite transactions** - delivery can spawn multiple ride-like segments. (3) **Resource pools per domain** - separate driver pools with crossover capability. (4) **Priority-based scheduling** - hot food delivery may preempt ride matching.

Service Adaptation: (1) **Strategy pattern per domain** - payment service has different auth strategies. (2) **Feature flags per domain** - enable/disable features independently. (3) **Domain-specific SLAs** - same service, different timeout/retry policies. (4) **Separate metric namespaces** - independent monitoring per business line.

Crossover Handling: (1) **Driver state synchronization** - atomic transition between ride/delivery modes. (2) **Earnings aggregation** - unified wallet despite different commission structures. (3) **Cross-domain reputation** - bad delivery rating affects ride matching priority. (4) **Hybrid operations** - package delivery during ride (with passenger consent).

Evolution Strategy: (1) **API versioning per domain** - rides on v2 while delivery moves to v3. (2) **Canary deployments per domain** - test delivery changes without affecting rides. (3) **Domain-specific feature teams** - independent development with shared platform team. (4) **Contract testing** - ensure shared services maintain compatibility with both domains.

SLA Differentiation: (1) **Dynamic timeout adjustment** - food delivery has stricter time constraints. (2) **Differentiated retry policies** - more aggressive retries for perishable deliveries. (3) **Domain-specific circuit breakers** - delivery failures don't affect ride availability. (4) **Resource reservation** - guarantee minimum capacity per business line.

Key Insight: Don't try to force-fit different business models into the same transaction pattern. Build a platform that supports multiple transaction personalities while sharing common infrastructure. The art is in identifying what's truly shared vs. what seems similar but has different semantics. A "driver" in rides has different constraints than a "driver" in delivery - model these as different specializations of a base concept rather than forcing them to be identical.
