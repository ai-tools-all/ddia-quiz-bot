### Goals & Context
* Build an extremely high-performance transactional database for a single data center that exploits modern RDMA networking hardware to achieve microsecond-level transaction latencies.
* Target CPU time as the primary bottleneck rather than network latency by co-locating replicas, achieving ~100x faster transactions than geographically distributed systems like Spanner (58 microseconds vs 10+ milliseconds).

### Architecture & Hardware Foundation
* All data sharded across primary-backup pairs within a single data center; configuration manager (using ZooKeeper) tracks which servers hold primary and backup roles for each shard.
* Non-volatile RAM stores all data in memory for fast access while surviving power failures through battery-backed DRAM; writes go directly to RAM without disk persistence.
* RDMA (Remote Direct Memory Access) NICs allow one-sided operations where clients directly read/write server memory without interrupting server CPUs, combined with kernel bypass for application-level network access.
* Each server memory contains object regions (with versioned objects having headers), plus per-client message queues and logs enabling N² communication channels across the cluster.

### Optimistic Concurrency Control (OCC) Motivation
* RDMA's one-sided operations fundamentally restrict design options—clients cannot check or acquire locks when directly reading memory, making traditional pessimistic locking infeasible.
* OCC allows transactions to read data optimistically without locking, buffer writes locally, then validate at commit time whether reads remained consistent; conflicts detected through version number mismatches trigger aborts.
* This design enables extremely fast one-sided RDMA reads during transaction execution, deferring all coordination and conflict detection to the commit phase.

### Transaction API & Execution Model
* Applications declare transactions with TX_create, read objects via TX_read (supplying object IDs), modify copies locally, and stage writes with TX_write before calling TX_commit.
* Object identifiers combine a region number (mapping to primary/backup servers) plus a memory address, enabling clients to route RDMA operations directly to the correct server locations.
* Transaction coordinators run on clients (often co-located with storage servers), orchestrating reads during execution and managing a variant of two-phase commit at commit time.

### Version Numbers and Lock Bits
* Every object header contains a version number incremented on each modification, plus a high-order lock bit set during commit processing to prevent concurrent conflicting transactions.
* Clients record the version number when first reading an object; at commit time, lock messages sent to primaries check that versions haven't changed and atomically set lock bits.
* Lock bit detection causes transactions to abort if another transaction is simultaneously committing changes to the same object, while version mismatches abort transactions that read stale data.

### Two-Phase Commit with OCC
* Commit begins by sending LOCK messages to all primaries holding objects the transaction will modify; primaries verify version numbers match expected values, set lock bits, and respond yes/no.
* If any primary responds no (due to version mismatch or lock bit already set), the coordinator aborts; otherwise it sends COMMIT-PRIMARY messages telling primaries to apply writes, increment versions, and clear locks.
* COMMIT-BACKUP messages propagate updates to backups for fault tolerance; the entire protocol ensures serializability by detecting and aborting conflicting concurrent transactions through version and lock checks.

### Validate Optimization for Read-Only Objects
* Transactions can use faster VALIDATE operations for objects they read but didn't modify, replacing lock messages with one-sided RDMA reads of object headers to check version numbers and lock bits.
* Read-only transactions execute entirely through RDMA reads followed by validation, avoiding expensive log appends and primary processing, dramatically improving throughput for read-heavy workloads.
* Validation still ensures serializability—transactions abort if lock bits are set or versions changed between initial read and validation, preventing them from committing computations based on inconsistent snapshots.

### Fault Tolerance & Recovery
* F+1 replicas per shard tolerate F failures; only one surviving replica needed (not a majority), and data center-wide power failures recoverable once machines restart since RAM contents persist.
* Write-ahead logs in per-client queues stored in server memory enable crash recovery, though details focus more on replication protocol than full recovery mechanisms.

### Performance & Trade-offs
* Massive sharding (90-way in experiments) provides automatic parallelism; all-data-in-RAM limits scale but enables sub-microsecond data access; RDMA hardware delivers unprecedented network performance.
* OCC performs well under low contention but aborts heavily conflicting transactions, requiring application-level retries (potentially with exponential backoff though not explicitly mentioned).
* Single-datacenter scope sacrifices geographic fault tolerance and disaster recovery compared to Spanner, trading global availability for extreme local performance.

### Key Takeaways
* FaRM demonstrates that RDMA hardware enables order-of-magnitude transaction speedups by eliminating server CPU involvement during reads, but requires adopting optimistic concurrency control since one-sided operations cannot enforce locks.
* Version numbers and lock bits provide lightweight serializability checking—transactions validate that data remained stable between read and commit, aborting on conflicts detected through version mismatches or concurrent lock attempts.
* The system illustrates a different point in the transaction design space than Spanner: prioritizing single-datacenter microsecond latencies over wide-area geographic replication and disaster tolerance.
