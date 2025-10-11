---
id: primary-backup-split-brain
day: 5
tags: [replication, split-brain, consistency]
related_stories:
  - vmware-ft
---

# Split-Brain Prevention

## question
How does VMware FT prevent split-brain during network partition?

## options
- A) Using heartbeat timeouts to detect failures
- B) Using atomic test-and-set on shared storage
- C) Letting clients decide which replica to send requests to

## answer
B

## explanation
VMware FT uses atomic test-and-set on shared storage (VMFS) to ensure only one replica becomes primary during partition.

## hook
What happens when primary and backup can't communicate?
