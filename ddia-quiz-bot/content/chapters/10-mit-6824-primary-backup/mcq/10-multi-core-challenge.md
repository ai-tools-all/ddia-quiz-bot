---
id: primary-backup-multicore
day: 10
tags: [replication, multi-core, performance]
related_stories:
  - vmware-ft
---

# Multi-Core Challenges

## question
Why did VMware FT initially support only single-core VMs?

## options
- A) Multi-core VMs consume too much memory
- B) Multi-core execution introduces memory access interleaving non-determinism

## answer
B

## explanation
Multi-core systems have non-deterministic memory access patterns that are extremely difficult to log and replay deterministically.

## hook
Why is multi-core replication more complex than single-core?
