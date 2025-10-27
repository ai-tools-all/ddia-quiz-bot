---
id: occ-read-only-validation
day: 3
tags: [occ, validation, read-only, rdma]
---

# Read-Only Validation

## question
How do read-only transactions validate correctness in this OCC design?

## options
- A) They acquire read locks on all objects they access
- B) They send LOCK messages to primaries and set lock bits temporarily
- C) They perform one-sided RDMA reads of object headers to re-check version numbers and lock bits before returning results
- D) They skip validation because reads never conflict

## answer
C

## explanation
Read-only transactions avoid locks. They re-read object headers (version, lock bit) via RDMA to ensure no concurrent writer changed or locked the objects since the initial read, otherwise they abort and retry.

## hook
Why does header re-checking preserve serializability without locks?
