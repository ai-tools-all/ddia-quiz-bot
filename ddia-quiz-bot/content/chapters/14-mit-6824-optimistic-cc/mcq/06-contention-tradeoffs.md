---
id: occ-contention-tradeoffs
day: 6
tags: [occ, performance, contention]
---

# Contention Trade-offs

## question
When does OCC typically deliver the best performance compared to lock-based schemes in this system?

## options
- A) Under high contention on hot keys where many transactions conflict
- B) Under low contention with read-heavy workloads where conflicts are rare
- C) When all transactions are long and update many objects
- D) Only when deployed across multiple data centers

## answer
B

## explanation
OCC shines when conflicts are infrequent: reads run with one-sided RDMA and validation succeeds. With heavy contention, aborts and retries can dominate.

## hook
What application patterns increase abort rates in OCC?
