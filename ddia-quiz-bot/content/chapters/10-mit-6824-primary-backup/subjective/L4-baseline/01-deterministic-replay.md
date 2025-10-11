---
id: primary-backup-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: primary-backup
subtopic: deterministic-replay
estimated_time: 6-8 minutes
---

# question_title - Deterministic Execution in Replicated State Machines

## main_question - Core Question
"In VMware FT's replicated state machine approach, what makes execution non-deterministic, and how does the system handle these challenges? Explain why deterministic replay is crucial for maintaining replica consistency."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Non-Deterministic Sources**: Time, randomness, I/O, interrupts, multi-core
- **Logging Approach**: Primary logs non-deterministic events with results
- **Replay Mechanism**: Backup replays using logged information
- **Consistency Requirement**: Both must reach identical state

### expected_keywords
- Primary keywords: deterministic, non-deterministic, logging, replay
- Technical terms: interrupts, random numbers, timing, instruction counter

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Instruction Counter**: Precise interrupt delivery timing
- **Output Commit Problem**: Delaying output until backup synchronized
- **Performance Impact**: Logging overhead and synchronization costs
- **Multi-core Challenges**: Why VMware FT initially single-core only

### bonus_keywords
- Implementation: log buffer, synchronization points, commit protocol
- Sources: RDTSC, RDRAND, network packets arrival, disk I/O completion
- Trade-offs: performance vs consistency, latency vs throughput

## sample_excellent - Example Excellence
"Deterministic replay is essential because replicated state machines require that given identical inputs, all replicas produce identical state. Non-deterministic events break this assumption. Sources include: timing of interrupts (network packets, disk I/O completion), random number generation, reading the timestamp counter, and in multi-core systems, memory access interleaving. VMware FT handles this by having the primary log all non-deterministic events with their values and the exact instruction count when they occurred. The backup doesn't generate these events independently but replays them from the log at the exact same instruction count. For example, when a network packet arrives at the primary at instruction 1000000, this is logged, and the backup delivers the same packet content at the same instruction count. This ensures both VMs remain in lockstep. The challenge is performance - all this logging and precise replay adds overhead, which is why VMware FT initially supported only single-core VMs where instruction ordering is deterministic."

## sample_acceptable - Minimum Acceptable
"Non-deterministic events include things like random numbers, current time readings, when interrupts happen, and I/O timing. These would cause primary and backup to diverge since they might happen differently on each. VMware FT solves this by having the primary log all these non-deterministic events with when they happened (instruction count), and the backup replays them at exactly the same point instead of experiencing them independently. This keeps both VMs in the exact same state."

## common_mistakes - Watch Out For
- Not explaining WHY non-determinism breaks RSM
- Missing critical sources like interrupt timing
- Confusing deterministic execution with identical hardware
- Not understanding instruction-level precision requirement

## follow_up_excellent - Depth Probe
**Question**: "VMware FT struggled with multi-core support. What additional non-determinism does multi-core introduce, and why is it harder to handle than single-core non-determinism?"
- **Looking for**: Memory access interleaving, cache coherence, race conditions
- **Red flags**: Not understanding concurrent execution challenges

## follow_up_partial - Guided Probe  
**Question**: "If the primary VM calls random() and gets 42, what exactly needs to be logged for the backup to stay synchronized?"
- **Hint embedded**: Value AND timing
- **Concept testing**: Understanding complete logging requirements

## follow_up_weak - Foundation Check
**Question**: "Two identical computers run the same program that prints the current time. Will they print the same value? What does this mean for replication?"
- **Simplification**: Basic non-determinism example
- **Building block**: Why logging is necessary

## bar_raiser_question - L4→L5 Challenge
"You're implementing replicated state machine for a database. A transaction reads current time, generates a UUID, and inserts a record. Network packet with commit arrives during UUID generation. Detail exactly what must be logged and when the backup processes each event to maintain consistency."

### bar_raiser_concepts
- Multiple non-deterministic events in sequence
- Interrupt timing precision
- Ordering of events
- State consistency verification
- Performance optimization opportunities

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Output commit, log optimization, multi-core challenges
