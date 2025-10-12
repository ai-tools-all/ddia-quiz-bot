### Goals & Context
* Build a globally distributed database supporting serializable transactions over data spread across multiple data centers worldwide while providing strong consistency guarantees.
* Google's advertising database drove initial design—data was manually sharded across many MySQL and Bigtable instances, requiring transparent automatic sharding and transactions spanning multiple shards.
* Workload dominated by read-only transactions (billions vs millions of read-write), making read performance optimization critical; external consistency required so committed transactions are immediately visible to subsequent transactions.

### Architecture & Data Distribution
* Data partitioned into shards, each replicated across multiple data centers using Paxos groups (typically 3-5 replicas per shard); every shard has a Paxos leader handling writes and coordinating replication.
* Servers spread globally across data centers; clients can read from nearby replicas for low latency while Paxos provides fault tolerance through majority-based replication.
* Replication enables geo-distribution for performance (data near users) and availability (surviving data center failures); Paxos majority requirement means one slow/failed data center won't block progress.

### Two-Phase Commit Over Paxos
* Transactions spanning multiple shards use two-phase commit for atomicity, but run it over Paxos-replicated participants rather than single servers to avoid blocking on coordinator crashes.
* Transaction coordinator (often co-located with one shard's Paxos leader) sends prepare messages; each shard's Paxos leader logs the prepare through its Paxos group before voting yes.
* After collecting unanimous yes votes, coordinator sends commit messages which are also logged through Paxos, ensuring crash recovery can complete transactions without blocking.

### Read-Only Transaction Challenge
* Paxos replication means minority replicas may lag behind and lack latest committed updates; naive local reads from nearest replica could return stale data violating consistency.
* Cannot simply read latest values across shards—might read from inconsistent time points (seeing writes from T1 on one shard, T2 on another when T2 happened first).
* Need mechanism ensuring reads see consistent snapshot of database as of specific point in time while still allowing fast local reads without contacting remote data centers.

### Snapshot Isolation with Timestamps
* Every transaction assigned timestamp; read-only transactions read data as of their timestamp, seeing most recent committed write with timestamp ≤ transaction's timestamp.
* Provides serializable snapshot isolation—transactions execute as if database frozen at specific timestamp, guaranteeing consistent view across all shards.
* Replicas track how up-to-date they are via timestamps; when read arrives requesting data at timestamp T, replica delays response until it has received all updates through timestamp T from Paxos leader.

### TrueTime API and Clock Uncertainty
* Physical clock synchronization across data centers impossible—network delays, clock drift, and relativity prevent perfect sync; best achievable is bounded uncertainty (Google targets ~7ms).
* TrueTime returns interval [earliest, latest] representing uncertainty bounds; actual time guaranteed within this interval using GPS receivers and atomic clocks at each data center.
* System architecture explicitly accounts for uncertainty rather than assuming synchronized clocks, using intervals to make correctness guarantees despite imperfect time.

### External Consistency via Start and Commit-Wait Rules
* Start rule: read-only transactions use latest bound of TrueTime interval as timestamp when starting; read-write transactions use latest bound when beginning commit.
* Commit-wait rule: after choosing timestamp, read-write transaction must delay committing until its timestamp < TrueTime.now().earliest, ensuring timestamp definitely in past before releasing locks.
* Together these rules guarantee external consistency—if T1 commits before T2 starts (in real time), T2's timestamp > T1's timestamp, so T2 sees T1's writes; commit-wait typically adds ~7-10ms latency.

### Key Takeaways
* Spanner demonstrates wide-area distributed transactions are practical by combining two-phase commit over Paxos for availability with timestamp-based snapshot isolation for efficient read-only transactions.
* TrueTime's explicit uncertainty bounds let system reason about time correctness despite imperfect clock synchronization, using commit-wait delays to bridge uncertainty gaps.
* Read-only transactions execute without locks or cross-data-center coordination by reading from local replicas at assigned timestamps, providing low latency for the common case while maintaining strict consistency.
