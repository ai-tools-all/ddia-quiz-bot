---
id: gfs-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: gfs
subtopic: failure-handling
estimated_time: 7-10 minutes
---

# question_title - Primary Replica Failure Handling

## main_question - Core Question
"A primary chunk server holding the lease for chunk C crashes during a write operation. Some secondaries received the mutation, others didn't. Walk through GFS's recovery process and explain what happens to in-flight writes."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Lease Expiration**: Wait 60 seconds for lease timeout
- **Version Number Increment**: Master assigns new version
- **Primary Selection**: Choose new primary from up-to-date replicas
- **In-flight Write Loss**: Some writes may be lost/incomplete

### expected_keywords
- Primary keywords: lease, version number, timeout, primary election
- Technical terms: mutation, replica state, consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Split-Brain Prevention**: Why waiting for lease expiration matters
- **Stale Replica Detection**: Version number comparison
- **Client Retry Logic**: Application-level handling
- **Write Ahead Log**: How primary tracks mutations
- **Network Partition**: vs actual failure scenarios

### bonus_keywords
- Implementation: heartbeat detection, lease renewal
- Edge cases: clock skew, partial writes
- Recovery: checkpoint, operation log

## sample_excellent - Example Excellence
"When the primary crashes, GFS follows a careful recovery process to prevent split-brain: 1) The master detects missing heartbeats but critically waits for the 60-second lease to expire, ensuring the old primary can't accept writes even if it's just network-partitioned. 2) Master increments the chunk version number and selects a new primary from replicas with the highest version. 3) Notifies new primary and secondaries of roles and new version. 4) New primary starts accepting writes. In-flight writes are lost - clients will timeout and retry, but the write might be partially applied to some replicas creating an inconsistent state. This is why GFS provides 'defined' but not 'consistent' semantics. Applications using record append will retry and may see duplicates. The version number increment ensures that if the old primary returns, its stale chunks are detected and garbage collected."

## sample_acceptable - Minimum Acceptable
"The master waits for the lease to expire (60 seconds) to avoid split-brain, then selects a new primary from available replicas and increments the version number. In-flight writes during the failure are lost and clients must retry."

## common_mistakes - Watch Out For
- Not mentioning lease timeout importance
- Missing version number increment
- Assuming writes are preserved
- Forgetting about client-side handling

## follow_up_excellent - Depth Probe
**Question**: "What if the master's clock is faster than the failed primary's clock by 5 seconds? Could this cause data corruption?"
- **Looking for**: Clock synchronization importance, conservative timeout buffers
- **Red flags**: Not recognizing distributed time challenges

## follow_up_partial - Guided Probe
**Question**: "You said the master waits 60 seconds. What bad thing would happen if it immediately assigned a new primary?"
- **Hint embedded**: Two primaries simultaneously
- **Concept testing**: Split-brain understanding

## follow_up_weak - Foundation Check
**Question**: "Imagine you're managing a team and the team lead is unreachable. When do you appoint a new lead - immediately or after waiting?"
- **Simplification**: Real-world leadership analogy
- **Building block**: Timeout reasoning

## bar_raiser_question - L4→L5 Challenge
"Design an improvement to GFS that reduces data loss during primary failures without significantly impacting write performance. Consider trade-offs."

### bar_raiser_concepts
- Chain replication for ordered delivery
- Synchronous replication to one secondary
- Shorter leases with adaptive renewal
- Write-ahead logging with replay capability
- Performance vs durability trade-offs

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-5 min discussion
- **Common next topics**: Raft/Paxos protocols, chain replication, exactly-once semantics
