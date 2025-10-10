# Subjective Question Format Proposal for GFS

## Proposed Markdown Structure

```markdown
---
id: gfs-subjective-[level]-[number]
type: subjective
level: L3|L4|L5|L6|L7
category: baseline|bar-raiser
topic: gfs
subtopic: [specific area like replication, consistency, etc.]
estimated_time: 5-10 minutes
---

# [Question Title]

## question
[Open-ended question that requires explanation]

## core_concepts
### Must Mention (Mandatory - 60% score)
- **Concept 1**: [Brief description of why this is essential]
- **Concept 2**: [Brief description]
- **Concept 3**: [Brief description]

### Keywords Expected
- Primary keywords: [master, chunk, replica, etc.]
- Technical terms: [lease, version number, etc.]

## peripheral_concepts
### Nice to Have (Bonus - 40% score)
- **Advanced Concept 1**: [Shows deeper understanding]
- **Trade-off Analysis**: [Demonstrates critical thinking]
- **Real-world Application**: [Shows practical knowledge]

### Bonus Keywords
- Implementation details: [pipelining, checkpointing, etc.]
- Related systems: [HDFS, Colossus, etc.]
- Performance metrics: [throughput, latency, etc.]

## evaluation_rubric
### Excellent (90-100%)
- Mentions ALL core concepts with clear understanding
- Includes 2+ peripheral concepts
- Makes connections between concepts
- Provides examples or analogies

### Good (70-89%)
- Mentions MOST core concepts (at least 2/3)
- Includes 1 peripheral concept
- Shows solid understanding

### Needs Improvement (50-69%)
- Mentions SOME core concepts (at least 1/3)
- May have minor inaccuracies
- Basic understanding evident

### Insufficient (0-49%)
- Misses most core concepts
- Shows fundamental misunderstanding
- Too vague or incorrect

## sample_answers
### Excellent Answer
"[Full sample answer demonstrating all core concepts and some peripheral ones]"

### Acceptable Answer
"[Sample showing minimum acceptable understanding]"

### Common Mistakes
- Mistake 1: [What candidates often get wrong]
- Mistake 2: [Another common error]

## follow_ups
### If Excellent Answer
**Depth Probe**: "[Question that goes deeper into the topic]"
- Looking for: [Advanced concept exploration]
- Red flags: [Over-engineering, missing subtleties]

### If Partial Answer (missed some core concepts)
**Guided Probe**: "[Hint-based question to steer toward missed concept]"
- Hint embedded: [The hint you're giving]
- Concept testing: [What core concept you're checking]

### If Weak Answer (missed most concepts)
**Foundation Check**: "[Step back to basics question]"
- Simplification: [Simpler version of the concept]
- Building block: [Fundamental understanding check]

## bar_raiser
### Next Level Question (L3→L4)
"[More complex scenario requiring next-level thinking]"
- Additional concepts needed: [What L4 would know]
- Complexity increase: [What makes this harder]

## interviewer_notes
- Time expectation: [How long should candidate take]
- Watch for: [Specific things to observe]
- Common paths: [How discussions typically evolve]
```

## Example Implementation for L3 GFS Question

