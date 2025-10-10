---
id: ch06-compound-partition-keys
day: 17
tags: [partitioning, compound-keys, schema-design]
related_stories: []
---

# Compound Partition Keys

## question
A messaging application uses (user_id, timestamp) as a compound partition key. What query pattern does this optimize for?

## options
- A) Finding all messages between any two users
- B) Retrieving a user's messages within a time range
- C) Searching messages by content globally
- D) Counting total messages in the system

## answer
B

## explanation
With (user_id, timestamp) as a compound key, user_id determines the partition and timestamp provides ordering within that partition. This optimizes for retrieving a specific user's messages within a time range - the query goes to a single partition and can efficiently scan the time-ordered messages. Cross-user queries would need to hit multiple partitions.

## hook
How does Cassandra's compound primary keys enable efficient time-series queries?
