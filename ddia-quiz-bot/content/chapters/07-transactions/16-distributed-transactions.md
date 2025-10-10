---
id: ch07-distributed-transactions
day: 16
tags: [transactions, distributed-systems, two-phase-commit, coordination]
related_stories: []
---

# Distributed Transactions

## question
Why are distributed transactions across multiple databases often avoided in modern microservices architectures?

## options
- A) They're impossible to implement correctly
- B) Performance overhead and reduced availability from coordination
- C) Modern databases don't support them
- D) Security concerns with data sharing

## answer
B

## explanation
Distributed transactions using protocols like two-phase commit (2PC) introduce significant performance overhead due to coordination messages and synchronous waiting. They also reduce availability - if any participant fails, the entire transaction blocks. Modern architectures often prefer eventual consistency with compensation/saga patterns over distributed ACID transactions for better performance and availability.

## hook
How do payment systems handle money transfers without distributed transactions between banks?
