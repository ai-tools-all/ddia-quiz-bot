---
id: ch05-conflict-resolution
day: 9
tags: [replication, conflicts, multi-leader, CRDTs, last-write-wins]
related_stories: []
---

# Write Conflict Resolution

## question
Two users concurrently edit the same document in different datacenters. What is the most dangerous conflict resolution strategy?

## options
- A) Last-write-wins (LWW) based on timestamp
- B) Merge the changes using CRDTs
- C) Keep all versions and let the application resolve
- D) Designate one datacenter as primary for conflict resolution

## answer
A

## explanation
Last-write-wins (LWW) is dangerous because it silently discards data - one write will overwrite the other based solely on timestamp, potentially losing important updates. This is especially problematic with clock skew between datacenters. CRDTs or application-level resolution are safer as they preserve information and handle conflicts explicitly.

## hook
How do collaborative editing apps like Google Docs handle concurrent edits without losing data?
