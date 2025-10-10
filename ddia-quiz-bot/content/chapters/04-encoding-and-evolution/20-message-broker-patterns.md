---
id: ch04-message-broker-patterns
day: 20
tags: [message-broker, async, pub-sub, dataflow]
related_stories: []
---

# Message Broker Communication Patterns

## question
In a message broker system using publish/subscribe, what compatibility requirement exists between publishers and subscribers?

## options
- A) Publishers and subscribers must use the exact same schema version
- B) Publishers and subscribers can be updated independently if they maintain backward/forward compatibility
- C) Publishers must always be updated before subscribers
- D) Message brokers don't require compatibility between publishers and subscribers

## answer
B

## explanation
Message brokers decouple publishers from subscribers, but they still need to maintain compatibility in their message schemas. Subscribers need backward compatibility (to read messages from older publishers) and publishers need forward compatibility (since older subscribers may still be running). The broker doesn't interpret the messages, just routes them.

## hook
What happens when a publisher starts sending new fields that old subscribers don't expect?
