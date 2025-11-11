---
id: farm-subjective-L4-004
type: subjective
level: L4
category: baseline
topic: performance
subtopic: retry-backoff
estimated_time: 8-10 minutes
---

# question_title - Transaction Retry Strategies

## main_question - Core Question
"Design a transaction retry strategy for FaRM when a transaction aborts due to conflicts. Consider exponential backoff, retry limits, and fairness. Analyze a scenario where transaction T keeps aborting because it conflicts with a stream of smaller, faster transactions on the same object. How would you ensure T eventually commits without starving other transactions?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Exponential backoff: increase delay after each abort (e.g., 1ms, 2ms, 4ms, 8ms)
- Retry limit: maximum attempts before giving up or escalating
- Without backoff, concurrent retries create thundering herd, worsening contention
- Fairness problem: long transactions may starve under high contention

### expected_keywords
- exponential backoff, retry limit, thundering herd, starvation, fairness, contention

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Randomized jitter to spread retries over time
- Priority-based scheduling for long-running transactions
- Adaptive backoff based on abort rate
- Fallback to pessimistic locking for chronically failing transactions
- Application-level circuit breakers

### bonus_keywords
- jitter, priority, adaptive, pessimistic fallback, circuit breaker, livelock prevention

## sample_excellent - Example Excellence
"Implement exponential backoff with jitter: after abort, wait rand(2^attempts * base_delay). Cap at max 100ms, limit to 10 retries. For fairness, track abort count per transaction—after 5 aborts, elevate priority or acquire pessimistic lock on conflicting objects. Example: T1 (long transaction) competes with T2,T3,T4 (short). T2-T4 complete quickly. T1 aborts 5 times, then acquires pessimistic locks, forcing others to wait. This prevents starvation while maintaining OCC benefits for common case."

## sample_acceptable - Minimum Acceptable
"Use exponential backoff to delay retries after aborts. Increase wait time after each failure. Limit maximum retries to prevent infinite loops. This helps reduce contention."

## common_mistakes - Watch Out For
- Using fixed delays (causes synchronized retries)
- No retry limit (infinite loops possible)
- Ignoring fairness—long transactions may never succeed
- Not considering application-level timeouts

## follow_up_excellent - Depth Probe
**Question**: "How would you detect and break livelock where two transactions repeatedly abort each other?"
- **Looking for**: Timestamp ordering, coordinator-assigned priorities, deterministic conflict resolution

## follow_up_partial - Guided Probe
**Question**: "Why add randomized jitter to exponential backoff?"
- **Hint embedded**: Prevents synchronized retries causing repeated conflicts

## follow_up_weak - Foundation Check
**Question**: "What is the thundering herd problem in the context of retries?"
