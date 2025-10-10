---
id: ch04-rest-advantages
day: 14
tags: [rest, rpc, api-design, services]
related_stories: []
---

# REST vs RPC Advantages

## question
What is an advantage of REST compared to RPC frameworks?

## options
- A) REST is always faster than RPC
- B) REST supports experimentation and debugging using standard tools (web browsers, curl)
- C) REST has better type safety
- D) REST uses less bandwidth

## answer
B

## explanation
REST's use of standard HTTP features means you can experiment with APIs using a web browser, curl, or other standard tools. The data formats (usually JSON) are human-readable. RPC frameworks often require custom clients and binary formats, making them harder to debug interactively. However, RPC frameworks often provide better type safety and code generation.

## hook
When was the last time you debugged an API using just curl and your eyes?
