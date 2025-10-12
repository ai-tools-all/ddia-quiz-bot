### Goals & Context
* Build a geo-replicated key-value store supporting local reads and writes at each data center, avoiding cross-datacenter coordination on the critical path to minimize latency and improve fault tolerance.
* Deliver stronger consistency than eventual consistency while maintaining performance comparable to systems like Dynamo, but without the anomalies that confuse application developers.
* Enable straightforward application code (like photo upload followed by list update) to work correctly without explicit synchronization barriers.

### The Eventual Consistency Problem
* Pure eventual consistency (Strawman 1) allows local reads and writes with asynchronous replication via queues between data centers, but produces anomalies where causally related operations appear out of order at remote sites.
* The canonical photo example fails: inserting a photo then adding it to a list may cause readers at other data centers to see the list entry before the photo itself arrives, breaking application logic.
* Using physical clocks or even Lamport clocks for version numbers helps order concurrent updates to the same key, but does not preserve cross-key causal relationships.

### Strawman 2: Explicit Synchronization
* Introducing a sync barrier (similar to fsync) forces a write to wait until it propagates to all data centers before returning, ensuring subsequent operations see it.
* This solves ordering but sacrifices the goal of local writes—clients must wait for cross-datacenter round-trips, negating the performance benefits and reducing fault tolerance.
* Systems like Spanner and Facebook's Memcache already incur similar write latencies through Paxos or primary-site coordination, so this approach offers no improvement.

### COPS Design: Dependency Tracking
* Each client maintains a context—a set of key-version pairs observed during gets—that accumulates causally prior reads as the client progresses through operations.
* When issuing a put, the client sends not just the key and value but also the current context as a dependency list, explicitly encoding "this write depends on having seen X version 2 and Y version 4."
* The local shard server immediately acknowledges the put and assigns a new version number, allowing the client to proceed without waiting; the server then asynchronously replicates the put with its dependencies to remote data centers.

### Maintaining Causal Order at Replicas
* Remote shard servers receiving a put do not immediately make it visible to local gets; instead, they wait until all dependencies in the attached list are satisfied (i.e., those versions exist and are visible locally).
* This deferred visibility ensures that any client reading the new version will also see (or have already seen) all causally prior writes, preserving the application's expected ordering.
* Gets return the highest visible version; the client library adds the returned key-version pair to the context, propagating causality forward into subsequent puts.

### Causal Consistency Guarantees
* COPS enforces causal consistency: if operation A causally precedes B (through same-client sequencing or cross-client read-write chains), then all clients observe A before B.
* Dependencies are transitive, so a write depending on a read of an earlier write inherits all of that earlier write's dependencies, forming causal chains across multiple clients and keys.
* This model sits between eventual consistency (too weak, anomalies common) and linearizability (too strong, requires coordination), offering a practical middle ground for geo-replication.

### Trade-offs & Limitations
* Cascading dependency waits can delay visibility at remote sites if upstream dependencies themselves are waiting, potentially causing long stalls even without failures.
* Network partitions blocking dependency propagation will stall dependent writes indefinitely, sacrificing availability for consistency in the affected key sets.
* Lamport clocks and last-writer-wins for concurrent updates to the same key mean true conflicts still discard one update, requiring application-level conflict resolution if needed.
* The system assumes cooperative use and does not address security or malicious clients; it also lacks transactional multi-key updates in the base design (though COPS-GT extends it).

### Key Takeaways
* COPS demonstrates that causal consistency is achievable with local reads and writes by embedding dependency metadata in puts and deferring visibility at replicas until dependencies are satisfied.
* The client library transparently tracks context and attaches dependencies, making causal ordering automatic for application code and eliminating common eventual-consistency anomalies.
* This approach offers a compelling design point for geo-replicated systems where low latency and intuitive semantics matter more than strict linearizability, though cascading delays and partition sensitivity remain practical concerns.
