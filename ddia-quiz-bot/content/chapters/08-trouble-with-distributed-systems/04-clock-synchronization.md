---
id: ch08-ntp-accuracy
day: 4
tags: [clocks, ntp, synchronization]
related_stories: []
---

# Clock Synchronization Accuracy

## question
What is the typical accuracy of NTP clock synchronization within a datacenter?

## options
- A) Nanoseconds
- B) Microseconds  
- C) Milliseconds
- D) Seconds

## answer
C

## explanation
Within a datacenter, NTP can typically synchronize clocks to within tens of milliseconds. Over the public internet, accuracy is usually worse (tens to hundreds of milliseconds). This level of accuracy is insufficient for many distributed systems algorithms that require precise ordering of events, which is why logical clocks like Lamport timestamps are often used instead.

## hook
Can you rely on synchronized clocks to order events in distributed systems?
