---
id: ch04-message-passing-dataflow
day: 15
tags: [dataflow, message-queues, async, message-broker]
related_stories: []
---

# Message-Passing Dataflow

## question
What is a key advantage of message-passing dataflow (message queues) over direct RPC?

## options
- A) Message queues are always faster than RPC
- B) Message queues act as a buffer when recipient is unavailable or overloaded
- C) Message queues don't require encoding
- D) Message queues guarantee the recipient will process the message

## answer
B

## explanation
A message broker (queue) acts as a buffer, improving reliability when the recipient is temporarily unavailable or overloaded. The sender can continue working while messages accumulate in the queue. This also decouples the sender from the recipient - the sender doesn't need to know the recipient's network location. However, message delivery is usually asynchronous with no response.

## hook
What happens to your RPC calls when the downstream service crashes?
