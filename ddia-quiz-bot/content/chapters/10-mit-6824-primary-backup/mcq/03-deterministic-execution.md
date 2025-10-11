---
id: primary-backup-determinism
day: 3
tags: [replication, determinism, consistency]
related_stories:
  - vmware-ft
---

# Deterministic Execution

## question
What must VMware FT log to handle non-deterministic events?

## options
- A) Only the occurrence of the event
- B) Both the event value and exact instruction count when it occurred
- C) Nothing; the backup infers timing from its own execution

## answer
B

## explanation
Non-deterministic events must be logged with their values AND precise timing (instruction count) for accurate replay.

## hook
Why is determinism crucial for replicated state machines?
