---
id: zookeeper-subjective-L4-bar-raiser-001
type: subjective
level: L4
category: bar-raiser
topic: zookeeper
subtopic: distributed-locks
estimated_time: 10-12 minutes
---

# question_title - Building a Fair Distributed Lock Service

## main_question - Core Question
"Design and implement a fair, reentrant distributed lock using Zookeeper that handles client failures gracefully and prevents starvation. Your lock should support both exclusive and shared (reader/writer) modes. Explain the edge cases and how your design handles them."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Sequential Nodes for Fairness**: FIFO ordering prevents starvation
- **Ephemeral Nodes**: Automatic release on client failure
- **Watch Mechanism**: Efficient notification without polling
- **Reader/Writer Separation**: Different node prefixes or paths
- **Reentrancy**: Same client can acquire multiple times

### expected_keywords
- Primary keywords: lock, fairness, ephemeral, sequential, reader, writer
- Technical terms: starvation, deadlock, reentrancy, mutual exclusion

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Lock Timeouts**: Maximum hold time to prevent indefinite locks
- **Lock Stealing**: Admin override for stuck locks
- **Deadlock Detection**: Tracking lock dependency graphs
- **Performance Optimization**: Batching, connection pooling
- **Monitoring**: Lock wait times, holder information
- **Graceful Degradation**: Fallback when Zookeeper unavailable
- **Priority Locks**: High-priority clients jump queue

### bonus_keywords
- Patterns: read-write locks, upgradeable locks, try-lock semantics
- Implementation: client ID tracking, session management, lock tokens
- Production: metrics, alerting, SLAs, debugging tools

## sample_excellent - Example Excellence
"Here's a complete fair distributed lock implementation:

Data Structure:
```
/locks/resource-1/
  ├── write-000001 (ephemeral, sequential)
  ├── read-000002 (ephemeral, sequential)
  ├── read-000003 (ephemeral, sequential)
  └── write-000004 (ephemeral, sequential)
```

Exclusive (Write) Lock Algorithm:
1. Create `/locks/resource-1/write-` (ephemeral, sequential)
2. Get all children, sort by sequence
3. If I'm the lowest number, acquire lock
4. Otherwise, watch the node immediately before me
5. On watch trigger, goto step 2

Shared (Read) Lock Algorithm:
1. Create `/locks/resource-1/read-` (ephemeral, sequential)
2. Get all children, sort by sequence
3. If no write nodes exist before me, acquire lock
4. Otherwise, watch the last write node before me
5. On watch trigger, goto step 2

Reentrancy:
- Track client ID in node data: `write-000001[clientId=abc]`
- Before creating new node, check if I already hold lock
- Maintain reference count for same client

Edge Cases Handled:
1. **Client Crash**: Ephemeral node auto-deleted, next waiter notified
2. **Network Partition**: Session timeout causes ephemeral cleanup
3. **Starvation**: Sequential ordering ensures FIFO fairness
4. **Deadlock**: Can't happen with single resource; for multiple, use lock ordering
5. **Thundering Herd**: Watch only specific predecessor, not all nodes
6. **Connection Loss**: Store node name client-side for recovery
7. **Zombie Clients**: Session timeout prevents indefinite holding

Reader/Writer Semantics:
- Multiple readers can hold simultaneously (no writes before them)
- Writers need exclusive access (no readers or writers before them)
- Writers don't starve: new readers queue behind pending writers

Advanced Features:
```python
class DistributedLock:
    def try_lock(timeout_ms):
        # Non-blocking acquire with timeout
        
    def lock_with_lease(lease_time_ms):
        # Auto-release after lease expires
        
    def upgrade_lock():
        # Convert read lock to write lock
        
    def get_lock_info():
        # Return current holder, queue length, wait estimate
```

Monitoring:
- Track metrics: acquisition time, hold time, queue length
- Alert on: locks held > 5 min, queue > 100, repeated failures
- Debug info: Who holds lock, waiting queue, acquisition history

This provides strong fairness guarantees while handling all failure modes gracefully."

## sample_acceptable - Minimum Acceptable
"Create ephemeral sequential nodes for lock requests. For exclusive locks, only the lowest sequence number can acquire. For shared locks, multiple readers can acquire if no writers are ahead. Use watches to get notified when to retry. Ephemeral nodes ensure locks are released if clients crash. Track the created node name for reentrancy checks."

## common_mistakes - Watch Out For
- Not preventing writer starvation with continuous readers
- Watching all nodes instead of predecessors
- No reentrancy support
- Not handling session recovery
- Missing reader/writer distinction

## follow_up_excellent - Depth Probe
**Question**: "How would you extend this to support distributed transactions across multiple resources while preventing deadlocks?"
- **Looking for**: Lock ordering, 2PL protocol, deadlock detection/prevention, timeout strategies
- **Red flags**: Not considering deadlock, no ordering strategy

## follow_up_partial - Guided Probe  
**Question**: "What happens if a reader is waiting but writers keep arriving? How do you ensure readers eventually get access?"
- **Hint embedded**: Need policy for reader/writer fairness
- **Concept testing**: Understanding starvation scenarios

## follow_up_weak - Foundation Check
**Question**: "Think about a library with study rooms. How would you fairly manage who gets to use them?"
- **Simplification**: Real-world queuing and fairness
- **Building block**: FIFO and resource allocation

## bar_raiser_question - L4→L5 Challenge
"Your distributed lock is being used by a critical payment system. Add capabilities for: (1) detecting and breaking deadlocks across multiple resources, (2) supporting priority inheritance to prevent priority inversion, and (3) providing a 'lock doctor' tool for production debugging. How do you implement these while maintaining correctness?"

### bar_raiser_concepts
- Wait-for graphs for deadlock detection
- Priority inheritance protocol
- Lock dependency tracking
- Diagnostic snapshots
- Safe lock breaking mechanisms
- Compensation for broken locks

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 5-7 min discussion
- **Common next topics**: Database locking, transaction protocols, consensus algorithms
