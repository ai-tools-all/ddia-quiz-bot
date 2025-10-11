---
id: primary-backup-system-goal
day: 12
tags: [replication, consistency, design-goals]
related_stories:
  - vmware-ft
---

# Essence of Primary-Backup

## question
What core property does the primary-backup design in MIT 6.824 aim to deliver?

## options
- A) Infinite horizontal scalability regardless of workload diversity
- B) Continued service availability with consistent state despite a single replica failure
- C) Elimination of the need for client-side retries and acknowledgements

## answer
B

## explanation
Primary-backup replication keeps a hot standby that mirrors the primary's operations so that service can continue with consistent state after a single node failure.

## hook
Why does consistency still matter even when availability is the priority?
