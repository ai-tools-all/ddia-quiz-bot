---
id: primary-backup-correlated-failures
day: 6
tags: [replication, failures, availability]
related_stories:
  - datacenter-disasters
---

# Correlated Failures

## question
Which scenario would likely defeat primary-backup replication?

## options
- A) Random hardware failure in one server
- B) Earthquake affecting the entire datacenter
- C) Planned OS upgrade on the primary with backup standing by

## answer
B

## explanation
Correlated failures like natural disasters affecting all replicas in the same location defeat replication's protection.

## hook
Why does replica placement matter for availability?
