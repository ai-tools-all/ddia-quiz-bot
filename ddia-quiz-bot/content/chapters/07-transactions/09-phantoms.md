---
id: ch07-phantoms
day: 9
tags: [transactions, phantoms, isolation, write-skew]
related_stories: []
---

# Phantom Reads

## question
A transaction checks that no meetings are booked in Room A from 2-3 PM, then books the room. Another concurrent transaction does the same check and also books Room A for 2-3 PM. What causes this double-booking?

## options
- A) Lost update on the booking record
- B) Phantom read - new records appearing that match a query condition
- C) Dirty read of uncommitted bookings
- D) Incorrect timestamp comparison

## answer
B

## explanation
This is a phantom read problem. The first transaction's query "find bookings for Room A from 2-3 PM" returns an empty set. But while it's executing, another transaction inserts a booking that would have matched that query - a "phantom" record. This leads to write skew. Solutions include predicate locks or serializable isolation.

## hook
How does MySQL's InnoDB engine use gap locks to prevent phantom reads?
