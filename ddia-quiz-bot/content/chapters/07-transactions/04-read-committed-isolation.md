---
id: ch07-read-committed-isolation  
day: 4
tags: [transactions, isolation, read-committed, dirty-reads]
related_stories: []
---

# Read Committed Isolation

## question
Under read committed isolation level, which phenomenon is prevented?

## options
- A) Dirty reads and dirty writes
- B) Lost updates
- C) Write skew
- D) Phantom reads

## answer
A

## explanation
Read committed isolation prevents dirty reads (reading uncommitted data from other transactions) and dirty writes (overwriting uncommitted data). However, it doesn't prevent lost updates, write skew, or phantom reads - these require stronger isolation levels. Read committed is the default isolation level in many databases like PostgreSQL and Oracle.

## hook
Why is "read committed" the default isolation level in PostgreSQL instead of stronger guarantees?
