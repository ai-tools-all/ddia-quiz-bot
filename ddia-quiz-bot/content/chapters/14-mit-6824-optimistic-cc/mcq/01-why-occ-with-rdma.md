---
id: farm-why-occ-rdma
day: 1
tags: [occ, rdma, locking, design]
---

# Why OCC with RDMA

## question
Why does the FaRM system use optimistic concurrency control (OCC) instead of traditional lock-based concurrency control?

## options
- A) One-sided RDMA operations can directly acquire server locks efficiently
- B) One-sided RDMA doesn't involve server CPUs, making centralized lock coordination infeasible; OCC defers checks to commit time
- C) Locks provide lower latency than version validation in this design
- D) OCC is required to support geo-replication across multiple data centers

## answer
B

## explanation
One-sided RDMA reads/writes bypass server CPUs entirely, so coordinating lock acquisition during the read phase is not practical. OCC allows clients to read optimistically and validate at commit time using version numbers and lock bits.

## hook
How does removing server CPUs from the read path fundamentally change your concurrency control choices?
