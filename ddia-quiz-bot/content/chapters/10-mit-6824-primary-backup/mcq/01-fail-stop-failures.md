---
id: primary-backup-fail-stop
day: 1
tags: [replication, failures, fault-tolerance]
related_stories:
  - vmware-ft
---

# Fail-Stop Failure Model

## question
What type of failures can primary-backup replication typically handle?

## options
- A) Software bugs that cause incorrect calculations
- B) Hardware failures that cause the server to stop executing
- C) Byzantine failures caused by actively malicious replicas

## answer
B

## explanation
Primary-backup replication handles fail-stop failures where servers stop cleanly, not bugs that would affect both replicas identically.

## hook
Can replication protect against all types of failures?
