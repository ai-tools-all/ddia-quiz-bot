---
id: ch08-clock-types
day: 3
tags: [clocks, time, distributed-systems]
related_stories: []
---

# Monotonic vs Wall-Clock Time

## question
Which type of clock should you use for measuring durations and timeouts in distributed systems?

## options
- A) Wall-clock time (time-of-day clock)
- B) Monotonic clock
- C) Either one works equally well
- D) System time

## answer
B

## explanation
Monotonic clocks are guaranteed to always move forward and are not affected by NTP adjustments or manual clock changes. They're ideal for measuring durations and timeouts. Wall-clock time can jump backwards or forwards due to NTP synchronization, making it unsuitable for measuring elapsed time but necessary for timestamps that need to correspond to actual calendar time.

## hook
Why can measuring elapsed time with System.currentTimeMillis() cause bugs?
