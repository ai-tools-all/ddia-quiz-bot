---
id: zookeeper-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: zookeeper
subtopic: leader-election
estimated_time: 8-10 minutes
---

# question_title - Implementing Leader Election with Zookeeper

## main_question - Core Question
"Implement a leader election algorithm using Zookeeper primitives. Explain how your solution handles network partitions, leader crashes, and ensures there's always exactly one leader."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Sequential Ephemeral Nodes**: Create ordered, temporary nodes for candidates
- **Lowest Sequence Wins**: Node with smallest sequence number becomes leader
- **Watch Previous Node**: Each node watches the one before it in sequence
- **Automatic Cleanup**: Ephemeral nodes disappear on disconnect/crash

### expected_keywords
- Primary keywords: leader, election, ephemeral, sequential, watch
- Technical terms: znode, session, sequence number, predecessor

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Herd Effect Prevention**: Why watch only previous, not all nodes
- **Session Timeout**: Configuring appropriate timeouts for failure detection
- **Leader Duties**: What leader should write to announce leadership
- **Split Brain Prevention**: How ephemeral nodes prevent multiple leaders
- **Graceful Handoff**: Voluntary leader stepdown patterns
- **Leader Epochs**: Versioning leadership terms

### bonus_keywords
- Implementation: session management, connection handling, retries
- Patterns: bully algorithm comparison, consensus relationship
- Edge cases: clock skew, network delays, GC pauses

## sample_excellent - Example Excellence
"Here's my leader election implementation using Zookeeper:

Algorithm:
1. Each candidate creates an ephemeral sequential node: `/election/node-000000001`
2. Get all children of `/election` and sort by sequence number
3. If you created the lowest numbered node, you're the leader
4. Otherwise, watch the node immediately before yours
5. When watched node disappears, go back to step 2

Key Properties:
- **Exactly One Leader**: Only the lowest sequence number can be leader
- **Automatic Failover**: When leader crashes, its ephemeral node vanishes, next-in-line gets notified via watch
- **No Split Brain**: Even during network partition, ephemeral nodes ensure crashed leader's node disappears

Network Partition Handling:
- If leader gets partitioned, its session times out, ephemeral node deleted
- Remaining nodes elect new leader from their partition
- Original leader discovers it's no longer leader when reconnecting (its node is gone)

Implementation Details:
```
create('/election/node-', ephemeral=true, sequential=true)
while true:
    children = get_children('/election')
    sorted_children = sort(children)
    if my_node == sorted_children[0]:
        become_leader()
        write('/current-leader', my_id)  # Announce leadership
    else:
        predecessor = find_predecessor(sorted_children, my_node)
        exists(predecessor, watch=true)
        wait_for_watch()
```

Optimizations:
- Watch only predecessor to avoid herd effect (O(1) vs O(n) notifications)
- Use reasonable session timeout (e.g., 10-30 seconds) balancing failure detection vs false positives
- Leader writes to `/current-leader` so others can discover who leads
- Add epoch/generation counter to prevent stale commands

This leverages Zookeeper's core guarantees: sequential ordering, ephemeral cleanup, and reliable watches."

## sample_acceptable - Minimum Acceptable
"Create ephemeral sequential nodes in `/election` directory. Each service gets a unique sequence number. The service with the lowest number becomes leader. Other services watch the node just before theirs. When a leader crashes, its ephemeral node disappears, and the next service gets notified through its watch and becomes the new leader."

## common_mistakes - Watch Out For
- Watching all nodes (causes thundering herd)
- Not using ephemeral nodes (orphaned nodes on crash)
- Not handling reconnection scenarios
- Forgetting sequential flag (no ordering)
- Not considering session timeout implications

## follow_up_excellent - Depth Probe
**Question**: "How would you prevent a 'zombie leader' - a leader that thinks it's still leader but the rest of the cluster has moved on?"
- **Looking for**: Epoch numbers, leader lease, heartbeats to Zookeeper, fencing tokens
- **Red flags**: Only relying on network timeouts

## follow_up_partial - Guided Probe  
**Question**: "Why watch only the previous node instead of watching the current leader node directly?"
- **Hint embedded**: What happens when many nodes watch one node
- **Concept testing**: Understanding the herd effect

## follow_up_weak - Foundation Check
**Question**: "In a group project, how would you choose who presents? What problems might occur if the presenter suddenly leaves?"
- **Simplification**: Real-world leader selection
- **Building block**: Single leader necessity

## bar_raiser_question - L4→L5 Challenge
"Extend your leader election to support 'warm standby' - the next-in-line pre-loads leader state for faster failover. How do you ensure consistency during the handoff?"

### bar_raiser_concepts
- State replication mechanisms
- Fence tokens for old leader
- Checkpoint consistency
- Dual writes during transition
- Client request routing

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 5-6 min discussion
- **Common next topics**: Raft consensus, primary-backup replication, service discovery
