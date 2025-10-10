---
id: ch08-process-pause-causes
day: 14
tags: [process-pauses, scheduling, distributed-systems]
related_stories: []
---

# Causes of Process Pauses

## question
Besides garbage collection, what can cause a running process to pause unexpectedly in a distributed system?

## options
- A) Only hardware failures
- B) Virtual machine suspension, page faults, CPU throttling, or context switches
- C) Network delays only
- D) Database locks only

## answer
B

## explanation
Process pauses can occur for many reasons beyond GC: VM suspension (in virtualized environments), page faults (when memory is swapped to disk), CPU throttling (resource limits), context switches (OS scheduling other processes), SIGSTOP signals, or even laptop lid closure. Any of these can cause a node to pause for seconds or minutes, making it appear dead to other nodes, potentially triggering incorrect failure detection and split-brain scenarios.

## hook
What happens when your cloud provider suspends your VM for live migration?
