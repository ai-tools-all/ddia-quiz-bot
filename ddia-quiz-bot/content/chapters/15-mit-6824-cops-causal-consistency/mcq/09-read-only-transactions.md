---
id: cops-read-only-transactions
day: 1
tags: [cops, cops-gt, transactions, multi-key-reads]
---

# Read-Only Transactions (COPS-GT)

## question
Why does COPS need the COPS-GT extension for consistent multi-key reads (read-only transactions)?

## options
- A) Base COPS doesn't track dependencies for get operations
- B) Without COPS-GT, reads of multiple keys might see different points in causal time (causally inconsistent snapshot)
- C) COPS-GT adds two-phase locking which base COPS lacks
- D) Base COPS cannot handle reads from multiple shards simultaneously

## answer
B

## explanation
In base COPS, individual gets return the highest visible version, but multiple gets might see causally inconsistent versions. For example: get(X) sees version X:5 (which depends on Y:3), but get(Y) sees Y:2 (older). This violates causal consistency across the read set. COPS-GT (Get Transaction) extends COPS to provide causally-consistent snapshots for multi-key reads by using a two-round protocol: first round gets versions, second round checks dependencies are satisfied. This ensures all keys in the transaction reflect a causally-consistent point in time.

## hook
How does COPS-GT's two-round protocol differ from Spanner's snapshot reads with TrueTime?
