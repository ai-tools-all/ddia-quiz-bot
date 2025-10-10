---
id: gfs-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: gfs
subtopic: replication
estimated_time: 5-7 minutes
---

# question_title - GFS Replication Strategy

## main_question - Core Question
"Explain how GFS ensures data durability when a chunk server fails. Walk me through what happens from the moment of failure."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **3x Replication**: GFS maintains 3 copies of each chunk by default
- **Master Detection**: Master detects failure through heartbeats
- **Re-replication**: Master initiates copying to restore replication factor

### expected_keywords
- Primary keywords: replica, chunk server, master, failure
- Technical terms: heartbeat, replication factor

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Rack Awareness**: Replicas spread across different racks
- **Priority Queue**: More under-replicated chunks get priority
- **Version Numbers**: Ensures stale replicas aren't used
- **Bandwidth Throttling**: Re-replication doesn't overwhelm network

### bonus_keywords
- Implementation: pipelining, chunk priorities
- Related: Colossus improvements, HDFS similarities
- Metrics: MTTR, availability targets

## sample_excellent - Example Excellence
"When a chunk server fails in GFS, several things happen: First, the master detects the failure when heartbeats stop arriving (typically after 60 seconds). The master then scans its metadata to identify all chunks that had replicas on the failed server. For each under-replicated chunk (now at 2 replicas instead of 3), the master adds it to a re-replication queue, prioritizing chunks with fewer replicas. The master then instructs healthy chunk servers holding replicas to copy their chunks to other available servers, spreading them across different racks when possible. This ensures the 3x replication factor is restored. The system uses version numbers to prevent stale replicas from being used if the 'failed' server comes back online."

## sample_acceptable - Minimum Acceptable
"GFS keeps 3 copies of each chunk on different servers. When a chunk server fails, the master notices through missing heartbeats and starts making new copies of the affected chunks on other servers to get back to 3 replicas."

## common_mistakes - Watch Out For
- Saying "backup" instead of "replica" (backups are different!)
- Forgetting the master's coordination role
- Thinking clients handle re-replication
- Confusing with RAID or local redundancy

## follow_up_excellent - Depth Probe
**Question**: "What happens if the master fails during re-replication? How does GFS handle that?"
- **Looking for**: Shadow master, operation log, checkpointing
- **Red flags**: Saying there's automatic master failover (there isn't!)

## follow_up_partial - Guided Probe  
**Question**: "You mentioned the master knows about failures. How would the master know that a chunk server that was working fine a minute ago has suddenly crashed?"
- **Hint embedded**: Time-based checking mechanism
- **Concept testing**: Heartbeat understanding

## follow_up_weak - Foundation Check
**Question**: "Let's start simpler - if you have important data, how many copies would you keep and where would you put them?"
- **Simplification**: Basic redundancy concept
- **Building block**: Understanding why multiple copies matter

## bar_raiser_question - L3→L4 Challenge
"Now consider this scenario: A rack switch fails, taking down 20 chunk servers at once. Each server held 1000 chunks. How should GFS prioritize the re-replication? What problems might arise?"

### bar_raiser_concepts
- Prioritization algorithms (chunks with fewer replicas first)
- Network bandwidth constraints
- Cascading failure risks
- System-wide resource management

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: CAP theorem, HDFS comparison
