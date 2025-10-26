---
id: fault-tolerance-subjective-L5-004
type: subjective
level: L5
category: baseline
topic: fault-tolerance
subtopic: read-paths-comparison
estimated_time: 7-9 minutes
---

# question_title - Read Paths: Commit, ReadIndex, and Leases

## main_question - Core Question
"Compare leader-commit reads, ReadIndex, and lease-based reads in Raft. For each, explain correctness assumptions, latency characteristics, and failure modes. When would you use each in production?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Leader-Commit Reads**: Strongest, piggyback on committed log
- **ReadIndex**: Confirm leader’s commit index without appending
- **Leases**: Time-bounded authority with clock assumptions

### expected_keywords
- Primary: linearizability, commit index, lease, clock skew
- Technical: heartbeat, quorum, bounded drift, fencing

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Failure Windows**: GC pauses, clock step, partition
- **Fencing Tokens**: Prevent stale leaders serving reads/writes
- **Follower Reads**: Via ReadIndex confirmation

### bonus_keywords
- Implementation: lease timeout, maxSkew, ReadIndex RPC
- Related: Spanner TrueTime, Jepsen tests

## sample_excellent - Example Excellence
"Commit reads route through the log and wait for commitIndex ≥ read point—highest safety, highest latency. ReadIndex avoids appending by confirming the leader’s commitIndex via a quorum round-trip; safety without clocks but requires leader availability. Lease reads are fastest but assume bounded clock skew and uninterrupted lease validity; unsafe if the leader’s clock runs fast or the node pauses, potentially overlapping with a new leader. Choose commit reads for critical mutations and compliance, ReadIndex for most reads in stable clusters, and leases only with robust time discipline and monitoring."

## sample_acceptable - Minimum Acceptable
"Commit reads are safest; ReadIndex reduces latency; leases are fastest but risk inconsistency if clocks or pauses are bad."

## common_mistakes - Watch Out For
- Assuming leases are safe without synchronized clocks
- Serving follower reads without commit confirmation
- Ignoring GC/stop-the-world effects on leases

## follow_up_excellent - Depth Probe
**Question**: "Design guardrails that make lease reads safe enough in practice."
- **Looking for**: maxSkew bounds, short leases, NTP/PTP, pause detector, fencing tokens, rapid revocation

## follow_up_partial - Guided Probe  
**Question**: "How can a follower serve a linearizable read?"
- **Hint**: Confirm commitIndex with ReadIndex or forward to leader

## follow_up_weak - Foundation Check
**Question**: "Why might a very fast read be unsafe?"
- **Simplification**: It might come from a stale replica

## bar_raiser_question - L5→L6 Challenge
"Propose an adaptive read path that switches among commit, ReadIndex, and leases based on health signals and load. What signals trigger transitions?"

### bar_raiser_concepts
- Health metrics: heartbeat RTTs, clock drift, pause time, leader stability
- SLOs: tail latency budgets vs correctness risk
- Safe fallbacks and hysteresis

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 3-4 min discussion
- **Common next topics**: Session guarantees, Jepsen scenarios, production SLOs

## assistant_answer
Commit reads for maximum safety; ReadIndex for common-case low-latency linearizable reads; leases only with strong clock bounds and pause detection. In production, default to ReadIndex, fall back to commit under instability, and use leases opportunistically when clocks are tight and leaders stable.

## improvement_suggestions
- Add specific thresholds (e.g., max skew, lease ≤ skew×2) and instrumentation requirements.
- Include a decision table mapping health signals to read path choices.

## improvement_exercises
### exercise_1 - Thresholds and Instrumentation
**Question**: "Propose concrete thresholds for lease safety and list instrumentation to enforce them."

**Sample answer**: "maxSkew ≤ 5ms, lease ≤ 2×maxSkew, heartbeat ≤ 100ms; require monotonic NTP/PTP, pause detector < 10ms p99, fencing tokens for write paths, and alarms if drift or pauses exceed thresholds; auto-disable leases on violations."

### exercise_2 - Decision Table
**Question**: "Draft a simple decision table for choosing read path based on health signals (stable leader, drift OK, pauses low, low RTT)."

**Sample answer**: "All green → lease; leader stable but drift unknown → ReadIndex; instability or high RTT → commit reads; any alarm → force commit reads until recovered."
