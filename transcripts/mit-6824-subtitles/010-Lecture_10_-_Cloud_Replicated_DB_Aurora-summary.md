### Goals & Context
* Explain how Amazon Aurora delivers a managed, crash-tolerant SQL database with cloud-native performance.
* Understand why general-purpose EC2 disks and EBS volumes could not meet customer expectations for durability across availability zones (AZs) and for write-heavy workloads.

### From EC2/EBS to Aurora
* Early EC2 databases relied on locally attached disks; losing a VM meant losing storage unless the operator managed backups manually.
* EBS introduced replicated block storage, but replicas sat in one AZ and chain-replicated 8KB pages, saturating networks and limiting throughput.
* Customer demand for cross-AZ durability and higher write rates pushed Amazon to abandon generic blocks and build an application-aware storage service.

### Architecture Highlights
* Aurora keeps MySQL-compatible compute on EC2 while offloading durability to a fleet of storage nodes; compute sends redo/undo log records instead of whole pages.
* Each protection group spans three AZs with six storage nodes; the database instance is the sole writer issuing monotonically numbered log records.
* The control plane monitors the writer, promotes a standby on failure, and reattaches it to the same replicated volume within tens of seconds.
* The storage service is purpose-built: it understands log structure, applies records to pages lazily, and supports features such as snapshots, continuous backups, and fast copy-on-write clones.

### Quorum-Based Storage Layer
* Writes require acknowledgements from any four of the six replicas, trading synchronous quorum for tolerance of one lost AZ plus an extra node.
* Reads normally go to a single up-to-date replica because the writer tracks each node’s durable log prefix; quorum reads are reserved for recovery or page repair probes.
* Quorum math (R + W > N) gives linearizable commits while letting the system ignore the slowest replicas and survive transient slowness.
* Heartbeats and “repair” reads detect lagging replicas, which request any missing log segments before they are allowed back into the fast path.

### Page Reads, Crash Recovery, and Log Application
* Storage nodes persist redo/undo log fragments and older page versions; they materialize pages on demand when a compute node requests them.
* The database cache keeps hot pages locally, while background threads issue “log apply” requests so storage shards eventually converge to a fresh base image.
* On crash, a replacement database instance performs quorum reads to find the highest contiguous log index, truncates partial transactions, and replays committed updates.
* Undo records allow Aurora to roll back in-flight transactions without waiting for background page flushing, preserving MySQL semantics on top of the distributed log.

### Scaling Storage and Fast Repair
* Data is striped into 10GB protection groups; large databases simply add more groups, each with its own six-node quorum and subset of the log.
* A storage server hosts fragments for hundreds of customers; when it fails, Aurora rebuilds each fragment on different machines in parallel, avoiding terabyte-scale copies over a single link.
* Segment-aware replication plus background scrubbing keeps durability high even with correlated hardware failures or AZ outages.
* Continuous backups stream log state to S3, enabling point-in-time restore without quiescing the writer or copying entire volumes.

### Read Replicas and Mini-Transactions
* Aurora supports up to 15 read-only replicas that share the storage fleet; replicas receive the log stream to keep caches warm but do not participate in quorum writes.
* Because pages might be mid-rebalance, the storage tier enforces “mini-transactions” (VDL/VCL) so replicas see either the before or after image, never torn B-tree state.
* Replica lag is small (tens of milliseconds) since only committed log records are applied, giving near-real-time analytics without burdening the primary writer.
* Cross-region replicas consume the same shared log stream, allowing globally distributed readers on top of a single durable storage backend.

### Key Takeaways
* Separating compute from an application-specific storage service slashes network traffic and enables dramatic (>30×) throughput gains over block storage.
* Quorum replication across AZs provides durability targets customers expect while maintaining write latency through selective acknowledgement.
* Purpose-built cloud storage lets Aurora offer managed failover, fast repairs, clones, and snapshots that would be onerous for self-managed MySQL on EC2.
