---
id: primary-backup-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: primary-backup
subtopic: replication-approaches
estimated_time: 5-7 minutes
---

# question_title - State Transfer vs Replicated State Machine

## main_question - Core Question
"Explain the difference between state transfer and replicated state machine approaches to replication. When would you choose one over the other?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **State Transfer**: Primary sends entire state/snapshots to backup
- **Replicated State Machine**: Primary sends operations/inputs to backup
- **Determinism Requirement**: RSM requires deterministic execution

### expected_keywords
- Primary keywords: state, operations, transfer, deterministic
- Technical terms: snapshot, checkpoint, input stream, execution

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Bandwidth Considerations**: State transfer can be expensive for large state
- **Recovery Speed**: State transfer might be faster for initial sync
- **Complexity**: RSM more complex but potentially more efficient
- **VMware FT Example**: Uses replicated state machine approach

### bonus_keywords
- Trade-offs: network bandwidth, storage, CPU usage
- Implementation: logging, replay, synchronization
- Examples: database replication, VM replication

## sample_excellent - Example Excellence
"State transfer and replicated state machine (RSM) are two fundamental approaches to replication. In state transfer, the primary periodically sends its entire state or state changes to the backup - like taking snapshots of memory and disk. The backup directly receives and applies this state. In RSM, the primary sends the sequence of operations or inputs it receives to the backup, and the backup executes these same operations independently. RSM requires deterministic execution - given the same inputs, both replicas must produce identical state. State transfer is simpler and good when state is small or changes infrequently, like configuration servers. RSM is more bandwidth-efficient for systems with small operations but large state, like databases or VMs. VMware FT uses RSM by sending input events rather than memory contents."

## sample_acceptable - Minimum Acceptable
"State transfer means the primary sends its actual state (memory, disk contents) to the backup. Replicated state machine means the primary sends the operations it receives to the backup, which then executes them. RSM needs deterministic execution so both get the same result. State transfer is simpler but uses more network bandwidth for large states."

## common_mistakes - Watch Out For
- Confusing state changes with full state transfer
- Not mentioning determinism requirement for RSM
- Missing bandwidth/efficiency trade-offs
- Thinking RSM always better than state transfer

## follow_up_excellent - Depth Probe
**Question**: "In VMware FT, why did they choose replicated state machine over periodic state transfer? What challenges does this create?"
- **Looking for**: Non-deterministic operations, performance impact, output commit
- **Red flags**: Not understanding determinism challenges

## follow_up_partial - Guided Probe  
**Question**: "If you're replicating a 1TB database, which approach would likely use less network bandwidth for normal operations?"
- **Hint embedded**: Consider operation size vs state size
- **Concept testing**: Understanding efficiency trade-offs

## follow_up_weak - Foundation Check
**Question**: "Imagine backing up your laptop. Would you send the entire disk contents every time, or just the list of files you changed?"
- **Simplification**: Incremental vs full backup analogy
- **Building block**: Efficiency of sending operations vs state

## bar_raiser_question - L3→L4 Challenge
"You're designing replication for a chess game server. Each game has small state (board position) but many moves. Players expect instant response. Should you use state transfer after each move, or replicated state machine? Consider failure recovery time."

### bar_raiser_concepts
- State size vs operation frequency
- Recovery time implications
- Network bandwidth optimization
- User experience during failover

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Deterministic replay, output commit, VMware FT architecture
