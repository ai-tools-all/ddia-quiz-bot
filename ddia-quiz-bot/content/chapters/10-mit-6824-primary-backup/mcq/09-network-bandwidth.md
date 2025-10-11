---
id: primary-backup-bandwidth
day: 9
tags: [replication, performance, network]
related_stories:
  - vmware-ft
---

# Network Bandwidth Usage

## question
When is replicated state machine more bandwidth-efficient than state transfer?

## options
- A) When state is small and changes frequently
- B) When state is large but operations are small
- C) When clients require causal consistency across regions

## answer
B

## explanation
RSM sends small operations rather than large state snapshots, making it efficient for systems with large state but small operations.

## hook
How do you choose between state transfer and replicated state machine?
