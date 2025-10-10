---
id: ch08-consensus-impossibility
day: 19
level: L5
tags: [consensus, flp-impossibility, theory]
related_stories: []
---

# Consensus in Asynchronous Systems

## question
The FLP impossibility result proves that consensus cannot be deterministically solved in asynchronous systems with one faulty process. How do practical systems like Raft and Paxos achieve consensus anyway?

## options
- A) They violate the FLP theorem
- B) They use randomization or timeouts to achieve probabilistic termination
- C) They only work in synchronous networks
- D) They don't actually achieve consensus

## answer
B

## explanation
FLP proves deterministic consensus is impossible in asynchronous systems with failures. Practical algorithms circumvent this by sacrificing deterministic termination. They use timeouts (making the system partially synchronous) or randomization to ensure consensus terminates with high probability. While they can't guarantee termination in bounded time, they maintain safety (never disagree) and achieve liveness (eventual termination) under reasonable assumptions about timing and failures.

## hook
How do real systems achieve "impossible" consensus?
