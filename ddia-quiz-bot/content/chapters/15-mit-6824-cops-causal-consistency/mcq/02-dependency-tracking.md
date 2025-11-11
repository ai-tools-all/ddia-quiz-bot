---
id: cops-dependency-tracking
day: 1
tags: [cops, context, dependencies, design]
---

# Dependency Tracking Mechanism

## question
In COPS, what information does a client maintain in its "context" to track causal dependencies?

## options
- A) A log of all SQL queries executed by the client
- B) A set of key-version pairs observed during gets
- C) A list of timestamps for each operation using GPS clocks
- D) A merkle tree hash of all data read by the client

## answer
B

## explanation
Each COPS client maintains a context consisting of key-version pairs that accumulate causally prior reads. When the client performs a get, the returned key-version pair is added to the context. When issuing a put, the client sends the current context as a dependency list, explicitly encoding "this write depends on having seen X version 2 and Y version 4." This allows the system to preserve causal ordering without using physical clocks or complex vector clocks.

## hook
Why is tracking key-version pairs sufficient to capture causal dependencies across multiple keys?
