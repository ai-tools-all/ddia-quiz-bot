---
id: ch05-read-after-write-consistency
day: 5
tags: [replication, consistency, read-after-write, user-experience]
related_stories: []
---

# Read-After-Write Consistency

## question
A user updates their profile on a social media site. When they immediately refresh the page, they don't see their changes. What consistency guarantee is missing?

## options
- A) Eventual consistency
- B) Read-after-write consistency
- C) Monotonic reads
- D) Consistent prefix reads

## answer
B

## explanation
Read-after-write consistency (also called read-your-writes) guarantees that users will always see their own writes immediately. Without this guarantee, a user might write to the leader but then read from a follower that hasn't yet received the update, making it appear as if their write was lost.

## hook
How would you implement read-after-write consistency for user profile updates?
