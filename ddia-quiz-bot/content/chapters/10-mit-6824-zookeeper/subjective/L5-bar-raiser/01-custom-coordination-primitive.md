---
id: zookeeper-subjective-L5-bar-raiser-001
type: subjective
level: L5
category: bar-raiser
topic: zookeeper
subtopic: custom-primitives
estimated_time: 12-15 minutes
---

# question_title - Designing Custom Coordination Primitives

## main_question - Core Question
"Design and implement a distributed semaphore with weighted permits using Zookeeper primitives. Your semaphore should support fair acquisition, permit weights (e.g., one client needs 3 permits, another needs 1), deadlock detection, and graceful degradation. Walk through the implementation and explain how you handle edge cases."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Permit Tracking**: Total permits and current allocation
- **Weighted Requests**: Variable permit requirements per client
- **Fair Queuing**: FIFO ordering with weight consideration
- **Deadlock Prevention**: Handling partial permit availability
- **Cleanup**: Releasing permits on client failure

### expected_keywords
- Primary keywords: semaphore, permits, weighted, fairness, deadlock
- Technical terms: acquire, release, queue, capacity, allocation

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Priority Classes**: Urgent vs normal requests
- **Permit Preemption**: Reclaiming permits from lower priority
- **Timeout Handling**: Maximum wait times
- **Monitoring**: Queue depth, wait times, utilization
- **Dynamic Resizing**: Adjusting total permits at runtime
- **Starvation Prevention**: Ensuring large requests eventually succeed
- **Distributed Deadlock Detection**: Cross-resource dependency tracking

### bonus_keywords
- Patterns: banker's algorithm, resource allocation graph
- Implementation: atomic operations, compensation logic
- Production: metrics, debugging tools, admin interface

## sample_excellent - Example Excellence
"Here's a complete weighted distributed semaphore implementation:

Data Structure:
```
/semaphore/resource-pool/
├── capacity (data: "100")  # Total permits
├── queue/
│   ├── request-000001 (data: {"client": "A", "weight": 30, "timestamp": T1})
│   ├── request-000002 (data: {"client": "B", "weight": 10, "timestamp": T2})
│   └── request-000003 (data: {"client": "C", "weight": 50, "timestamp": T3})
└── holders/
    ├── client-A (data: {"permits": 30, "acquired": T4})
    └── client-D (data: {"permits": 20, "acquired": T5})
```

Core Algorithm:
```python
class WeightedSemaphore:
    def acquire(client_id, weight, timeout=None):
        # 1. Create queue entry (ephemeral, sequential)
        path = create('/semaphore/resource-pool/queue/request-', 
                     data={'client': client_id, 'weight': weight},
                     ephemeral=True, sequential=True)
        
        while True:
            # 2. Check if we can acquire
            capacity = read('/semaphore/resource-pool/capacity')
            holders = get_children('/semaphore/resource-pool/holders')
            used = sum(read(h)['permits'] for h in holders)
            available = capacity - used
            
            # 3. Get queue and check our position
            queue = sorted(get_children('/semaphore/resource-pool/queue'))
            my_position = queue.index(path)
            
            # 4. Calculate if we can acquire
            can_acquire = True
            needed = 0
            for i in range(my_position + 1):
                request = read(queue[i])
                needed += request['weight']
                if needed > available:
                    can_acquire = False
                    break
            
            if can_acquire and my_position == 0:
                # 5. Acquire permits
                create('/semaphore/resource-pool/holders/' + client_id,
                       data={'permits': weight, 'acquired': now()},
                       ephemeral=True)
                delete(path)  # Remove from queue
                return True
            
            # 6. Wait for change
            if timeout and elapsed > timeout:
                delete(path)  # Clean up our request
                return False
                
            # Watch for changes in holders or queue
            watch_path = find_blocking_element(queue, holders, my_position)
            wait_for_change(watch_path)
    
    def release(client_id):
        delete('/semaphore/resource-pool/holders/' + client_id)
```

