---
id: primary-backup-output-commit
day: 4
tags: [replication, consistency, performance]
related_stories:
  - vmware-ft
---

# Output Commit Problem

## question
When can the primary send output to external clients in VMware FT?

## options
- A) Immediately after processing the request
- B) Only after the backup acknowledges receiving the log entries

## answer
B

## explanation
The primary must delay output until the backup confirms receipt of log entries to prevent inconsistency after failover.

## hook
Why does replication add latency to client responses?
