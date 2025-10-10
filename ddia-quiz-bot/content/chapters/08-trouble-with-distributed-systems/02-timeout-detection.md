---
id: ch08-timeout-detection
day: 2
tags: [timeouts, failure-detection, distributed-systems]
related_stories: []
---

# Timeout-Based Failure Detection

## question
Why are timeouts used as the primary mechanism for detecting failures in distributed systems?

## options
- A) Because timeouts are 100% accurate in detecting failures
- B) Because there's no other way to distinguish between a slow response and no response
- C) Because timeouts are easy to implement
- D) Because network protocols require them

## answer
B

## explanation
In an asynchronous network, it's impossible to distinguish between a node that has failed, a network that has failed, or a response that is just slow. Timeouts provide a practical heuristic: if we don't receive a response within a certain time, we assume something has failed. However, this can lead to false positives when the system is just slow.

## hook
How do you know if a remote node has crashed or is just slow?
