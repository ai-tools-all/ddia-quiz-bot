---
id: ch07-single-vs-multi-object
day: 3
tags: [transactions, multi-object, atomicity, nosql]
related_stories: []
---

# Single vs Multi-Object Transactions

## question
Many NoSQL databases only provide single-object atomicity. What is the main implication of this limitation?

## options
- A) Better performance for all operations
- B) Cannot maintain relationships between objects consistently
- C) Reduced storage requirements
- D) Simpler query syntax

## answer
B

## explanation
Without multi-object transactions, you cannot atomically update multiple related objects together. This means operations like transferring money between accounts (debiting one, crediting another) or updating a record and its index cannot be done atomically. If a failure occurs between updates, the database can be left in an inconsistent state with only partial updates applied.

## hook
How does MongoDB's multi-document transactions (added in v4.0) change its competitive position vs traditional databases?
