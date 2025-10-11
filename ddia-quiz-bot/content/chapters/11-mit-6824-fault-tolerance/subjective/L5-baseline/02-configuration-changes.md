---
id: fault-tolerance-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: fault-tolerance
subtopic: configuration-changes
estimated_time: 7-9 minutes
---

# question_title - Dynamic Membership Changes

## main_question - Core Question
"Explain the challenges of changing cluster membership (adding/removing servers) in a running Raft cluster. Why can't we just immediately switch from old to new configuration, and how does Raft's joint consensus approach solve this?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Two Majorities Problem**: Old and new configs might elect different leaders
- **Joint Consensus**: Transitional config requiring both majorities
- **Two-Phase Change**: Old → Old+New → New
- **Safety During Transition**: No period with two leaders

### expected_keywords
- Primary keywords: configuration, membership, joint consensus, majority
- Technical terms: C-old, C-new, configuration entry

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Single-Server Changes**: Alternative simpler approach
- **Availability Considerations**: Maintaining quorum during changes
- **New Server Catch-up**: Log replication before voting rights
- **Configuration as Log Entry**: Configs committed like any operation
- **Leader Restrictions**: Which config can elect leaders when

### bonus_keywords
- Implementation: non-voting members, learner phase
- Related: consensus reconfiguration, view changes
- Optimizations: pre-vote, transfer leadership

## sample_excellent - Example Excellence
"The core challenge in membership changes is preventing 'split-brain' during transition. If we immediately switched from an old configuration to new, there could be a moment where the old configuration's servers elect one leader while the new configuration's servers elect another, violating safety. For example, switching from 3 servers to 5 servers, the original 3 might maintain their old leader while the 2 new servers plus 1 old elect a different leader. Raft solves this using joint consensus: a two-phase transition through a combined configuration C-old+new. In joint consensus, agreement requires separate majorities from both old AND new configurations. The transition works as: 1) Leader logs C-old+new entry, 2) Once C-old+new commits, no leader can be elected using C-old alone, 3) Leader logs C-new entry, 4) Once C-new commits, transition completes. During joint consensus, the leader must get agreement from both majorities for any decision. This prevents split-brain because there's no point where disjoint majorities can act independently. The configuration changes are treated as special log entries, leveraging Raft's existing consensus mechanism for agreement on the configuration itself."

## sample_acceptable - Minimum Acceptable
"Changing membership is dangerous because during transition, different subsets of servers might recognize different configurations, potentially electing multiple leaders. Raft uses joint consensus where for a period, decisions require majority approval from both old and new configurations. This involves two phases: first committing a combined configuration, then committing the final new configuration. This ensures there's never a time when two leaders can be elected."

## common_mistakes - Watch Out For
- Thinking you can just add/remove servers immediately
- Not understanding why two phases are necessary
- Missing that configs are log entries
- Forgetting about availability during transition

## follow_up_excellent - Depth Probe
**Question**: "During joint consensus with C-old={A,B,C} and C-new={A,B,D,E,F}, what's the minimum number of servers needed for commitment? Which specific combinations work?"
- **Looking for**: Need 2 from old AND 3 from new, minimum 4 servers total
- **Red flags**: Not understanding both majorities required

## follow_up_partial - Guided Probe  
**Question**: "Why does Raft need two phases? Why not just use the joint configuration permanently?"
- **Hint embedded**: Joint consensus is complex and inefficient
- **Concept testing**: Understanding transitional nature

## follow_up_weak - Foundation Check
**Question**: "Imagine changing the rules of a game while people are playing. What problems might occur if not everyone switches to new rules simultaneously?"
- **Simplification**: Inconsistent rule application
- **Building block**: Atomic transitions

## bar_raiser_question - L5→L6 Challenge
"A 3-server cluster {A,B,C} needs to migrate to a completely different set {D,E,F} with no overlap. Design a safe migration strategy. What are the availability implications at each step?"

### bar_raiser_concepts
- Cannot directly transition (no overlap for joint consensus)
- Must add new servers first: {A,B,C} → {A,B,C,D,E,F}
- Then remove old: {A,B,C,D,E,F} → {D,E,F}
- Or use multiple smaller transitions
- Availability risks when removing servers
- Need majority from 6 servers in middle phase

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 3-4 min discussion
- **Common next topics**: Byzantine consensus, blockchain consensus, Multi-Paxos
