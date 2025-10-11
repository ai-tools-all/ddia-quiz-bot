---
id: craq-subjective-L5-003
type: subjective
level: L5
category: baseline
topic: craq
subtopic: configuration-failure
estimated_time: 8-10 minutes
---

# question_title - Configuration Manager Failure Drill

## main_question - Core Question
"Conduct a failure drill where the configuration manager loses quorum mid-failover. Describe the safe recovery procedure for CRAQ and relate it to the failure detection and fencing strategies from DDIA's coordination chapters." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Quorum Loss Response**: Freeze chain membership changes; halt new writes
- **Fencing Tokens**: Prevent stale managers or nodes from rejoining without authorization
- **Recovery Steps**: Restore quorum, replay pending configuration changes, revalidate dirty replicas
- **DDIA Connection**: Mirrors Zookeeper-style fencing and lease expiration logic

### expected_keywords
- Primary keywords: quorum loss, fencing, lease, failover, recovery runbook
- Technical terms: epoch, monotonic sequence, consensus, coordination service

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Client Behaviour**: Retry with exponential backoff, surface degraded mode
- **Audit Trail**: Persist configuration events for reconciliation
- **Testing Playbooks**: Regular chaos drills, automated verification
- **Operational Risk**: Danger of split brain if fencing ignored

### bonus_keywords
- Implementation: Raft term, Zookeeper ephemeral znode, witness node, shadow manager
- Scenarios: network partition of manager cluster, software upgrade gone wrong, clock skew
- Trade-offs: downtime vs risk, manual vs automated intervention

## sample_excellent - Example Excellence
"If the configuration manager loses quorum while promoting a new head, we immediately freeze write traffic—just like DDIA warns in the coordination chapter. Nodes keep their last valid epoch; any attempt to operate with expired fencing tokens is rejected. Recovery involves restoring the manager's quorum, replaying the configuration intent log, verifying that no replica processed a higher epoch, then issuing a new configuration with incremented epoch. Only after all replicas acknowledge the new view do we resume writes. This is the same lease/fencing discipline Zookeeper uses to prevent stale leaders from acting." 

## sample_acceptable - Minimum Acceptable
"When the configuration manager loses quorum we halt writes, prevent stale nodes from acting with fencing tokens, restore the manager, replay the last intent, and only then resume. This follows the lease and fencing advice from DDIA's coordination chapter." 

## common_mistakes - Watch Out For
- Allowing chain to self-elect during manager outage
- Forgetting to mention fencing tokens or epochs
- Restarting writes without reconciling pending dirty replicas
- Not tying approach back to DDIA coordination practices

## follow_up_excellent - Depth Probe
**Question**: "How would you detect if any replica acted on the half-completed failover during the outage?"
- **Looking for**: Audit logs, epoch comparison, dirty entry reconciliation
- **Red flags**: Blind trust, lack of verification

## follow_up_partial - Guided Probe  
**Question**: "What client-facing behaviour reassures users during the freeze?"
- **Hint embedded**: Graceful errors, status page, read-only mode
- **Concept testing**: Operational empathy

## follow_up_weak - Foundation Check
**Question**: "If the referee leaves mid-game, why should play stop until a new referee arrives?"
- **Simplification**: Configuration manager analogy
- **Building block**: Prevent chaos by waiting for authority

## bar_raiser_question - L5→L6 Challenge
"Design guardrails that automatically roll back a partially applied CRAQ reconfiguration if the manager quorum disappears, using ideas from DDIA's transactional sagas." 

### bar_raiser_concepts
- Compensating actions, sagas
- Event sourcing of configuration changes
- Automatic rollback vs manual intervention
- Ensuring atomicity of topology updates

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Automated recovery, chaos engineering, audit infrastructure
