---
id: spanner-truetime-uncertainty
day: 4
tags: [spanner, truetime, clocks, uncertainty]
---

# TrueTime and Uncertainty Bounds

## question
What does Spanner’s TrueTime API return and why?

## options
- A) A single NTP-synchronized timestamp to minimize latency
- B) A logical clock value that never goes backwards
- C) An interval [earliest, latest] capturing bounded clock uncertainty
- D) A vector clock capturing causality across shards

## answer
C

## explanation
TrueTime exposes time uncertainty explicitly by returning an interval within which the real time lies, based on GPS/atomic clocks. The system reasons about correctness using these bounds rather than assuming perfectly synchronized physical clocks.

## hook
What can a database do with time intervals that it cannot safely do with a single unsynchronized timestamp?
