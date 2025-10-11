---
id: craq-subjective-L6-002
type: subjective
level: L6
category: baseline
topic: craq
subtopic: failure-injection
estimated_time: 9-12 minutes
---

# question_title - CRAQ Failure Injection Program

## main_question - Core Question
"Design a chaos engineering program for CRAQ that validates its guarantees under compounded failures (e.g., tail crash + configuration manager leader change + Kafka outage for CDC). Relate each experiment to DDIA's failure modes and observability recommendations." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Failure Matrix**: Enumerate combined scenarios (replica crash, manager failover, downstream outage)
- **Hypotheses and Metrics**: Define expected behaviour (e.g., write halt, dirty duration spike) tied to DDIA metrics
- **Automation & Safeguards**: Controlled chaos with blast-radius containment and rollback
- **Learning Loop**: Feed findings into runbooks and design improvements

### expected_keywords
- Primary keywords: chaos experiment, hypothesis, blast radius, observability
- Technical terms: failure injection, game day, rollback, SLOs, dirty duration

## peripheral_concepts - Nice to Have (40%)
- **Experiment Scheduling**: Rolling across regions, non-peak windows
- **Data Integrity Checks**: Compare ledger snapshots pre/post events
- **Cross-Team Coordination**: Align with downstream consumers (Kafka)
- **Documentation**: Postmortem templates, risk acceptance

### bonus_keywords
- Implementation: fault injection proxy, config toggles, canary region, anomaly detection
- Scenarios: network partition, disk corruption, configuration drift
- Trade-offs: experimentation frequency vs stability

## sample_excellent - Example Excellence
"I'd build a failure matrix combining tail crashes, configuration manager leadership changes, and CDC sink outages. For each scenario we define hypotheses: writes halt within 2 s when tail unreachable (per DDIA's CAP discussion), configuration manager fencing tokens prevent split brain, CDC backlog stays below threshold with graceful degradation. We automate injections via feature flags, monitor dirty duration, tail lag, and outbox growth, and rehearse rollback steps. Findings feed back into runbooks and alert tuning." 

## sample_acceptable - Minimum Acceptable
"Plan chaos experiments combining tail crash, config manager failover, and CDC outage, define expected behaviours (writes halt, fencing works, CDC backlog manageable), run them with monitoring on dirty duration and tail lag, and update runbooks—following DDIA's failure testing advice." 

## common_mistakes - Watch Out For
- Running single-failure tests only
- Lacking hypotheses tied to metrics
- Ignoring downstream integrations during chaos
- No rollback plan or documentation

## follow_up_excellent - Depth Probe
**Question**: "How would you ensure experiments don't violate customer SLAs while still providing signal?"
- **Looking for**: Canary regions, time-boxing, staged rollouts, kill switches
- **Red flags**: Running full-scale experiments in production without guardrails

## follow_up_partial - Guided Probe  
**Question**: "What is the success metric for an experiment where CDC is down but CRAQ stays healthy?"
- **Hint embedded**: Outbox backlog, no loss, eventual catch-up time
- **Concept testing**: Targeted validation

## follow_up_weak - Foundation Check
**Question**: "Why practice fire drills when the building isn't burning?"
- **Simplification**: Chaos rehearsal analogy
- **Building block**: Preparedness

## bar_raiser_question - L6→L7 Challenge
"Develop an enterprise governance policy that mandates cross-service chaos experiments, integrating CRAQ, streaming, and cache tiers using DDIA's holistic operations model."

### bar_raiser_concepts
- Enterprise resiliency programs
- Cross-service coordination
- Policy and governance
- Continuous improvement loops

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: Resilience testing, governance, cross-functional drills
