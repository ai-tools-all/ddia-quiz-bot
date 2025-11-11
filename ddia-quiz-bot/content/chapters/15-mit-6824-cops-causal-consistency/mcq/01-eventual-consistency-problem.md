---
id: cops-eventual-consistency-problem
day: 1
tags: [cops, eventual-consistency, causal-consistency, anomalies]
---

# Eventual Consistency Problem

## question
What is the key problem with pure eventual consistency (Strawman 1) that the COPS paper addresses?

## options
- A) It doesn't provide fault tolerance across data centers
- B) It has high write latency due to cross-datacenter coordination
- C) Causally related operations can appear out of order at remote sites, causing application anomalies
- D) It cannot handle network partitions between data centers

## answer
C

## explanation
Pure eventual consistency allows local reads and writes with asynchronous replication, but produces anomalies where causally related operations (like inserting a photo then adding it to a list) appear out of order at remote data centers. A reader might see the list entry before the photo itself arrives, breaking application logic. COPS solves this by tracking and enforcing causal dependencies.

## hook
How does the photo example illustrate why eventual consistency is insufficient for real applications?
