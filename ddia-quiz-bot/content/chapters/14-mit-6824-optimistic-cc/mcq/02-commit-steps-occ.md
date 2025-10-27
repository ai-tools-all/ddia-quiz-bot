---
id: occ-commit-steps
day: 2
tags: [occ, commit, 2pc, version, lock-bit]
---

# OCC Commit Steps

## question
Which sequence best describes the commit path for a read-write transaction in this OCC system?

## options
- A) Validate read versions via RDMA, broadcast writes to all replicas, then set lock bits
- B) Send LOCK to primaries (check versions and set lock bit); if all yes, COMMIT to primaries (apply writes, bump version, clear lock); then COMMIT to backups
- C) Acquire global locks first, then read data, then write directly to backups
- D) Write to backups first to ensure durability, then attempt to lock primaries

## answer
B

## explanation
At commit, the coordinator locks each primary while verifying versions. If any check fails, abort. Otherwise it applies writes on primaries (incrementing versions, clearing locks) and propagates to backups for fault tolerance.

## hook
Why set the lock bit before applying writes?
