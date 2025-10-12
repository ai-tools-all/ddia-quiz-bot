### Goals & Context
* Scale Facebook's infrastructure to handle billions of user requests while keeping response times low, accepting eventual consistency for most operations while ensuring users always see their own writes.
* Leverage memcached as a look-aside cache layer between stateless web servers and MySQL database servers to achieve read throughput that databases alone cannot deliver economically.
* The paper demonstrates how aggressive caching creates consistency challenges that require careful protocol design to balance performance with acceptable staleness bounds.

### Architectural Evolution & Look-Aside Caching
* Facebook evolved from single-server setups through sharded databases to a memcached layer absorbing 99%+ of reads, leaving databases handling only writes and cache misses.
* Look-aside caching places responsibility on front-end PHP code: on reads, check memcached first; on misses, fetch from MySQL and populate cache; clients manage the relationship between cached and database state.
* This differs from look-through caching where the cache itself fetches from backing stores; the decoupling allows front-ends to cache arbitrary transformations of database records rather than raw rows.

### Write Invalidation Protocol
* Front-ends send writes to MySQL, then immediately issue deletes to memcached for affected keys, forcing subsequent readers to fetch fresh data from the database and repopulate the cache.
* Invalidation (delete) is preferred over updates (set) because concurrent writes can leave stale values in memcached if two clients' sets arrive out of order relative to their database commits.
* MySQL servers also asynchronously send deletes by monitoring the replication log, ensuring invalidations propagate even if front-end deletes are lost; front-end deletes primarily ensure clients see their own writes immediately.

### Regional Replication Architecture
* Facebook deployed multiple regions (West Coast primary, East Coast secondary) each holding complete replicas of all data to reduce cross-country latency for user-facing reads.
* All writes flow to the primary region's MySQL master; asynchronous replication propagates updates to secondary regions with seconds of lag, accepting stale reads for users far from the primary.
* Reads remain local within each region (both memcached and MySQL queries), exploiting the read-heavy workload while tolerating brief inconsistency across coasts; serving a single page often requires hundreds of data items, making local access critical.

### Intra-Region Clustering Strategy
* Within each region, front-ends and memcached servers group into independent clusters to limit N-squared TCP connection overhead and reduce incast congestion (hundreds of responses arriving simultaneously).
* Data partitions across memcached servers within a cluster, but popular keys benefit from replication across clusters, allowing parallel serving of hot items that would bottleneck a single sharded server.
* A regional pool of memcached servers shared across clusters caches infrequently accessed items, avoiding wasteful replication of cold data while dedicating per-cluster RAM to hot keys.

### Partitioning vs. Replication Trade-offs
* Databases use sharding (partitioning) to split write load and data volume; memcached employs both sharding within clusters for capacity and replication across clusters for handling popular keys.
* Sharding provides independent parallelism and scales total capacity but cannot alleviate hot spots; replication multiplies serving capacity for popular data at the cost of RAM and consistency complexity.
* The tension between performance (which replication aids by distributing load) and consistency (which replication undermines by creating multiple copies to synchronize) pervades the design.

### Consistency Model & User Expectations
* Facebook accepts seconds of staleness for casual browsing (news feeds, profiles) because users rarely notice slightly outdated content, prioritizing low-latency reads over linearizability.
* The critical consistency requirement is read-your-writes: after updating data, that user must immediately see the new value; front-end deletes after database writes guarantee this by forcing cache misses.
* Long-term stale data cached indefinitely would harm user experience, so invalidation protocols and timeouts ensure caches refresh within reasonable windows even if some deletes fail.

### Key Takeaways
* Memcached's look-aside pattern and invalidation-on-write scheme enable Facebook to absorb massive read load on commodity hardware while databases handle durable writes and provide ground truth.
* Regional replication with local reads trades global consistency for user-perceived latency, accepting that East Coast users see slightly stale data replicated from the West Coast primary.
* Clustering within regions and combining sharding with replication balances hot-key performance, memory efficiency, and network overhead, illustrating practical engineering tradeoffs at extreme scale.
