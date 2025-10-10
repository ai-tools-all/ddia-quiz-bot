---
id: ch07-write-skew
day: 8
tags: [transactions, write-skew, isolation, concurrency]
related_stories: []
---

# Write Skew

## question
Two doctors simultaneously check if at least one doctor is on-call. Both see two doctors on-call, so both remove themselves from the on-call list. Now zero doctors are on-call. What is this problem called?

## options
- A) Lost update
- B) Dirty write
- C) Write skew
- D) Phantom read

## answer
C

## explanation
This is write skew - a transaction reads some data, makes a decision based on that data, then writes based on that decision. When two transactions do this concurrently, they can each make decisions based on a premise that the other transaction invalidates. Write skew is harder to prevent than lost updates and requires serializable isolation or explicit locking of the read set.

## hook
How would you prevent the "on-call doctor" problem in a real hospital scheduling system?
