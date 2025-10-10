---
id: ch07-materializing-conflicts
day: 10
tags: [transactions, materializing-conflicts, concurrency-control, phantoms]
related_stories: []
---

# Materializing Conflicts

## question
To prevent phantom reads in a meeting room booking system, you pre-create rows for every room and time slot, marking them as available/booked. What is this technique called?

## options
- A) Pessimistic locking
- B) Materializing conflicts
- C) Snapshot isolation
- D) Index locking

## answer
B

## explanation
Materializing conflicts involves creating explicit records for all possible states that could conflict, turning phantom prevention into a simpler update conflict problem. Instead of checking for the absence of records (which creates phantoms), you update existing records (which can be properly locked). While effective, it can lead to table bloat and isn't always practical.

## hook
When is materializing conflicts a good solution versus using serializable isolation?
