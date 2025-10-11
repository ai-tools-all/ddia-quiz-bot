# Task: Create Raft Subjective Questions

## Date: 2025-10-10
## Category: Task
## Status: Completed

## Overview
Created subjective interview questions about Raft consensus algorithm based on MIT 6.824 lecture transcripts and the provided guidelines.

## Input Resources
1. **Guidelines**: `/ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective/GUIDELINES.md`
   - Evaluation rubric with scoring criteria (excellent, good, needs improvement, insufficient)
   - Follow-up patterns for different answer levels
   - Timing guidelines and scoring weights

2. **Transcripts**: `/transcripts/mit-6824-subtitles/`
   - Lecture 5: Go Threads and Raft
   - Lecture 6: Fault Tolerance - Raft 1
   - Lecture 7: Fault Tolerance - Raft 2
   - Additional lectures covering distributed systems concepts

## Questions Created

### L3 Baseline
- **File**: `L3-baseline/03-raft-intro.md`
- **Topic**: Introduction to Raft Consensus
- **Focus**: Basic understanding of consensus problem and why Raft is needed

### L4 Baseline  
- **File**: `L4-baseline/03-raft-basics.md`
- **Topic**: Understanding Raft's Core Components
- **Focus**: Three server states (follower, candidate, leader) and transitions

### L5 Baseline
- **File**: `L5-baseline/03-raft-consensus.md`
- **Topic**: Raft Leader Election and Split Brain Prevention
- **Focus**: Majority voting, term numbers, network partition handling

### L5 Bar Raiser
- **File**: `L5-bar-raiser/03-raft-log-replication.md`
- **Topic**: Raft Log Replication and Consistency Guarantees
- **Focus**: Log matching property, AppendEntries protocol, commitment rules

### L6 Baseline
- **File**: `L6-baseline/02-raft-performance.md`
- **Topic**: Raft Performance Optimization and Trade-offs
- **Focus**: Bottlenecks, production optimizations, safety vs performance

### L7 Bar Raiser
- **File**: `L7-bar-raiser/02-consensus-evolution.md`
- **Topic**: Evolution from Classical Consensus to Modern Production Systems
- **Focus**: Paxos to Raft to EPaxos evolution, workload-driven protocol selection

## Key Concepts Covered

### Core Raft Concepts
- Leader election and prevention of split brain
- Log replication and consistency
- Term numbers and voting mechanisms
- Three server states and transitions
- Majority requirements and quorum

### Advanced Topics
- Performance optimizations (batching, pipelining, parallel apply)
- Comparison with other consensus protocols (Paxos, EPaxos, Flexible Paxos)
- Production system implementations
- Workload-specific protocol selection
- Geo-distributed optimizations

## Question Structure
Each question follows the standardized format:
- Main question with clear scope
- Core concepts (60% weight) that must be mentioned
- Peripheral concepts (40% weight) for differentiation
- Sample excellent and acceptable answers
- Common mistakes to watch for
- Three types of follow-up questions (excellent/partial/weak)
- Bar raiser challenge for level progression

## Implementation Notes
1. Questions progress from basic understanding (L3) to system design (L7)
2. Each level has appropriate complexity and time expectations
3. Follow-up questions test depth and recovery abilities
4. Bar raiser questions bridge to next level expectations
5. Real-world examples from production systems included

## Quality Checks
- ✅ Aligned with evaluation rubric from GUIDELINES.md
- ✅ Based on actual Raft implementation details from MIT lectures
- ✅ Progressive difficulty across levels
- ✅ Mix of theoretical and practical aspects
- ✅ Industry-relevant examples and optimizations
