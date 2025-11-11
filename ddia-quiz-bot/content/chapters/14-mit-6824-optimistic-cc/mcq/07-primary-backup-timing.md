---
id: farm-primary-backup-timing
day: 4
tags: [replication, primary-backup, commit, durability]
---

# Primary-Backup Replication Timing

## question
When does FaRM send COMMIT-BACKUP messages to backup replicas during the commit protocol?

## options
- A) Before sending LOCK messages to primaries
- B) After all primaries respond yes to LOCK, but before sending COMMIT-PRIMARY
- C) After sending COMMIT-PRIMARY messages to primaries
- D) Only after the entire transaction has committed and locks are released

## answer
C

## explanation
FaRM follows a specific commit order: first LOCK messages validate and lock primaries, then COMMIT-PRIMARY tells primaries to apply writes and increment versions, and finally COMMIT-BACKUP propagates updates to backups for fault tolerance. This ordering ensures primaries commit first, establishing the authoritative state before replicating to backups.

## hook
Why is it safe to send COMMIT-BACKUP after COMMIT-PRIMARY rather than waiting for backup acknowledgments before responding to the client?
