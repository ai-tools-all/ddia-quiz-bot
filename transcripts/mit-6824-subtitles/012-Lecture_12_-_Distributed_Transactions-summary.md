### Goals & Context
* Build transactional systems that shard data across multiple servers while hiding distribution complexity from application programmers, preserving correctness despite concurrent access and partial failures.
* Support bank transfers, vote counting, or any multi-record operations where some records may reside on different machines, requiring atomicity across machine boundaries.
* Traditional database ACID guarantees must extend to distributed settings where network partitions and individual server crashes are common.

### ACID Properties
* **Atomicity** ensures that multi-step transactions either complete fully or leave no trace—failures mid-transaction must not leave half-updated state visible.
* **Consistency** (application-level invariants) is mostly out of scope; the lecture focuses on isolation and atomicity.
* **Isolation** (serializability) means concurrent transactions produce results equivalent to some serial (one-at-a-time) execution order, preventing intermediate state from leaking between transactions.
* **Durability** guarantees committed transactions survive crashes by persisting updates to non-volatile storage like disk.

### Serializability Definition
* An execution is serializable if there exists a serial order of the same transactions yielding identical results—both final database state and any printed output.
* For two transactions T1 (transfer) and T2 (audit), only two outcomes are legal: T1-then-T2 or T2-then-T1; interleaved executions producing other outputs violate serializability.
* This model is programmer-friendly: write transactions as if running alone, and the system ensures other transactions cannot interfere or observe partial updates.
* Allows true parallelism when transactions touch disjoint data, enabling performance gains in sharded systems.

### Concurrency Control: Pessimistic vs Optimistic
* **Pessimistic (locking)** schemes acquire locks before accessing data and block conflicting transactions, trading performance for safety when conflicts are frequent.
* **Optimistic** schemes execute without locks, check for conflicts at commit time, and abort if interference detected—faster when conflicts are rare but wasteful under contention.
* Today's focus: pessimistic two-phase locking, the most common database approach; optimistic schemes (e.g., FaRM) covered in later lectures.

### Two-Phase Locking (2PL)
* **Rule 1:** Acquire a lock before reading or writing any record.
* **Rule 2:** Hold all locks until transaction commits or aborts—no early release.
* Locks associated per-record (or coarser granularity); transactions race for locks, forcing a serial execution order dynamically.
* Early lock release risks non-serializable interleavings (e.g., T2 reading X after T1 modifies it but before T1 updates Y) or exposing phantom values if T1 aborts after T2 reads its changes.
* Deadlocks are common (T1 locks X then wants Y; T2 locks Y then wants X); systems detect cycles or use timeouts, aborting one transaction to break the cycle.

### Distributed Transactions & Atomic Commit
* Sharded data means a single transaction may modify records on multiple servers (e.g., X on server 1, Y on server 2).
* Partial failures (server 2 crashes or network partitions) can leave only half the transaction applied, violating atomicity.
* Need for **atomic commit protocol**: either all participants execute their part or none do, despite crashes, aborts, or missing records.

### Two-Phase Commit (2PC) Protocol
* **Participants:** Transaction coordinator (TC) executes transaction logic; participants (servers) hold sharded data and apply operations.
* **Phase 1 (Prepare):** TC sends "prepare" messages to all participants after completing transaction logic; each participant checks if it can commit (e.g., no deadlock, no missing record) and replies "yes" (vote to commit) or "no" (must abort).
* **Phase 2 (Commit/Abort):** If all votes are "yes," TC sends "commit" to all; if any vote is "no," TC sends "abort" to all; participants unlock data only after receiving commit or abort.
* Participants reply with acknowledgments (ACKs) after processing commit/abort, allowing TC to finalize.
* Correctness: unanimous "yes" votes ensure all participants can commit; any "no" triggers global abort, preserving all-or-nothing atomicity even when individual servers face local issues.

### Key Takeaways
* Two-phase locking enforces serializability by dynamically ordering transactions through lock conflicts, ensuring isolation despite concurrency.
* Two-phase commit achieves distributed atomicity by coordinating a voting protocol: all participants must agree to commit before any applies changes durably.
* These primitives layer to provide full ACID semantics in sharded systems, though performance costs (lock contention, prepare-phase latency) and availability concerns (blocking on coordinator) motivate alternative designs like optimistic concurrency control and consensus-based commits in modern systems.