```markdown
---
id: gfs-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: gfs
subtopic: replication
estimated_time: 5-7 minutes
---

# GFS Replication Strategy

## question
"Explain how GFS ensures data durability when a chunk server fails. Walk me through what happens from the moment of failure."

## core_concepts
### Must Mention (Mandatory - 60% score)
- **3x Replication**: GFS maintains 3 copies of each chunk by default
- **Master Detection**: Master detects failure through heartbeats
- **Re-replication**: Master initiates copying to restore replication factor

### Keywords Expected
- Primary keywords: replica, chunk server, master, failure
- Technical terms: heartbeat, replication factor

## peripheral_concepts  
### Nice to Have (Bonus - 40% score)
- **Rack Awareness**: Replicas spread across different racks
- **Priority Queue**: More under-replicated chunks get priority
- **Version Numbers**: Ensures stale replicas aren't used
- **Bandwidth Throttling**: Re-replication doesn't overwhelm network

### Bonus Keywords
- Implementation: pipelining, chunk priorities
- Related: Colossus improvements, HDFS similarities
- Metrics: MTTR, availability targets

## evaluation_rubric
### Excellent (90-100%)
- Explains 3-replica strategy clearly
- Describes heartbeat detection mechanism
- Details re-replication process
- Mentions rack awareness or version numbers
- May discuss trade-offs or limitations

### Good (70-89%)
- Mentions 3 replicas
- Understands master coordinates recovery
- Describes basic re-replication
- Minor gaps in mechanism details

### Needs Improvement (50-69%)
- Knows about replication but vague on count
- Understands failure detection happens
- Missing key details about recovery

### Insufficient (0-49%)
- Confuses replication with backup
- No clear understanding of master's role
- Cannot explain recovery process

## sample_answers
### Excellent Answer
"When a chunk server fails in GFS, several things happen: First, the master detects the failure when heartbeats stop arriving (typically after 60 seconds). The master then scans its metadata to identify all chunks that had replicas on the failed server. For each under-replicated chunk (now at 2 replicas instead of 3), the master adds it to a re-replication queue, prioritizing chunks with fewer replicas. The master then instructs healthy chunk servers holding replicas to copy their chunks to other available servers, spreading them across different racks when possible. This ensures the 3x replication factor is restored. The system uses version numbers to prevent stale replicas from being used if the 'failed' server comes back online."

### Acceptable Answer  
"GFS keeps 3 copies of each chunk on different servers. When a chunk server fails, the master notices through missing heartbeats and starts making new copies of the affected chunks on other servers to get back to 3 replicas."

### Common Mistakes
- Saying "backup" instead of "replica" (backups are different!)
- Forgetting the master's coordination role
- Thinking clients handle re-replication

## follow_ups
### If Excellent Answer
**Depth Probe**: "What happens if the master fails during re-replication? How does GFS handle that?"
- Looking for: Shadow master, operation log, checkpointing
- Red flags: Saying there's automatic master failover (there isn't!)

### If Partial Answer (missed detection mechanism)
**Guided Probe**: "You mentioned the master knows about failures. How would the master know that a chunk server that was working fine a minute ago has suddenly crashed?"
- Hint embedded: Time-based checking mechanism
- Concept testing: Heartbeat understanding

### If Weak Answer (vague on replication)
**Foundation Check**: "Let's start simpler - if you have important data, how many copies would you keep and where would you put them?"
- Simplification: Basic redundancy concept
- Building block: Understanding why multiple copies matter

## bar_raiser
### Next Level Question (L3→L4)
"Now consider this scenario: A rack switch fails, taking down 20 chunk servers at once. Each server held 1000 chunks. How should GFS prioritize the re-replication? What problems might arise?"
- Additional concepts: Prioritization algorithms, network bandwidth limits, cascading failures
- Complexity increase: Scale, resource constraints, system-wide thinking

## interviewer_notes
- Time expectation: 2-3 minutes initial answer, 3-4 minutes discussion
- Watch for: Confusion between replication and sharding, understanding of eventual consistency
- Common paths: Usually leads to discussion of CAP theorem or comparison with HDFS
```

## Format Benefits

1. **Clear Scoring**: Objective criteria for evaluation
2. **Adaptive**: Follow-ups adjust to candidate's level
3. **Teaching Moment**: Even weak answers lead to learning
4. **Bar Raiser**: Tests readiness for next level
5. **Interviewer Friendly**: Notes help maintain consistency

## Proposed File Structure
```
09-distributed-systems-gfs/
├── subjective/
│   ├── L3-baseline/
│   │   ├── 01-replication-basics.md
│   │   ├── 02-consistency-understanding.md
│   │   └── 03-failure-handling.md
│   ├── L3-bar-raiser/
│   │   └── 01-scale-challenges.md
│   ├── L4-baseline/
│   │   ├── 01-trade-offs.md
│   │   └── 02-design-decisions.md
│   └── L4-bar-raiser/
│       └── 01-system-evolution.md
└── multiple-choice/
    └── [existing quiz questions]
```

## Questions Per Level Breakdown
- **L3**: 3 baseline + 1 bar raiser
- **L4**: 3 baseline + 1 bar raiser  
- **L5**: 3 baseline + 1 bar raiser
- **L6**: 2 baseline + 1 bar raiser
- **L7**: 2 baseline + 1 bar raiser

Total: 20 subjective questions with full follow-up trees

Does this format work for your needs? Should I adjust anything before implementing the GFS subjective questions?
