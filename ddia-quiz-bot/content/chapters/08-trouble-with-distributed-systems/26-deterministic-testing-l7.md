---
id: ch08-deterministic-testing-l7
day: 26
level: L7
tags: [testing, determinism, distributed-systems, principal-engineer]
related_stories: []
---

# Deterministic Testing of Non-Deterministic Systems

## question
As Principal Engineer, you're seeing 10% of production incidents caused by race conditions and timing-dependent bugs that never appeared in testing. Your test environments can't reproduce the failures even with the same inputs. Design a testing strategy that can systematically find these distributed systems bugs before production.

## expected_concepts
- Deterministic simulation testing (like FoundationDB)
- Model checking and formal methods (TLA+)
- Fault injection and chaos engineering
- Property-based testing
- Linearizability checking
- Jepsen-style testing

## answer
The fundamental problem: distributed systems bugs emerge from the interaction of non-determinism (message ordering, timing, failures) with system logic. Traditional testing explores a tiny fraction of possible interleavings. The solution requires multiple complementary approaches:

(1) **Deterministic Simulation**: Build a simulation layer that controls all non-determinism - time, network, thread scheduling. Run millions of test scenarios with different interleavings, failure patterns, and timings. FoundationDB found thousands of bugs this way. Key insight: make the non-determinism controllable and reproducible.

(2) **Model Checking**: Write formal specifications in TLA+ for critical algorithms. Model checkers exhaustively explore all possible states to find violations. Amazon uses this for S3, DynamoDB, and EBS. Catches design bugs, not just implementation bugs.

(3) **Controlled Chaos**: Unlike random chaos engineering, systematically explore failure combinations. Use fault injection that can reproduce exact failure sequences. The key is deterministic replay - when you find a bug, you can reproduce it reliably.

(4) **Property-Based Testing**: Instead of testing specific scenarios, test invariants that should always hold (linearizability, serializability). Generate random operations and verify the system maintains its guarantees.

The architectural requirement: systems must be designed for deterministic testing from day one. This means explicit handling of all non-determinism, mockable time, and observable internal state.

## hook
How do you repeatedly test something that never happens the same way twice?

## follow_up
Your deterministic testing found a subtle bug that only occurs when: (1) Network partition happens, (2) During leader election, (3) While a transaction is committing, (4) And a node experiences a GC pause, (5) In a specific ordering. It took 10 million simulation runs to find. The team argues this is too rare to fix given the engineering cost. How do you evaluate whether to fix it?

## follow_up_answer
The rarity argument is flawed due to scale effects: a one-in-10-million event happens 100 times per day at billion-request scale. The evaluation framework: (1) **Blast Radius**: Does it cause data loss, corruption, or just temporary unavailability? Data corruption bugs must be fixed regardless of rarity, (2) **Cascade Potential**: Can this trigger failure of other systems? Rare bugs that cause cascading failures are system-killers, (3) **Recovery Complexity**: Can operators easily recover, or does it require manual data repair? Complex recovery makes rare bugs extremely expensive, (4) **Confidence Erosion**: Knowing about unfixed race conditions creates technical debt and fear around the codebase. Teams become hesitant to make changes. The meta-insight: in distributed systems, "rare" is relative to your scale. At sufficient scale, every possible race condition will happen. The question isn't whether to fix it, but whether to fix it now or after it causes an incident. Generally, the cost of fixing known bugs is lower than the cost of debugging them under pressure during an outage.
