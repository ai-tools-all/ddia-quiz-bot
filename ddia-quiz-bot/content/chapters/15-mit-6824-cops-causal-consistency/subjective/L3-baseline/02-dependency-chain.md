---
id: cops-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: cops
subtopic: dependency-visibility
estimated_time: 6-8 minutes
---

# question_title - Dependency Satisfaction and Visibility

## main_question - Core Question
"When a remote shard server receives a put operation with dependencies, what does it do before making the put visible to local gets? Why is this necessary?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Remote server does NOT immediately make put visible
- Server waits until all dependencies in the list are satisfied
- A dependency is satisfied when that key-version exists and is visible locally
- Gets return highest visible version
- This ensures causal ordering: readers see all prior causally-related writes

### expected_keywords
- deferred visibility, dependency satisfaction, visible version, causal ordering

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Prevents photo-list anomaly (seeing list entry before photo)
- Dependency checking is local to the shard server
- Transitivity of dependencies

### bonus_keywords
- anomaly prevention, causal consistency guarantee, local checking

## sample_excellent - Example Excellence
"When a remote shard server receives a replicated put with dependencies, it does not immediately make this put visible to clients. Instead, it checks whether all dependencies (key-version pairs) in the attached list are satisfied—meaning those specific versions exist and are already visible locally. Only when all dependencies are met does the server make this new put visible to subsequent gets. This deferred visibility mechanism ensures that any client reading the new version will have already seen (or will see) all causally prior writes, preventing anomalies like reading a photo list entry before the photo itself arrives."

## sample_acceptable - Minimum Acceptable
"Server waits for all dependencies to be satisfied (present and visible) before making the put visible. This preserves causal order so clients don't see effects before causes."

## common_mistakes - Watch Out For
- Saying server immediately applies the put
- Not explaining what "satisfied" means for a dependency
- Missing the connection to causal consistency guarantees

## follow_up_excellent - Depth Probe
**Question**: "What happens if a dependency never arrives due to a network partition?"
- **Looking for**: Availability sacrifice, indefinite stall, consistency over availability

## follow_up_partial - Guided Probe
**Question**: "How does deferred visibility solve the photo-list example from the paper?"
- **Hint embedded**: Order of visibility at remote site

## follow_up_weak - Foundation Check
**Question**: "What is a 'visible' version versus a version that has been received?"
