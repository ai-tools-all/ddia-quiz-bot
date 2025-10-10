---
id: ch07-deadlock-detection
day: 12
tags: [transactions, deadlocks, locking, timeout]
related_stories: []
---

# Deadlock Detection

## question
Transaction A holds lock on row 1, wants row 2. Transaction B holds lock on row 2, wants row 1. How do databases typically resolve this deadlock?

## options
- A) Wait indefinitely until one transaction voluntarily releases
- B) Abort the transaction that has done less work
- C) Abort both transactions immediately
- D) Automatically split the transactions

## answer
B

## explanation
Databases detect deadlocks using a wait-for graph or timeout mechanism. When detected, they choose a victim transaction to abort (often the one with less work done or lower priority), allowing the other to proceed. The aborted transaction's changes are rolled back and it typically retries. This is preferable to indefinite waiting or aborting both transactions.

## hook
How does PostgreSQL's deadlock_timeout setting balance between quick detection and false positives?
