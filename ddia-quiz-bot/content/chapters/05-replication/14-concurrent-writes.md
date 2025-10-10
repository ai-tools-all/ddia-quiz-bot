---
id: ch05-concurrent-writes
day: 14
tags: [replication, concurrent-writes, version-vectors, causality]
related_stories: []
---

# Detecting Concurrent Writes

## question
How do version vectors help detect concurrent writes in multi-replica systems?

## options
- A) They timestamp each write with wall clock time
- B) They track version numbers per replica to identify causal relationships
- C) They prevent concurrent writes from happening
- D) They automatically merge concurrent updates

## answer
B

## explanation
Version vectors maintain a version number for each replica. When a write happens, the replica increments its version number. By comparing version vectors, the system can determine if one write happened before another (causally related) or if they were concurrent (neither knew about the other). Concurrent writes require conflict resolution, while causally related writes can be ordered.

## hook
Why can't we just use timestamps to order distributed writes?
