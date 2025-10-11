---
id: zookeeper-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: zookeeper
subtopic: configuration-management
estimated_time: 8-10 minutes
---

# question_title - Dynamic Configuration Management with Zookeeper

## main_question - Core Question
"Design a dynamic configuration management system using Zookeeper where hundreds of services need to receive configuration updates within seconds. Explain your design choices and how you handle edge cases like partial updates, rollbacks, and service restarts."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Configuration Storage**: Store configs as znodes with versioning
- **Watch-Based Updates**: Services watch config nodes for changes
- **Atomic Updates**: Ensure all-or-nothing config changes
- **Service Registration**: Services register to track who's using configs

### expected_keywords
- Primary keywords: configuration, watch, atomic, versioning, rollback
- Technical terms: znodes, ephemeral nodes, sequential nodes, watches

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Staged Rollout**: Using separate paths for canary deployments
- **Validation**: Config schema validation before deployment
- **Audit Trail**: Sequential nodes for configuration history
- **Health Checks**: Ephemeral nodes to track service health
- **Two-Phase Updates**: Prepare/commit pattern for coordinated changes
- **Rate Limiting**: Preventing thundering herd on updates

### bonus_keywords
- Patterns: blue-green deployment, feature flags, circuit breakers
- Implementation: watch re-registration, bulk operations, transactions
- Operations: rollback mechanisms, monitoring, alerting

## sample_excellent - Example Excellence
"I'd design a hierarchical configuration system with these components:

Structure:
- `/config/prod/serviceA/v001` - versioned configuration data
- `/config/prod/serviceA/current` - pointer to active version
- `/services/serviceA/instances/` - ephemeral nodes for live instances

Update Flow:
1. Write new config to versioned node (e.g., v002)
2. Update 'current' pointer atomically
3. Watches on 'current' fire, notifying all instances
4. Services sync+read to get latest, validate, then apply

Edge Cases:
- **Partial Updates**: Use Zookeeper's multi-operation transactions to ensure atomicity
- **Rollbacks**: Simply update 'current' pointer to previous version; all services automatically revert
- **Service Restarts**: On startup, read 'current' + set watch; ephemeral node tracks liveness
- **Validation Failures**: Services write to `/config/prod/serviceA/v002/failed/instance-xyz` to report issues

Optimizations:
- Implement exponential backoff for watch re-registration to prevent storms
- Use batching for services starting simultaneously
- Add '/config/prod/serviceA/staging' for canary testing
- Keep last 10 versions for quick rollback
- Use sequential ephemeral nodes to implement rolling updates

This design leverages Zookeeper's strengths: strong consistency for critical config data, watches for immediate propagation, and ephemeral nodes for service discovery."

## sample_acceptable - Minimum Acceptable
"Store configurations as znodes in Zookeeper with paths like `/config/serviceName`. Services watch these nodes and get notified of changes. For updates, write new config atomically and let watches trigger service reloads. Handle rollbacks by keeping previous versions and switching back when needed. On restart, services read current config and set up watches again."

## common_mistakes - Watch Out For
- Not handling watch re-registration after fires
- Ignoring the thundering herd problem
- No versioning or rollback strategy
- Missing atomic update considerations
- Not considering service health/readiness

## follow_up_excellent - Depth Probe
**Question**: "How would you handle a scenario where a configuration change requires coordinated updates across multiple services with dependencies?"
- **Looking for**: Two-phase commit, dependency ordering, canary deployments, rollback coordination
- **Red flags**: Not considering failure scenarios, no ordering mechanism

## follow_up_partial - Guided Probe  
**Question**: "What happens if 100 services all try to read the new configuration at the exact same moment after a watch fires?"
- **Hint embedded**: Load spike on Zookeeper
- **Concept testing**: Understanding scalability concerns

## follow_up_weak - Foundation Check
**Question**: "If you update a configuration file on disk, how do services know to reload it? How does Zookeeper make this better?"
- **Simplification**: Push vs pull, polling vs notifications
- **Building block**: Event-driven advantages

## bar_raiser_question - L4→L5 Challenge
"Your configuration system needs to support A/B testing where 10% of services get configuration variant A and 90% get variant B, with the ability to dynamically adjust percentages and promote winners. Design this using Zookeeper primitives."

### bar_raiser_concepts
- Consistent hashing for deterministic assignment
- Sequential nodes for ordered percentage boundaries
- Atomic updates for percentage changes
- Service registration with deterministic IDs
- Monitoring and metrics integration

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 5-6 min discussion
- **Common next topics**: Service mesh, feature flags, GitOps, configuration as code
