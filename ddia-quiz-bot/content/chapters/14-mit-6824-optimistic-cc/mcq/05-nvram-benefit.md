---
id: farm-nvram-benefit
day: 3
tags: [nvram, durability, hardware, performance]
---

# Non-Volatile RAM Benefits

## question
What is the primary advantage of using non-volatile RAM (battery-backed DRAM) in FaRM?

## options
- A) It allows data to persist across power failures while providing fast in-memory access without disk writes
- B) It reduces the cost of storage hardware compared to traditional SSDs
- C) It enables geo-replication across multiple continents
- D) It automatically compresses data to save space

## answer
A

## explanation
Non-volatile RAM combines the speed of DRAM (sub-microsecond access) with durability through battery backup. This allows FaRM to achieve extremely fast transaction latencies by avoiding disk I/O entirely, while still surviving data center-wide power failures since RAM contents persist when machines restart.

## hook
What trade-off does non-volatile RAM impose compared to disk-based storage?
