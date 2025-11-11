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
- [x] Generate MCQ questions
- [x] Generate subjective questions
- [ ] Commit and push

## Generated Content

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

### Quiz Structure Plan
- MCQ: 6 questions covering core concepts
- Subjective L3: 2 questions (6-8 min each) - basics
- Subjective L4: 2 questions (8-10 min each) - intermediate
- Subjective L5: 2 questions (10-12 min each) - advanced design
- Subjective L6: 2 questions (12-15 min each) - system design and trade-offs
