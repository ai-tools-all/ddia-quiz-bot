# Generate OCC Quiz - 2025-11-11

## Task
Generate MCQ and subjective quiz for MIT 6.824 Lecture 14: Optimistic Concurrency Control

## Files to Process
- transcripts/mit-6824-subtitles/014-Lecture_14_-_Optimistic_Concurrency_Control.en.srt
- transcripts/mit-6824-subtitles/014-Lecture_14_-_Optimistic_Concurrency_Control-summary.md

## Reference
- Example folder: ddia-quiz-bot/content/chapters/14-mit-6824-optimistic-cc
- Prompts: prompts/*.md

## Progress
- [x] Read reference prompts
- [x] Examine example folder structure
- [x] Read transcript and summary
- [x] Generate MCQ questions (Round 1: 6 questions)
- [x] Generate subjective questions (Round 1: 8 questions)
- [x] Generate additional questions (Round 2: 10 questions)
- [ ] Commit and push

## Generated Content - Round 1 (Initial 14 questions)

### MCQ Questions (6 total)
1. `01-why-occ-with-rdma.md` - Why FaRM uses OCC instead of locks
2. `02-version-lock-bits.md` - Purpose of lock bits in commit protocol
3. `03-validate-optimization.md` - VALIDATE optimization for read-only transactions
4. `04-two-phase-commit.md` - LOCK phase of 2PC protocol
5. `05-nvram-benefit.md` - Non-volatile RAM advantages
6. `06-contention-aborts.md` - OCC performance under contention

### Subjective Questions (8 total)

#### L3-baseline (2 questions, 6-8 min each)
1. `01-occ-workflow.md` - OCC workflow and conflict example
2. `02-rdma-one-sided-ops.md` - RDMA one-sided operations and design constraints

#### L4-baseline (2 questions, 8-10 min each)
1. `01-validate-vs-lock.md` - Comparing VALIDATE and LOCK operations
2. `02-2pc-coordinator-crash.md` - Coordinator crash scenarios in 2PC

#### L5-baseline (2 questions, 10-12 min each)
1. `01-version-lock-atomicity.md` - Object header design and atomic CAS operations
2. `02-phantom-reads-occ.md` - Phantom read problem in OCC

#### L6-baseline (2 questions, 12-15 min each)
1. `01-farm-vs-spanner-tradeoffs.md` - FaRM vs Spanner design trade-offs
2. `02-occ-contention-mitigation.md` - Strategies for mitigating OCC contention

## Generated Content - Round 2 (Additional 10 questions for practical understanding)

### MCQ Questions (3 additional)
7. `07-primary-backup-timing.md` - When COMMIT-BACKUP messages are sent
8. `08-object-id-structure.md` - Object ID structure enabling direct RDMA access
9. `09-wal-placement.md` - Write-ahead log placement in per-client queues

### Subjective Questions (7 additional)

#### L3-baseline (2 additional)
3. `03-commit-backup-durability.md` - COMMIT-BACKUP protocol and fault tolerance
4. `04-serializability-validation.md` - How version validation ensures serializability

#### L4-baseline (2 additional)
3. `03-wal-crash-recovery.md` - WAL-based crash recovery process
4. `04-retry-strategies.md` - Transaction retry with backoff and fairness

#### L5-baseline (2 additional)
3. `03-sharding-coordination.md` - 90-way sharding and multi-shard coordination
4. `04-read-write-optimization.md` - Optimizing mixed read-write transactions

#### L6-baseline (1 additional)
3. `03-partition-handling.md` - Network partition and split-brain prevention

## Notes

### Summary Analysis
Key topics from the lecture:
1. RDMA hardware and one-sided operations (kernel bypass, NIC)
2. Non-volatile RAM (battery-backed DRAM)
3. OCC motivation - why not locking with RDMA
4. Transaction API (TX_create, TX_read, TX_write, TX_commit)
5. Version numbers and lock bits in object headers
6. Two-phase commit protocol with OCC (LOCK, COMMIT-PRIMARY, COMMIT-BACKUP)
7. Validate optimization for read-only transactions
8. Fault tolerance (F+1 replicas, WAL in per-client queues)
9. Performance trade-offs (contention, single datacenter vs geo-distributed)

### Final Quiz Structure
- **MCQ: 9 questions** covering core concepts, replication, addressing, and logging
- **Subjective L3: 4 questions** (6-8 min each) - OCC basics, RDMA, replication, serializability
- **Subjective L4: 4 questions** (8-10 min each) - VALIDATE/LOCK, coordinator crashes, WAL recovery, retries
- **Subjective L5: 4 questions** (10-12 min each) - atomic operations, phantoms, sharding, optimization
- **Subjective L6: 3 questions** (12-15 min each) - system comparisons, contention mitigation, partitions
- **Total: 24 questions** comprehensively testing FaRM's OCC implementation
