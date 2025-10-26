---
id: spanner-commit-wait
day: 5
tags: [spanner, commit-wait, truetime, consistency]
---

# Commit-Wait Rule

## question
What is the purpose of Spanner’s commit-wait in read-write transactions?

## options
- A) To batch multiple commits for higher throughput
- B) To ensure the chosen commit timestamp is definitely in the past before releasing locks
- C) To wait for all replicas globally to acknowledge the commit
- D) To allow followers to serve reads immediately after the write

## answer
B

## explanation
After selecting a commit timestamp using TrueTime’s latest bound, Spanner waits until TrueTime.now().earliest exceeds that timestamp. This guarantees the commit’s timestamp is in the past in real time, preserving external consistency ordering.

## hook
Why does a ~7–10ms delay buy you strong ordering guarantees across the world?
