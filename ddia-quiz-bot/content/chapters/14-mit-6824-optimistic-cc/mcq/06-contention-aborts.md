---
id: farm-contention-aborts
day: 3
tags: [occ, contention, performance, trade-offs]
---

# OCC Contention and Aborts

## question
Under what workload conditions does FaRM's optimistic concurrency control perform poorly?

## options
- A) When transactions are geographically distributed across multiple data centers
- B) When there is high contention with many transactions writing to the same objects, causing frequent aborts
- C) When transactions only perform read operations without any writes
- D) When using RDMA hardware instead of traditional TCP/IP networking

## answer
B

## explanation
OCC performs well under low contention but struggles when many transactions conflict by writing to the same objects. Conflicting transactions detect version mismatches or locked objects at commit time and must abort, requiring application-level retries. High abort rates waste resources and reduce throughput.

## hook
How would you design an application to minimize abort rates in an OCC system?
