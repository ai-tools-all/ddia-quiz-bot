---
id: cops-subjective-L3-003
type: subjective
level: L3
category: advanced
topic: cops
subtopic: dependency-transitivity
estimated_time: 6-8 minutes
---

# question_title - Transitive Dependency Chains

## main_question - Core Question
"Explain what 'transitive dependencies' means in COPS. Give a concrete example with 3 clients and 3 keys showing how a dependency chain forms across clients and how this ensures causal consistency."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Transitivity: if A depends on B, and B depends on C, then A depends on C
- Example scenario: Client 1 writes X, Client 2 reads X and writes Y, Client 3 reads Y and writes Z
- Z's dependency list includes both Y and X (transitive)
- Ensures proper ordering: remote replicas see X before Y before Z
- Dependency inheritance through context accumulation

### expected_keywords
- transitive, dependency chain, causal chain, context propagation, inheritance

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Multi-hop causality across datacenters
- Dependency list grows with chain depth
- Difference between direct and transitive dependencies

### bonus_keywords
- multi-hop, causal history, direct dependency, inherited dependency

## sample_excellent - Example Excellence
"Transitive dependencies mean if operation A causally precedes B, and B causally precedes C, then A causally precedes C.

**Example:**
- Client 1 (DC-East) writes photo X:v1
- Client 2 (DC-West) reads X:v1, then writes album list Y:v1 with dependency {X:v1}
- Client 3 (DC-Asia) reads Y:v1 (which has dep {X:v1}), then writes notification Z:v1

When Client 3 reads Y:v1, the context gains {X:v1, Y:v1}. So Z:v1's dependency list is {X:v1, Y:v1}—it inherits X:v1 transitively through Y.

At remote DC-Europe, when Z:v1 arrives, it cannot become visible until both Y:v1 AND X:v1 are visible. This ensures correct causal order: users see photo (X), then album list (Y), then notification (Z). Transitivity prevents 'notification about an album that references a photo you can't see.'"

## sample_acceptable - Minimum Acceptable
"Transitive dependencies mean A→B→C implies A→C. Example: Client 1 writes X, Client 2 reads X writes Y, Client 3 reads Y writes Z. Z depends on both Y and X transitively. Remote replicas must see X before Y before Z."

## common_mistakes - Watch Out For
- Not providing concrete example with 3 operations
- Missing that dependencies are inherited through context
- Not explaining why transitivity is necessary for correctness

## follow_up_excellent - Depth Probe
**Question**: "What happens to the dependency list size as causal chains grow longer?"
- **Looking for**: Linear growth, metadata overhead, need for optimization

## follow_up_partial - Guided Probe
**Question**: "If Client 3 directly writes Z without reading Y, does Z depend on X?"
- **Hint embedded**: Causal vs concurrent

## follow_up_weak - Foundation Check
**Question**: "What does 'A causally precedes B' mean?"
