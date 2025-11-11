---
id: farm-object-id-structure
day: 4
tags: [object-id, addressing, sharding, rdma]
---

# Object ID Structure and Routing

## question
How are object identifiers structured in FaRM, and why is this structure important for RDMA operations?

## options
- A) Object IDs are simple integers used to index into a hash table on the server
- B) Object IDs combine a region number (mapping to primary/backup servers) plus a memory address, enabling direct RDMA access
- C) Object IDs are UUIDs that require a lookup service to locate the object
- D) Object IDs contain only the server hostname and port number

## answer
B

## explanation
Object IDs in FaRM encode both the region/shard number and the actual memory address where the object resides. This allows clients to route one-sided RDMA operations directly to the correct server and memory location without server CPU involvement or additional lookups, which is essential for achieving microsecond latencies.

## hook
What happens if an object is moved to a different memory location during compaction?
