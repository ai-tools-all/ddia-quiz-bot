---
id: occ-why-rdma-occ
day: 1
tags: [occ, rdma, locking, design]
---

# Why OCC with RDMA

## question
Why does the RDMA-based system use optimistic concurrency control (OCC) instead of traditional lock-based concurrency control?

## options
- A) One-sided RDMA operations can directly acquire server locks efficiently
- B) One-sided RDMA doesn’t involve server CPUs, making centralized lock coordination infeasible; OCC defers checks to commit time
- C) Locks provide lower latency than version validation in this design
- D) OCC is required to support geo-replication across data centers

## answer
B

## explanation
One-sided RDMA reads/writes bypass server CPUs, so coordinating lock acquisition is not practical. OCC allows clients to read optimistically and validate at commit using version numbers and lock bits.

## hook
How does removing server CPUs from the read path change your concurrency control choices?
