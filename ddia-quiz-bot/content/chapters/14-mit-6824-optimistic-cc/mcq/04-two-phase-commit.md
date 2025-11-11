---
id: farm-two-phase-commit
day: 2
tags: [2pc, commit, protocol, distributed]
---

# Two-Phase Commit with OCC

## question
In the FaRM commit protocol, what happens during the first phase (LOCK messages) sent to primaries?

## options
- A) Primaries apply the writes and increment version numbers immediately
- B) Primaries verify version numbers match expected values, atomically set lock bits, and respond yes/no
- C) Primaries forward updates to backup replicas for durability
- D) Primaries execute the transaction's read operations

## answer
B

## explanation
The LOCK phase is the prepare phase of two-phase commit. Primaries check that the object's current version matches what the transaction read (no intermediate updates), atomically set the lock bit, and respond with yes (can commit) or no (must abort). Only after all primaries respond yes does the coordinator proceed to the COMMIT-PRIMARY phase.

## hook
Why must the version check and lock bit setting be atomic?
