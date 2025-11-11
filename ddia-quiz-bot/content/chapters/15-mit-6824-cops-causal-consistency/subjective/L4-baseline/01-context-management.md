---
id: cops-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: cops
subtopic: context-management
estimated_time: 8-10 minutes
---

# question_title - Client Context Tracking and Propagation

## main_question - Core Question
"Explain how a COPS client library tracks and updates its context as it performs get and put operations. How does the context enable transitive causal dependencies across multiple clients?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Context starts empty, accumulates key-version pairs from gets
- Each get adds returned key-version pair to client's context
- Context is sent with every put as dependency list
- Write depends on all reads observed so far by that client
- Transitive dependencies: if A reads X (written by B who read Y), then A's write depends on both X and Y

### expected_keywords
- context accumulation, get response, put dependencies, transitivity, causal chain

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Context can grow large (optimization: garbage collection of old versions)
- Cross-client causal chains formed through read-write relationships
- Client library transparently manages context
- Different clients maintain separate contexts

### bonus_keywords
- context pruning, library transparency, per-client state, causally prior

## sample_excellent - Example Excellence
"A COPS client library maintains a context set that initially is empty. When the client performs a get operation, the server returns the key-value pair along with its version number, and the library adds this key-version pair to the context. When the client performs a put, it sends the entire current context as the dependency list for that put. This means the write depends on all keys the client has read. Transitive dependencies emerge naturally: if Client A reads key X (written by Client B, who had read Y), then when A writes key Z, Z's dependencies include both X and Y (since X's dependencies included Y). This forms multi-hop causal chains across clients."

## sample_acceptable - Minimum Acceptable
"Context tracks key-version pairs from gets. Each get adds to context, each put sends context as dependencies. Dependencies are transitive because a write inherits dependencies from all keys it read."

## common_mistakes - Watch Out For
- Not explaining how transitivity works across multiple clients
- Forgetting that context accumulates across operations
- Missing that only gets add to context, not puts

## follow_up_excellent - Depth Probe
**Question**: "How would you optimize context size if a client performs millions of reads before a write?"
- **Looking for**: Version pruning, dependency minimization, last-write tracking

## follow_up_partial - Guided Probe
**Question**: "If Client A reads X (written by B who read Y), then A writes Z, what are Z's dependencies?"
- **Hint embedded**: Trace through the causal chain

## follow_up_weak - Foundation Check
**Question**: "What operations add entries to the client's context?"
