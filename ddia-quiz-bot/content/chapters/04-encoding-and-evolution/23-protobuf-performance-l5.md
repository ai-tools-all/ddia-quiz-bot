---
id: ch04-protobuf-performance-l5
day: 23
level: L5
tags: [protobuf, performance, json, trade-offs, microservices]
related_stories: []
---

# Protobuf vs JSON Performance Trade-offs

## question
Your microservices architecture currently uses JSON over HTTP. The team proposes switching to gRPC with Protocol Buffers, claiming 10x performance improvement. When would this migration NOT be worth the cost?

## options
- A) When services communicate across the internet with external third-party clients
- B) When message size is large and bandwidth is the bottleneck
- C) When services make thousands of calls per second internally
- D) When you need strict type safety and code generation

## answer
A

## explanation
Protobuf's benefits (compact size, fast parsing, type safety) are most valuable for high-volume internal service communication. For external APIs with third-party clients, JSON's human readability, browser compatibility, and ubiquitous tooling often outweigh performance gains. The migration cost includes losing curl/browser debugging, requiring client code generation, and reduced ecosystem compatibility. This tests understanding of trade-offs between performance and operational concerns.

## hook
How do you debug a gRPC/Protobuf call in production when you can't use curl?
