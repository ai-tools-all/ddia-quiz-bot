---
id: ch07-repeatable-read
day: 17
tags: [transactions, isolation-levels, repeatable-read, consistency]
related_stories: []
---

# Repeatable Read Isolation

## question
Under repeatable read isolation, a transaction reads a value twice. Between the reads, another transaction commits a change to that value. What does the first transaction see on its second read?

## options
- A) The new committed value
- B) The original value it read first
- C) An error is thrown
- D) Null value

## answer
B

## explanation
Repeatable read isolation ensures that if a transaction reads a value, subsequent reads within the same transaction will see the same value, even if other transactions commit changes in between. This prevents non-repeatable reads (seeing different values on re-reading) but doesn't prevent phantom reads (new rows appearing that match query conditions).

## hook
Why does MySQL's repeatable read behave differently than the SQL standard definition?
