---
id: ch04-service-dataflow
day: 13
tags: [dataflow, services, rest, rpc, microservices]
related_stories: []
---

# Dataflow Through Services

## question
What is a key difference in compatibility requirements between database dataflow and service (RPC/REST) dataflow?

## options
- A) Services don't need any compatibility guarantees
- B) Services need only backward compatibility, not forward compatibility
- C) Services often have clients and servers updated independently, requiring backward and forward compatibility
- D) Services always use the same encoding format

## answer
C

## explanation
In service dataflow (RPC/REST APIs), clients and servers are often updated independently. You might update servers first (requiring backward compatibility for old clients) or do rolling updates where old and new versions run simultaneously. Unlike databases where you control both sides, you often don't control all API clients, making compatibility critical.

## hook
Can you safely deploy a breaking API change if you control the server?
