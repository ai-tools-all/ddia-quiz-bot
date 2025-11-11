---
id: farm-validate-optimization
day: 2
tags: [occ, read-only, optimization, performance]
---

# Validate Optimization

## question
How does the VALIDATE optimization improve performance for read-only transactions in FaRM?

## options
- A) It allows read-only transactions to skip version checking entirely
- B) It replaces expensive LOCK messages with one-sided RDMA reads to check versions and lock bits, avoiding log appends
- C) It caches object values on clients to eliminate all network round trips
- D) It uses pessimistic locking instead of optimistic validation

## answer
B

## explanation
For objects that a transaction reads but doesn't modify, VALIDATE uses one-sided RDMA reads of object headers instead of sending LOCK messages to primaries. This avoids expensive log appends and primary CPU processing while still ensuring serializability by checking that versions haven't changed and lock bits aren't set.

## hook
Can a read-only transaction using VALIDATE still abort? Under what conditions?
