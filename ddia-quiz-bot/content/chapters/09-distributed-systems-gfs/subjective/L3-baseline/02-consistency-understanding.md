---
id: gfs-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: gfs
subtopic: consistency
estimated_time: 5-7 minutes
---

# question_title - GFS Consistency Model

## main_question - Core Question
"A developer complains that two clients reading the same file in GFS sometimes see different data. Is this a bug? Explain why this might happen."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Relaxed Consistency**: GFS doesn't guarantee strong consistency
- **Concurrent Writes**: Can lead to undefined regions
- **Application Responsibility**: Apps must handle inconsistencies

### expected_keywords
- Primary keywords: consistency, replicas, concurrent, undefined
- Technical terms: relaxed consistency, defined regions

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Record Append**: Atomic operation for better consistency
- **Stale Reads**: Reading from outdated replicas
- **Version Numbers**: Mechanism to detect staleness
- **Performance Trade-off**: Consistency vs availability/performance

### bonus_keywords
- Implementation: primary replica, lease
- Related: CAP theorem, eventual consistency
- Use cases: MapReduce tolerance

## sample_excellent - Example Excellence
"This is not a bug but expected behavior in GFS's relaxed consistency model. GFS prioritizes availability and performance over strong consistency. When multiple clients write concurrently to the same region, it becomes 'undefined' - different replicas may have different data. Additionally, clients might read from different chunk server replicas that haven't fully synchronized. GFS made this design choice because its target applications like MapReduce can tolerate or work around these inconsistencies. For operations requiring stronger guarantees, GFS provides record append which is atomic."

## sample_acceptable - Minimum Acceptable
"It's not a bug. GFS uses relaxed consistency, which means different clients might see different data temporarily. This is a trade-off GFS made to achieve better performance and availability."

## common_mistakes - Watch Out For
- Thinking GFS provides strong consistency like traditional filesystems
- Not understanding this is a deliberate design choice
- Confusing with network errors or actual bugs
- Missing the application-level handling aspect

## follow_up_excellent - Depth Probe
**Question**: "Given this consistency model, how would you design an application to work reliably on GFS?"
- **Looking for**: Idempotent operations, checksums, record append usage
- **Red flags**: Trying to force strong consistency at app level

## follow_up_partial - Guided Probe
**Question**: "You mentioned it's for performance. What would happen if GFS tried to keep all replicas perfectly synchronized all the time?"
- **Hint embedded**: Coordination overhead
- **Concept testing**: Understanding trade-offs

## follow_up_weak - Foundation Check
**Question**: "Imagine you and a friend are both editing the same Google Doc. What are different ways the system could handle your simultaneous edits?"
- **Simplification**: Real-world analogy
- **Building block**: Consistency models basics

## bar_raiser_question - L3→L4 Challenge
"Your team needs to build a user authentication service storing session tokens in GFS. What problems would you face with GFS's consistency model, and how would you work around them?"

### bar_raiser_concepts
- Critical data vs eventual consistency mismatch
- Need for external coordination (Chubby/Zookeeper)
- Alternative storage systems for different requirements
- Architectural patterns for strong consistency needs

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: CAP theorem, consistency levels, distributed consensus
