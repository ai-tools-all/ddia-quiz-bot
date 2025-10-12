### Goals & Context
* Build a successor to MapReduce that supports iterative algorithms and multi-step data flow, targeting data center computations that loop over massive datasets repeatedly.
* Enable expressive programming through flexible transformation chains while giving the system visibility into the entire computation for optimization and fault recovery.

### Programming Model & Lineage Graphs
* Applications construct a lineage graph describing transformations rather than immediately executing operations; the driver program builds this recipe using functional operators like map, filter, join, and reduceByKey.
* Transformations are "lazy"—execution defers until an action (collect, count, save) forces materialization; this delayed evaluation lets Spark analyze and optimize the entire workflow before running anything.
* Each transformation produces a Resilient Distributed Dataset (RDD) representing a partitioned collection; RDDs track their lineage back to base inputs, enabling recomputation on failure without checkpointing every intermediate result.

### Narrow vs. Wide Dependencies
* Narrow dependencies occur when each output partition depends on at most one input partition (e.g., map, filter); workers process these locally without communication, enabling efficient pipelined execution.
* Wide dependencies require shuffling data across workers because output partitions depend on multiple input partitions (e.g., groupByKey, reduceByKey, join); these operations trigger expensive network transfers and sorting.
* Spark schedules narrow transformations in stages that run within worker memory, batching I/O and avoiding intermediate file writes; wide dependencies force stage boundaries where data must be exchanged and materialized.

### Persistence & Caching
* By default, Spark discards intermediate results after passing them to the next transformation, conserving memory; iterative algorithms that reuse datasets (like PageRank's link graph) must explicitly cache or persist data.
* The `cache()` call tells Spark to retain an RDD in worker memory across iterations; persist offers finer control, allowing storage in memory, disk, or replicated HDFS depending on size and fault-tolerance needs.
* Cached data eliminates redundant recomputation in loops—without caching, each iteration would reread the base input and rerun all preparatory transformations; with caching, the loop only recomputes updated ranks.

### Fault Tolerance via Recomputation
* Spark recovers from worker failures by replaying lost partitions using lineage; because RDDs remember their parent transformations, a crashed worker's tasks can rerun on survivors by reprocessing the same input partitions.
* Narrow dependencies localize recovery—only the failed partition's lineage needs recomputation; wide dependencies complicate recovery because rebuilding one partition may require re-executing transformations on all workers if intermediate outputs were discarded.
* Checkpointing mitigates recovery cost for long lineages by periodically writing RDD snapshots to replicated HDFS; after a checkpoint, recovery reads from stable storage instead of recomputing from the original input, trading storage overhead for shorter recovery chains.

### Execution & Optimization
* The driver compiles the lineage graph into Java bytecode stages and ships them to workers; workers execute transformations on their local partitions, reading from HDFS or cached memory and writing shuffle outputs for downstream stages.
* Spark optimizes by co-locating partitions that share keys (e.g., ensuring grouped data and ranks hash to the same workers for join), avoiding unnecessary shuffles when consecutive wide operations use compatible partitioning.
* PageRank demonstrates these patterns: map/filter steps pipeline within stages, groupByKey and join trigger shuffles, caching the link structure avoids rereading the web graph each iteration, and periodic checkpoints cap recovery time for hundred-iteration runs.

### Key Takeaways
* Spark generalizes MapReduce's rigid two-phase model into a flexible DAG of transformations, enabling iterative and multi-stage workflows with far less boilerplate and intermediate I/O.
* Lineage-based fault tolerance replaces eager replication with lazy recomputation, keeping intermediate data in memory and rerunning only lost partitions when failures occur.
* Explicit caching and strategic checkpointing balance memory usage, performance, and recovery cost, letting programmers tune persistence for workload characteristics while the system handles distribution transparently.
