---
id: ch08-partial-failures
day: 12
tags: [partial-failures, distributed-systems, reliability]
related_stories: []
---

# Partial Failures

## question
What makes partial failures particularly challenging in distributed systems compared to single-node systems?

## options
- A) They happen more frequently
- B) They're deterministic
- C) Some parts of the system work while others fail, creating inconsistent states
- D) They're easier to debug

## answer
C

## explanation
In a single-node system, failures tend to be total - the whole system either works or crashes. In distributed systems, partial failures mean some nodes/networks/operations succeed while others fail, creating complex inconsistent states. A request might succeed on some replicas but fail on others, or a multi-step operation might complete partially. This non-deterministic behavior makes distributed systems much harder to reason about and debug.

## hook
Why is a half-failed system worse than a fully failed system?
