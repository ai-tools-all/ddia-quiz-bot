### Goals
* Build coordination services that stay highly available despite machine failures.
* Compare ZooKeeper’s replicated log approach with chain-replication families (Chain Replication and CRAQ) for scalability and consistency.

### ZooKeeper Recap
* Leader-based service on Raft; all replicas execute writes in identical order, preserving per-client sequencing.
* Hierarchical znode namespace lets many apps share one service while isolating state.
* Followers may serve stale reads; applications must choose between freshness and load distribution.

### Coordination Primitives
* **Ephemeral znodes** vanish with client session death → easy master election/test-and-set.
* **Sequential znodes** give unique monotonic numbering for lock queues or task IDs.
* **Watches** deliver one-shot notifications; clients must re-register after each trigger.
* **Version numbers** enable optimistic concurrency (compare-and-set style deletes/writes).

### Limits of ZooKeeper Locks
* Locks only serialize access; protected state lives elsewhere, so crash recovery must repair partial updates.
* Real systems prefer idempotent operations or cleanup logic before releasing locks.

### Chain Replication Primer
* Linear chain (head → … → tail). Writes enter at head, propagate sequentially, commit when tail applies and replies.
* Reads hit the tail for latest committed value; failure handling mainly reassigns head/tail roles or replays buffered writes.
* Requires an external **configuration manager** (often ZooKeeper/Raft) to avoid split brain and publish membership.

### CRAQ (Chain Replication with Apportioned Queries)
* Objective: scale read throughput while keeping linearizable semantics.
* Each replica tracks per-object **version + clean/dirty flag**.
* Writes still traverse head→tail, leaving intermediates dirty until tail confirms commit.
* Clean replicas answer reads locally; dirty replicas forward to tail, which either serves the read or marks them clean once the commit arrives.
* Duplicate-suppression at the head handles client retries after failures.

### Operational Considerations
* Chain-style systems are sensitive to slow nodes because every write touches every link.
* Must pair fast data-plane pipelines with robust, consensus-backed configuration control.
* CRAQ shines when read load dominates and network topology keeps replicas close to readers.
