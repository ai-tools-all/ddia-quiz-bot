### Fault Tolerance and High Availability
* **Goal:** Keep a service running despite machine or network failures by provisioning redundant replicas.
* **Technique:** Replicate stateful servers so a backup can continue execution if the primary halts.

### Failure Model
* **Fail-Stop Focus:** Replication assumes components either run correctly or halt (e.g., power loss, NIC unplugged, thermal shutdown).
* **Out of Scope:** Bugs in the replicated software, latent CPU design flaws, or correlated disasters (shared power domains, earthquakes) can still wipe out all replicas.

### Replication Approaches
* **State Transfer:** Periodically copy full primary state to backup (simple but bandwidth heavy, during failover backup may start from slightly stale snapshot).
* **Replicated State Machine (RSM):** Keep replicas in sync by replaying the same input stream in identical order; relies on deterministic execution.

### VMware FT Architecture
* **VM-Level Replication:** Primary and secondary run identical VMs under VMware’s VMM; transparency allows unmodified applications/OSs.
* **Logging Channel:** Primary VMM captures all non-deterministic events (device interrupts, network packets, CPU results) as log entries and streams them to the backup.
* **Deterministic Replays:** Backup VMM enforces identical timing—interrupts injected at the same instruction number; non-deterministic instructions use primary-provided results.

### Handling Outputs and Failover
* **Output Rule:** Primary withholds external outputs until backup acknowledges receipt of corresponding log entries, guaranteeing the backup can resume without missing client-visible actions.
* **Failure Detection:** Heartbeats or lost TCP keep-alives reveal primary death; backup promotes its VM and begins emitting buffered outputs.
* **Performance Cost:** Output commit rule adds latency; batching and pipelined log delivery help but don’t remove the primary-to-backup round trip.

### Split-Brain Avoidance
* **Risk:** Network partition could leave both primary and backup believing the other is down, producing divergent outputs.
* **Solution:** Use an external test-and-set/arbitration service (often on shared storage) so only one replica can seize the “I am primary” token.

### Practical Concerns
* **Checkpointing State:** Occasional snapshots keep log size bounded and accelerate recovery.
* **I/O Diversity:** Need drivers that expose enough nondeterminism to the VMM; DMA arrivals, device IRQ ordering, and timer jitter must all be captured.
* **Resource Provisioning:** Primary and backup must have comparable CPU, memory, and NIC capacity to stay instruction-locked.
