---
id: raft-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: raft-consensus
subtopic: fundamental-concepts
estimated_time: 8-10 minutes
---

# question_title - Understanding Raft's Core Components

## main_question - Core Question
"Explain the three server states in Raft (follower, candidate, leader) and how servers transition between them. What triggers each transition?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Follower State**: Default state, responds to leader
- **Candidate State**: Seeking to become leader
- **Leader State**: Handles all client requests
- **Timeout Triggers**: Election timeout causes follower→candidate
- **Vote Outcomes**: Majority votes for candidate→leader

### expected_keywords
- Primary keywords: follower, candidate, leader, election, timeout
- Technical terms: election timeout, heartbeat, term, vote

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Heartbeat Mechanism**: AppendEntries with no entries
- **Term Discovery**: Higher term causes reversion to follower
- **Split Vote**: No majority leads to new election
- **Randomized Timeouts**: Prevents repeated split votes
- **Initial State**: All servers start as followers
- **Client Redirection**: Non-leaders redirect to leader

### bonus_keywords
- Implementation details: RequestVote RPC, AppendEntries RPC
- Timing: 150-300ms timeout range, heartbeat interval
- Edge cases: Network partition behavior, simultaneous elections

## sample_excellent - Example Excellence
"Raft servers operate in exactly one of three states:

**Follower** (default/passive):
- All servers start as followers
- Respond to RPCs from leaders and candidates
- If election timeout elapses without hearing from leader → become candidate
- Accept log entries from leader
- Vote for at most one candidate per term

**Candidate** (transitional/active):
- Initiated by election timeout while follower
- Increment term, vote for self, request votes from others
- Three possible outcomes:
  1. Win election (receive majority votes) → become leader
  2. Another server wins (receive AppendEntries from new leader) → become follower
  3. Split vote/timeout → start new election as candidate with higher term

**Leader** (active/authoritative):
- Handle all client requests
- Send periodic heartbeats (empty AppendEntries) to maintain authority
- Replicate log entries to followers
- Revert to follower if discover higher term

**Transition triggers**:
- Follower→Candidate: Election timeout (no heartbeat from leader)
- Candidate→Leader: Receive majority votes
- Candidate→Follower: Receive valid AppendEntries or see higher term
- Leader→Follower: Discover server with higher term
- Candidate→Candidate: Election timeout without decision

The randomized election timeouts (typically 150-300ms) are crucial for breaking symmetry and avoiding repeated split votes."

## sample_acceptable - Minimum Acceptable
"Raft has three states: follower, candidate, and leader. Followers are passive and respond to the leader. If a follower doesn't hear from a leader within the election timeout, it becomes a candidate and requests votes. If it gets majority votes, it becomes leader. Leaders handle all client requests and send heartbeats to maintain leadership. Any server reverts to follower if it sees a higher term."

## common_mistakes - Watch Out For
- Missing timeout as transition trigger
- Not mentioning term numbers
- Confusing heartbeats with normal AppendEntries
- Forgetting that leaders can revert to followers

## follow_up_excellent - Depth Probe
**Question**: "How does Raft ensure there's at most one leader per term?"
- **Looking for**: Majority requirement, single vote per term, persistent vote storage
- **Red flags**: Not understanding the mathematical impossibility of two majorities

## follow_up_partial - Guided Probe
**Question**: "What happens if two candidates start elections at exactly the same time?"
- **Hint embedded**: Split vote possibility, randomized timeouts
- **Concept testing**: Understanding election conflict resolution

## follow_up_weak - Foundation Check
**Question**: "Why do you think all servers start as followers rather than immediately trying to become leader?"
- **Simplification**: System initialization reasoning
- **Building block**: Understanding default passive behavior

## bar_raiser_question - L4→L5 Challenge
"How would you modify Raft's state transitions to support a 'read-only' mode where the system can serve reads but not accept writes during network partitions?"

### bar_raiser_concepts
- Degraded operation modes
- Read availability vs write consistency
- Partition detection mechanisms
- State machine separation
- Client experience during degradation

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: Paxos comparison, leader election optimizations, Byzantine fault tolerance
