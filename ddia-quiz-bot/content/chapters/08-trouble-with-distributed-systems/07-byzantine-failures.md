---
id: ch08-byzantine-failures
day: 7
tags: [byzantine-failures, fault-tolerance, distributed-systems]
related_stories: []
---

# Byzantine Failures

## question
What distinguishes a Byzantine failure from a regular node failure in distributed systems?

## options
- A) Byzantine failures are more common
- B) Byzantine failures involve nodes sending arbitrary or malicious incorrect messages
- C) Byzantine failures are easier to detect
- D) Byzantine failures only occur in blockchain systems

## answer
B

## explanation
Byzantine failures occur when nodes don't just fail by stopping (fail-stop) but instead send arbitrary, contradictory, or malicious messages to different nodes. This could be due to bugs, hardware corruption, or malicious actors. Most distributed systems assume non-Byzantine (crash-stop) failures because Byzantine fault tolerance requires much more complex algorithms and typically needs 3f+1 nodes to tolerate f Byzantine nodes.

## hook
When can a faulty node lie to other nodes in the system?
