---
id: primary-backup-failover-trigger
day: 11
tags: [replication, availability, fault-tolerance]
related_stories:
  - vmware-ft
---

# Failover Trigger

## question
What event typically causes the backup to take over as primary in a primary-backup system?

## options
- A) Completion of a scheduled checkpoint window
- B) Timeout on heartbeat messages from the current primary
- C) Manual operator approval through a control console

## answer
B

## explanation
Backups monitor periodic heartbeat messages from the current primary and assume leadership when those heartbeats stop, indicating a primary failure.

## hook
How do replicas agree that the primary has actually failed?
