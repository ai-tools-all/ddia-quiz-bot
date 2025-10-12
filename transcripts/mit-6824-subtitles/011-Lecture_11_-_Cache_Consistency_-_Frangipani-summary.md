### Goals & Context
* Build a shared network file system for a trusted research lab that feels like local Unix storage while supporting collaboration across dozens of workstations.
* Deliver strong consistency, fast interactive performance, and transparent crash recovery despite aggressive client-side write-back caching.

### Architecture & Workload Assumptions
* Every workstation runs a full Frangipani kernel module that implements filesystem logic and caches metadata/data in local memory; Petal supplies the shared virtual disk backing store with replicated blocks across server pairs.
* A separate lock service tracks which workstation holds each inode-level lock; Frangipani assumes cooperative users and omits strong security isolation.
* Typical workloads are read-heavy on a user’s own files with bursts of collaborative edits, motivating large caches and optimistic local execution.

### Cache Coherence via Leased Locks
* Clients must acquire a lock before caching any block; cached state is tied to the lock entry, which records whether the workstation is actively mutating it.
* The coherence protocol revolves around four RPCs: `request` (client → lock server), `grant`, `revoke`, and `release`; revocations force the current holder to flush dirty blocks to Petal before handing off.
* Locks support shared-read vs exclusive-write modes and an “idle” state so workstations can retain recently used locks without chatter; leases cap how long an unresponsive client may hold a lock.

### Atomicity Through Transactions
* File operations (create, rename, delete) lock all affected inodes/directories up front, perform modifications against cached copies, then write them to Petal before releasing locks.
* Holding the entire lock set until the operation is durable ensures other clients either see the previous state or the fully updated state—never a torn directory entry.
* The same lock table underpins both coherence (propagating visibility) and isolation (masking intermediate updates), yielding transactional behavior with little extra machinery.

### Write-Back Caching, Logging, and Durability
* Workstations buffer writes locally and lazily flush dirty blocks (at least every ~30 seconds) to match standalone Unix semantics while keeping interactive latency low.
* Before any dirty block protected by a revocation is emitted, the workstation appends a write-ahead log entry describing every metadata update in the operation; only after the log is stable in Petal are the corresponding blocks pushed.
* Each workstation owns a circular log region in Petal, indexed by monotonically increasing log sequence numbers and checksums; entries enumerate block IDs, version numbers, and new contents for metadata structures (directories, inodes, allocation maps).

### Crash Recovery Workflow
* Lock leases detect an unresponsive client; the lock server nominates a live workstation to replay the crashed node’s log directly from Petal before releasing its locks.
* Recovery scans the per-client log until sequence numbers stop increasing, redos each fully written entry, and ignores any torn records, ensuring partially written operations remain invisible.
* Because logs sit in shared storage, recovery can complete even if the crashed machine is offline or its local disk is lost; repeated replays are idempotent thanks to per-block version checks.

### Operational Considerations
* Heavy sharing can induce lock revocation storms; the paper notes batching revocations and using shared locks for read-mostly files to keep throughput reasonable.
* Log cleaning must keep ahead of new traffic so circular regions do not wrap while entries are still needed for recovery; background cleaners discard entries once all referenced blocks are safely checkpointed.
* Petal’s replication handles disk failures, while Frangipani’s design focuses on workstation crashes; security and untrusted client behavior remain out of scope.

### Key Takeaways
* Frangipani demonstrates that strong cache coherence, transactional metadata updates, and crash recovery can be layered on a simple shared-disk substrate by centralizing coordination in a lock service.
* Leased locks double as coherence tokens and transaction guards, providing fresh reads and atomic mutation without constant server mediation.
* Per-client write-ahead logs stored in shared storage let any survivor repair partial updates and free locks, giving the cluster forward progress even when a workstation dies mid-operation.