Deadlock Prevention:
1. **All-or-Nothing**: Only acquire if full weight available
2. **Timeout**: Prevent infinite waiting
3. **Fairness**: FIFO prevents livelock
4. **Abort-Retry**: Client can cancel and retry with smaller weight

Starvation Prevention:
```python
def anti_starvation_check():
    # Promote long-waiting large requests
    queue = get_queue_with_times()
    for request in queue:
        wait_time = now() - request['timestamp']
        if wait_time > MAX_WAIT and request['weight'] > capacity * 0.5:
            # Reserve capacity for this request
            set_reservation(request['client'], request['weight'])
```

Edge Cases Handled:
1. **Client Crash**: Ephemeral nodes auto-release permits
2. **Partial Availability**: Wait until enough permits free
3. **Weight > Capacity**: Reject immediately or queue for admin
4. **Dynamic Resizing**: Admin can update capacity, trigger re-evaluation
5. **Priority Inversion**: Optional priority boosting for waited requests
6. **Thundering Herd**: Watch specific blockers, not all changes

Advanced Features:
```python
class EnhancedSemaphore(WeightedSemaphore):
    def try_acquire_with_preemption(client_id, weight, priority):
        # Try normal acquire first
        if try_acquire(client_id, weight, timeout=0):
            return True
            
        # Check for preemptible lower priority holders
        victims = find_preemptible_victims(weight, priority)
        if victims:
            preempt(victims)
            return acquire(client_id, weight, timeout=SHORT)
        return False
    
    def resize_capacity(new_capacity):
        # Atomic capacity change with safety checks
        with distributed_lock('/semaphore/admin-lock'):
            current_used = calculate_used()
            if new_capacity < current_used:
                # Options: wait, preempt, or reject
                handle_capacity_reduction()
            update('/semaphore/resource-pool/capacity', new_capacity)
            notify_waiters()
```

Monitoring:
- Queue depth and composition
- Wait time percentiles by weight class
- Utilization over time
- Failed acquisition rate
- Deadlock detection runs

This provides a production-ready weighted semaphore with strong fairness guarantees and comprehensive edge case handling."

## sample_acceptable - Minimum Acceptable
"Create a queue using sequential ephemeral nodes where each node stores the client ID and required permits. Clients check if enough permits are available by summing current holders' permits and comparing to capacity. If available and they're first in queue, acquire by creating a holder node. Otherwise, watch for changes. Release by deleting holder node. Ephemeral nodes ensure cleanup on client failure."

## common_mistakes - Watch Out For
- Not handling weight > total capacity
- No fairness mechanism (FIFO ordering)
- Missing deadlock consideration
- No cleanup on client failure
- Not preventing partial acquisitions

## follow_up_excellent - Depth Probe
**Question**: "How would you extend this to support semaphore groups where acquiring from one affects availability in related semaphores?"
- **Looking for**: Resource dependency graphs, multi-resource transactions, deadlock detection across groups
- **Red flags**: Not considering cross-semaphore deadlocks

## follow_up_partial - Guided Probe  
**Question**: "What happens if a client requests 100 permits but the semaphore only has 50 total capacity?"
- **Hint embedded**: Request can never be satisfied
- **Concept testing**: Understanding impossible requests

## follow_up_weak - Foundation Check
**Question**: "Think about a parking garage with different sized vehicles. How would you fairly manage spaces?"
- **Simplification**: Real-world resource allocation
- **Building block**: Weighted resource management

## bar_raiser_question - L5→L6 Challenge
"Extend your semaphore to support distributed transactions where multiple semaphores must be acquired atomically across different Zookeeper clusters. Design the two-phase commit protocol and recovery mechanisms for partial failures."

### bar_raiser_concepts
- Two-phase commit across clusters
- Distributed transaction coordination
- Saga pattern for compensation
- Cross-cluster consistency
- Failure recovery protocols
- Transaction timeout handling

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 7-9 min discussion
- **Common next topics**: Distributed transactions, resource scheduling, workflow orchestration
