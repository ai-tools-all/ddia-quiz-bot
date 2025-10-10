---
id: ch04-database-schema-migration
day: 12
tags: [databases, schema-migration, encoding]
related_stories: []
---

# Database Schema Migration

## question
What is a challenge with database schema migrations compared to application code updates?

## options
- A) Databases cannot be updated while running
- B) Most relational databases apply schema changes to the entire dataset at once, which can be slow
- C) Database schemas cannot include optional fields
- D) Schema changes require restarting all clients

## answer
B

## explanation
Most relational databases apply schema changes (like adding a column) to the entire table at once with an ALTER TABLE statement. On large tables this can mean minutes or even hours of downtime, or at least a table lock. Schema-on-read databases (like document databases) make schema changes easier - the schema is only applied when data is read.

## hook
How long would it take to add a column to your largest production table?
