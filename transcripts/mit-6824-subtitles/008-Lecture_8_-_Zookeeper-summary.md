### Goals & Context
* Review linearizability as the reference model for “strong consistency,” then examine how ZooKeeper trades some guarantees for performance.
* Understand ZooKeeper as a generalized coordination service layered on Zab (a Raft-like replication protocol) and why it became a popular cloud primitive.

### Linearizability Refresher
* Linearizability demands that every completed operation appear in a single, real-time-respecting order; no read may observe a value older than the latest completed write.
* Histories can be validated either by constructing a total order that satisfies these constraints or by deriving graph cycles that prove inconsistency.
* Duplicate or retransmitted RPCs expand an operation’s real-time interval from first send to final response; servers must replay prior results to remain linearizable.

### ZooKeeper’s Service Model
* ZooKeeper exposes a hierarchical namespace of znodes storing small data items, akin to a filesystem shared by many clients.
* Sessions authenticate clients and anchor ephemeral znodes that vanish on session expiry, enabling lease-like behavior.
* Sequential znodes append monotonically increasing suffixes, offering simple global counters and FIFO queues without transactions.

### Consistency Guarantees
* Writes flow through a single leader and are replicated via Zab; they are linearizable and persisted before acknowledgement.
* Reads may be served by any replica; followers can lag, so ZooKeeper only promises per-client FIFO ordering plus “read-your-writes” when reads hit the leader.
* Clients needing fresh data invoke `sync` to force their session’s next read to observe the leader’s latest zxid.

### Zab Replication & Performance
* Zab mirrors Raft: the leader orders log entries, followers append and acknowledge, and majority confirmation commits entries.
* Because every write still hits the leader, adding replicas improves fault tolerance but not write throughput; however, offloading reads to followers scales read-heavy workloads nearly linearly.
* Crash recovery replays the log to rebuild the tree, preserving zxid ordering so that clients see monotonically increasing versions.

### Coordination Primitives
* Ephemeral znodes underpin master election: clients attempt to create an ephemeral node; the survivor detects leadership via watches and zxid comparisons.
* Sequential znodes implement fair locks or barriers by creating ordered children and watching the predecessor to avoid herding.
* Watches are one-shot notifications that trigger on znode changes; clients must re-register and tolerate races between notification delivery and new state.

### Operational Considerations
* Sessions require regular heartbeats; network partitions or long GC pauses can expire sessions and delete ephemeral state, so application logic must tolerate abrupt leadership changes.
* Metadata is small (KBs per znode); ZooKeeper is ill-suited for large payloads but excels at storing configuration hashes, shard maps, and coordination metadata.
* Deployments typically use odd-sized ensembles (3, 5, 7) in a single region to balance failure tolerance and latency; cross-region stretches worsen follower staleness.

### Key Takeaways
* ZooKeeper packages consensus into a reusable service, letting applications reuse battle-tested primitives instead of embedding Raft themselves.
* The system embodies a pragmatic consistency spectrum: linearizable writes plus mostly-fresh reads, with explicit APIs (`sync`, watches) for applications that need stronger guarantees.
* By combining hierarchical state, ephemeral resources, and ordered znodes, ZooKeeper supplies the building blocks for leader election, configuration management, and distributed locking at scale.
