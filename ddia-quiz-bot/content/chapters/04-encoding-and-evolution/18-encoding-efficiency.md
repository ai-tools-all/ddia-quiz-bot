---
id: ch04-encoding-efficiency
day: 18
tags: [encoding, performance, binary-formats, efficiency]
related_stories: []
---

# Encoding Efficiency Trade-offs

## question
Your service sends millions of messages per second between microservices. What is the MOST important factor in choosing an encoding format?

## options
- A) Human readability for debugging
- B) Compatibility with web browsers
- C) Encoding/decoding performance and message size
- D) Alphabetical sorting of field names

## answer
C

## explanation
At high message volumes, encoding/decoding performance and message size directly impact throughput, latency, and network costs. While debugging is important, at millions of messages per second, the performance impact of inefficient encoding dominates. Binary formats like Protobuf, Thrift, or Avro significantly outperform JSON/XML in both speed and size.

## hook
How much money are you spending on network bandwidth to transmit field names in JSON?
