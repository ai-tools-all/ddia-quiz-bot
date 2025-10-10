---
id: ch04-dataflow-through-databases
day: 11
tags: [dataflow, databases, encoding, compatibility]
related_stories: []
---

# Dataflow Through Databases

## question
In database dataflow, who is the "sender" and who is the "receiver" of data?

## options
- A) The sender is the client writing data; the receiver is the database
- B) The sender is the process that writes to the database; the receiver is a future process that reads from it
- C) The sender is the database; the receiver is the client
- D) There is no sender/receiver in database dataflow

## answer
B

## explanation
In database dataflow, the process writing to the database is the sender, encoding the data. The receiver is a future process (possibly the same application at a later time, or a different application) that reads and decodes the data. This can span different versions of the application, requiring both backward and forward compatibility.

## hook
Why do you need forward compatibility when the writer and reader are the same application?
