---
id: ch07-saga-pattern
day: 19
tags: [transactions, saga-pattern, compensation, microservices]
related_stories: []
---

# Saga Pattern

## question
In microservices, a saga breaks a distributed transaction into local transactions with compensating actions. What happens when step 3 of a 5-step saga fails?

## options
- A) All steps are retried from the beginning
- B) Compensating transactions undo steps 2 and 1 in reverse order
- C) The entire saga is abandoned without cleanup
- D) Only step 3 is retried indefinitely

## answer
B

## explanation
The saga pattern handles failures by executing compensating transactions for completed steps in reverse order. When step 3 fails, the saga runs compensating transactions for steps 2 and 1 (in that order) to undo their effects. This achieves eventual consistency without distributed transactions. Each step must be idempotent and have a defined compensation action.

## hook
How does Uber handle multi-service operations like ride booking without distributed transactions?
