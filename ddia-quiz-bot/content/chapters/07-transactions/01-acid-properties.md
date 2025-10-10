---
id: ch07-acid-properties
day: 1
tags: [transactions, acid, consistency, databases]
related_stories: []
---

# ACID Properties

## question
Which ACID property ensures that concurrent transactions appear to execute in some serial order, even if they actually run in parallel?

## options
- A) Atomicity
- B) Consistency
- C) Isolation
- D) Durability

## answer
C

## explanation
Isolation is the ACID property that ensures concurrent transactions don't interfere with each other. It makes each transaction appear as if it's the only one running in the system, even when multiple transactions execute simultaneously. Different isolation levels provide different guarantees about how much transactions can see of each other's uncommitted changes.

## hook
Why do most databases default to weaker isolation levels than full serializability?
