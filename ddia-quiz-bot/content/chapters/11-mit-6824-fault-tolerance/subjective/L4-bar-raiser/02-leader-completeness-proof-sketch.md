---
id: fault-tolerance-subjective-L4-004
type: subjective
level: L4
category: bar-raiser
topic: fault-tolerance
subtopic: leader-completeness-proof
estimated_time: 8-10 minutes
---

# question_title - Proving Leader Completeness (Sketch)

## main_question - Core Question
"Provide a proof sketch of Raft’s Leader Completeness property and explain how the election restriction and commitment rule interact to preserve it across leader changes."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Leader Completeness**: Committed entries appear in all future leaders’ logs
- **Quorum Intersection**: Majorities intersect
- **Election Restriction**: Up-to-date check (term then index)
- **Commitment Rule**: Commit current-term entries; earlier commit indirectly

### expected_keywords
- Primary: completeness, quorum, majority, safety
- Technical: lastLogTerm, lastLogIndex, commitIndex

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Figure 8 Edge Case**: Why current-term-only is required
- **Inductive Argument**: Prefix grows monotonically
- **Crash/Recovery**: Persistence of term and log

### bonus_keywords
- Implementation: leaderCommit, matchIndex
- Related: Paxos safety, quorum certificates

## sample_excellent - Example Excellence
"Let E be a committed entry at index i in term t. Commitment implies E is stored on a majority M. Any future leader must gain votes from a majority V, and M∩V≠∅. Consider voter x∈M∩V that has E. By the up-to-date rule, x will not vote for any candidate whose log lacks E (their lastLogTerm/index would be behind). Therefore, the elected leader’s log contains E. The current-term-only commit rule ensures that when a leader advances commitIndex, at least one entry from its term is majority-replicated, which—via the same intersection—locks in the prefix up to that index. By induction over successive leaders, the committed prefix is contained in all future leaders’ logs."

## sample_acceptable - Minimum Acceptable
"Committed entries live on a majority; future leaders need a majority. Because the two sets overlap, at least one voter would reject a candidate missing the committed entry. Hence the new leader must have all committed entries."

## common_mistakes - Watch Out For
- Ignoring term-first comparison in up-to-date rule
- Forgetting why current-term-only commits are necessary
- Treating majority as fixed set (it can vary but still intersects)

## follow_up_excellent - Depth Probe
**Question**: "How does the proof fail if voters granted votes based on log length only?"
- **Looking for**: Counterexample where older longer log wins over newer short log and loses committed entry

## follow_up_partial - Guided Probe  
**Question**: "Explain how indirect commitment works and why it matters for the proof."
- **Hint**: Commit current-term entry implies all previous indices are committed

## follow_up_weak - Foundation Check
**Question**: "Why must two majorities overlap?"
- **Simplification**: Pigeonhole principle on >N/2 sets

## bar_raiser_question - L4→L5 Challenge
"Construct a concrete 5-node scenario that would violate safety if leaders were allowed to commit prior-term entries directly. Formalize the violation."

### bar_raiser_concepts
- Safety violation via retroactive uncommit
- Necessity of current-term-only commitment
- Demonstrate quorum intersection reasoning

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Formal proofs, TLA+ modeling, Paxos comparison

## assistant_answer
Committed entries are present on a majority; any elected leader must win a (possibly different) majority, which intersects the first. Voters in the intersection reject candidates lacking the committed entry due to the up-to-date rule, forcing any new leader to include it. Current-term-only commitment locks in the prefix and prevents Figure 8-style retroactive uncommit.

## improvement_suggestions
- Add a small TLA+ style invariant list to anchor the proof.
- Include a concrete counterexample showing failure under length-only voting.

## improvement_exercises
### exercise_1 - Invariant Set
**Question**: "Propose 3 invariants sufficient to maintain Leader Completeness."

**Sample answer**: "(I1) Majority intersection. (I2) Up-to-date voting (term, then index). (I3) CommitIndex advances only when current-term entry is majority-replicated; applied ≤ commitIndex."

### exercise_2 - Length-Only Failure Case
**Question**: "Provide a minimal example where length-only voting elects a leader missing a committed entry."

**Sample answer**: "Committed entry at index 5 term 3 on S1,S2,S3. Candidate S4 has long log ending term 2 at index 6; S5 empty. Length-only lets S4 win without entry @5, losing the committed entry."
