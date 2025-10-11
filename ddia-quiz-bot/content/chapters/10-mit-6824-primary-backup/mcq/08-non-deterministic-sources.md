---
id: primary-backup-non-determinism
day: 8
tags: [replication, determinism, vmware-ft]
related_stories:
  - vmware-ft
---

# Sources of Non-Determinism

## question
Which is NOT a source of non-determinism that VMware FT must handle?

## options
- A) Arithmetic operations like addition
- B) Random number generation
- C) Timing variations in device interrupts

## answer
A

## explanation
Basic arithmetic operations are deterministic. Random numbers, time reads, and interrupt timing are non-deterministic.

## hook
What makes execution non-deterministic in replicated systems?
