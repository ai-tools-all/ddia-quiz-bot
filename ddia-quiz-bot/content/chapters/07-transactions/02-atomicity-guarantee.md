---
id: ch07-atomicity-guarantee
day: 2
tags: [transactions, atomicity, rollback, fault-tolerance]
related_stories: []
---

# Atomicity Guarantee

## question
What happens when a database crashes in the middle of a multi-statement transaction that has atomicity guarantee?

## options
- A) Partial changes are kept and the rest are lost
- B) All changes are rolled back as if the transaction never happened
- C) The database becomes corrupted
- D) Only the last statement's changes are preserved

## answer
B

## explanation
Atomicity means "all or nothing" - either all operations in a transaction succeed and are committed, or none of them take effect. If a crash occurs mid-transaction, the database ensures that any partial changes are rolled back during recovery, maintaining data integrity. This is typically implemented using write-ahead logging (WAL).

## hook
How does PostgreSQL's WAL (Write-Ahead Logging) ensure atomicity even during power failures?
