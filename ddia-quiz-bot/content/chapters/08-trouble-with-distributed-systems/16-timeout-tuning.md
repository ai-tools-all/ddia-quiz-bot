---
id: ch08-timeout-tuning
day: 16
level: L5
tags: [timeouts, failure-detection, trade-offs]
related_stories: []
---

# Timeout Tuning Trade-offs

## question
What is the fundamental trade-off when choosing timeout values for failure detection in distributed systems?

## options
- A) Memory usage vs CPU usage
- B) Fast failure detection vs false positive rate
- C) Network bandwidth vs storage
- D) Consistency vs partition tolerance

## answer
B

## explanation
Short timeouts detect genuine failures quickly, minimizing downtime and allowing fast failover. However, they increase false positives - declaring healthy but slow nodes as dead (due to GC pauses, network congestion, etc.). Long timeouts reduce false positives but mean genuine failures take longer to detect, increasing recovery time. Finding the right balance requires understanding your system's latency distribution, acceptable downtime, and the cost of false positives.

## hook
Why might aggressive timeouts cause more outages than they prevent?
