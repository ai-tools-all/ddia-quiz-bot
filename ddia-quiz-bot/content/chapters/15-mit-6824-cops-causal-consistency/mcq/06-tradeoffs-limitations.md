---
id: cops-tradeoffs-limitations
day: 1
tags: [cops, tradeoffs, limitations, partitions]
---

# Trade-offs and Limitations

## question
What is a key limitation of COPS regarding network partitions?

## options
- A) COPS requires GPS clocks which fail during network partitions
- B) Network partitions blocking dependency propagation will stall dependent writes indefinitely
- C) COPS automatically switches to eventual consistency during partitions
- D) Partitions cause COPS to lose all dependency metadata

## answer
B

## explanation
Network partitions blocking dependency propagation will stall dependent writes indefinitely at affected replicas, sacrificing availability for consistency in those key sets. Additionally, cascading dependency waits can delay visibility at remote sites if upstream dependencies themselves are waiting, potentially causing long stalls even without failures. COPS uses Lamport clocks and last-writer-wins for concurrent updates to the same key, meaning true conflicts still discard one update.

## hook
Is COPS an AP or CP system according to the CAP theorem, and why?
