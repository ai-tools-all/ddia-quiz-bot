---
id: primary-backup-replication-approach
day: 2
tags: [replication, state-machine, architecture]
related_stories:
  - vmware-ft
---

# Replication Approaches

## question
Which replication approach does VMware FT use?

## options
- A) State transfer - sending memory snapshots
- B) Replicated state machine - sending operations/inputs

## answer
B

## explanation
VMware FT uses replicated state machine, sending input events to the backup rather than memory/state snapshots.

## hook
How does VMware FT keep replicas synchronized?
