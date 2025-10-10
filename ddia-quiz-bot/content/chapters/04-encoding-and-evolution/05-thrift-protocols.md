---
id: ch04-thrift-protocols
day: 5
tags: [thrift, binary-protocol, compact-protocol]
related_stories: []
---

# Thrift Binary Protocols

## question
What is the difference between Thrift's BinaryProtocol and CompactProtocol?

## options
- A) BinaryProtocol is human-readable, CompactProtocol is not
- B) CompactProtocol uses variable-length integers and bit packing for more compact encoding
- C) BinaryProtocol is faster but uses more space
- D) CompactProtocol doesn't support nested structures

## answer
B

## explanation
Thrift's CompactProtocol uses variable-length integers and bit packing to achieve more compact encoding compared to BinaryProtocol. BinaryProtocol uses fixed-width integers which are simpler but less space-efficient. Both support the same data structures.

## hook
How much bandwidth could you save by switching from BinaryProtocol to CompactProtocol?
