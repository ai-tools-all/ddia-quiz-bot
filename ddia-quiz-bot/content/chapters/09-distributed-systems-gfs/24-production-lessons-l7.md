---
id: ch09-production-lessons-l7
day: 1
tags: [gfs, production-systems, architectural-philosophy, L7]
related_stories: []
---

# GFS Production Lessons (L7)

## question
GFS demonstrated that "tailoring the design to the workload" was more important than general-purpose solutions. What broader architectural principle does this represent for building large-scale systems?

## options
- A) Always use industry-standard protocols for compatibility
- B) Design for specific use cases and co-design with applications
- C) Build the most flexible system possible to handle any workload
- D) Prioritize correctness over performance in all scenarios

## answer
B

## explanation
GFS's success came from co-designing the storage system with its applications (like MapReduce). By accepting relaxed consistency and optimizing for large sequential operations, it achieved massive scale economically. This principle - that specialized systems outperform general-purpose ones at scale - influenced many later systems like Bigtable, Spanner, and specialized ML infrastructure.

## hook
Should you build one perfect system for everyone or ten specialized systems for different needs?
